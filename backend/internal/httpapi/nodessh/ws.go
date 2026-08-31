package nodessh

import (
"context"
"database/sql"
"encoding/base64"
"encoding/json"
"errors"
"fmt"
"io"
"net"
"net/http"
"strings"
"sync"
"sync/atomic"
"time"

"exodus/internal/config"
"exodus/internal/logger"

"golang.org/x/crypto/ssh"
"nhooyr.io/websocket"
"nhooyr.io/websocket/wsjson"
)

// Client message structs
type SshClientMsg struct {
T         string   `json:"t"`
ID        int      `json:"id,omitempty"`
Host      string   `json:"host,omitempty"`
Port      int      `json:"port,omitempty"`
Username  string   `json:"username,omitempty"`
Cols      int      `json:"cols,omitempty"`
Rows      int      `json:"rows,omitempty"`
Keys      []string `json:"keys,omitempty"`
Signature string   `json:"signature,omitempty"`
Accept    bool     `json:"accept,omitempty"`
Message   string   `json:"message,omitempty"`
}

// Server message structs
type SshServerMsg struct {
T           string  `json:"t"`
ID          int     `json:"id,omitempty"`
PublicKey   string  `json:"publicKey,omitempty"`
Data        string  `json:"data,omitempty"`
Hash        *string `json:"hash,omitempty"`
Algo        string  `json:"algo,omitempty"`
Fingerprint string  `json:"fingerprint,omitempty"`
Key         string  `json:"key,omitempty"`
Code        *int    `json:"code,omitempty"`
Signal      *string `json:"signal,omitempty"`
Message     string  `json:"message,omitempty"`
}

type sessionHub struct {
ws        *websocket.Conn
ctx       context.Context
cancel    context.CancelFunc
reqMu     sync.Mutex
nextID    int32
identCh   map[int]chan []string
signCh    map[int]chan string
hostCh    map[int]chan bool
errCh     map[int]chan string
resizeCh  chan [2]int
openCh    chan SshClientMsg
stdinPipe io.WriteCloser
pipeMu    sync.Mutex
log       *logger.Logger
}

func newSessionHub(parentCtx context.Context, ws *websocket.Conn, log *logger.Logger) *sessionHub {
ctx, cancel := context.WithCancel(parentCtx)
return &sessionHub{
ws:       ws,
ctx:      ctx,
cancel:   cancel,
identCh:  make(map[int]chan []string),
signCh:   make(map[int]chan string),
hostCh:   make(map[int]chan bool),
errCh:    make(map[int]chan string),
resizeCh: make(chan [2]int, 10),
openCh:   make(chan SshClientMsg, 1),
log:      log,
}
}

func (h *sessionHub) getNextID() int {
return int(atomic.AddInt32(&h.nextID, 1))
}

func (h *sessionHub) sendJSON(v any) error {
ctx, cancel := context.WithTimeout(h.ctx, 5*time.Second)
defer cancel()
return wsjson.Write(ctx, h.ws, v)
}

func (h *sessionHub) sendBinary(b []byte) error {
ctx, cancel := context.WithTimeout(h.ctx, 5*time.Second)
defer cancel()
return h.ws.Write(ctx, websocket.MessageBinary, b)
}

func (h *sessionHub) sendError(msg string) {
_ = h.sendJSON(SshServerMsg{T: "error", Message: msg})
}

func (h *sessionHub) setStdinPipe(w io.WriteCloser) {
h.pipeMu.Lock()
defer h.pipeMu.Unlock()
h.stdinPipe = w
}

func (h *sessionHub) writeStdin(data []byte) {
h.pipeMu.Lock()
defer h.pipeMu.Unlock()
if h.stdinPipe != nil {
_, _ = h.stdinPipe.Write(data)
}
}

// NodeSSHWSHandler handles the WebSocket stream bridging xterm.js to Node SSH pty
func NodeSSHWSHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
return func(w http.ResponseWriter, r *http.Request) {
var log *logger.Logger
if cfg != nil && cfg.Logger != nil {
log = cfg.Logger.RoleService("API", "NodeSSH")
}

if log != nil {
log.Debug("Incoming WebSocket connection request", "uri", r.URL.RequestURI(), "remoteAddr", r.RemoteAddr)
}

// Parse subprotocols: "ex, <ticket>, <token>" or query params
protoHeader := r.Header.Get("Sec-WebSocket-Protocol")
parts := strings.Split(protoHeader, ",")
for i := range parts {
parts[i] = strings.TrimSpace(parts[i])
}

ticket := r.URL.Query().Get("ticket")
selectedProto := "ex"
if len(parts) >= 2 {
selectedProto = parts[0]
ticket = parts[1]
} else if len(parts) == 1 && parts[0] != "" {
if parts[0] == "ex" || parts[0] == "rw" {
selectedProto = parts[0]
} else {
ticket = parts[0]
}
}

if ticket == "" {
if log != nil {
log.Warn("WebSocket rejected: missing ticket header or query")
}
http.Error(w, "missing ticket", http.StatusBadRequest)
return
}

ticketInfo, ok := consumeTicket(ticket)
if !ok {
if log != nil {
log.Warn("WebSocket rejected: invalid or expired ticket", "ticket", ticket)
}
http.Error(w, "invalid or expired ticket", http.StatusUnauthorized)
return
}

if log != nil {
log.Debug("Ticket valid", "adminUuid", ticketInfo.AdminUUID, "nodeUuid", ticketInfo.NodeUUID, "subproto", selectedProto)
}

opts := &websocket.AcceptOptions{
Subprotocols:       []string{selectedProto, "ex", "rw"},
InsecureSkipVerify: true,
}

conn, err := websocket.Accept(w, r, opts)
if err != nil {
if log != nil {
log.Error("WebSocket accept error", "error", err)
}
return
}
defer conn.Close(websocket.StatusNormalClosure, "session closed")

hub := newSessionHub(r.Context(), conn, log)
defer hub.cancel()

// Run SSH session in separate goroutine
go runSshSession(hub, db, cfg, ticketInfo)

// Main reader loop
for {
msgType, data, err := conn.Read(hub.ctx)
if err != nil {
if log != nil {
log.Debug("WebSocket client disconnected", "error", err)
}
break
}

if msgType == websocket.MessageBinary {
hub.writeStdin(data)
continue
}

if msgType == websocket.MessageText {
var msg SshClientMsg
if err := json.Unmarshal(data, &msg); err != nil {
if log != nil {
log.Debug("Malformed control message", "raw", string(data), "error", err)
}
continue
}

hub.reqMu.Lock()
switch msg.T {
case "open":
select {
case hub.openCh <- msg:
default:
}
case "identities":
if ch, exists := hub.identCh[msg.ID]; exists {
ch <- msg.Keys
}
case "sign":
if ch, exists := hub.signCh[msg.ID]; exists {
ch <- msg.Signature
}
case "hostkey":
if ch, exists := hub.hostCh[msg.ID]; exists {
ch <- msg.Accept
}
case "error":
if ch, exists := hub.errCh[msg.ID]; exists {
ch <- msg.Message
}
case "resize":
select {
case hub.resizeCh <- [2]int{msg.Rows, msg.Cols}:
default:
}
}
hub.reqMu.Unlock()
}
}
}
}

func runSshSession(hub *sessionHub, db *sql.DB, cfg *config.BackendConfig, ticketInfo TicketInfo) {
var openMsg SshClientMsg
select {
case openMsg = <-hub.openCh:
case <-time.After(15 * time.Second):
if hub.log != nil {
hub.log.Warn("Open message timeout")
}
hub.sendError("Timeout waiting for open message")
hub.cancel()
return
case <-hub.ctx.Done():
return
}

if hub.log != nil {
hub.log.Debug("SSH Open requested", "host", openMsg.Host, "port", openMsg.Port, "username", openMsg.Username, "cols", openMsg.Cols, "rows", openMsg.Rows)
}

// Fetch node address
var nodeAddress string
err := db.QueryRowContext(hub.ctx, `SELECT address FROM nodes WHERE uuid = $1`, ticketInfo.NodeUUID).Scan(&nodeAddress)
if err != nil {
if hub.log != nil {
hub.log.Error("Node not found in database", "error", err, "nodeUuid", ticketInfo.NodeUUID)
}
hub.sendError("Node not found")
hub.cancel()
return
}

targetHost := strings.TrimSpace(openMsg.Host)
targetPort := openMsg.Port
if targetPort == 0 {
targetPort = 22
}

// 2. Setup SSH ClientConfig with Browser Signer
authMethod := ssh.PublicKeysCallback(func() ([]ssh.Signer, error) {
reqID := hub.getNextID()
idCh := make(chan []string, 1)
errCh := make(chan string, 1)

hub.reqMu.Lock()
hub.identCh[reqID] = idCh
hub.errCh[reqID] = errCh
hub.reqMu.Unlock()

defer func() {
hub.reqMu.Lock()
delete(hub.identCh, reqID)
delete(hub.errCh, reqID)
hub.reqMu.Unlock()
}()

if err := hub.sendJSON(SshServerMsg{T: "agent-identities", ID: reqID}); err != nil {
return nil, err
}

select {
case keys := <-idCh:
var signers []ssh.Signer
for _, kStr := range keys {
pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(kStr))
if err != nil {
continue
}
signers = append(signers, &browserSigner{
hub:    hub,
pubKey: pubKey,
})
}
return signers, nil
case errMsg := <-errCh:
return nil, errors.New(errMsg)
case <-time.After(30 * time.Second):
return nil, errors.New("agent identities timeout")
case <-hub.ctx.Done():
return nil, errors.New("session closed")
}
})

sshConfig := &ssh.ClientConfig{
User: openMsg.Username,
Auth: []ssh.AuthMethod{authMethod},
HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
reqID := hub.getNextID()
hostCh := make(chan bool, 1)

hub.reqMu.Lock()
hub.hostCh[reqID] = hostCh
hub.reqMu.Unlock()

defer func() {
hub.reqMu.Lock()
delete(hub.hostCh, reqID)
hub.reqMu.Unlock()
}()

fp := ssh.FingerprintSHA256(key)
keyB64 := base64.StdEncoding.EncodeToString(key.Marshal())

if err := hub.sendJSON(SshServerMsg{
T:           "hostkey",
ID:          reqID,
Algo:        key.Type(),
Fingerprint: fp,
Key:         keyB64,
}); err != nil {
return err
}

select {
case accept := <-hostCh:
if !accept {
return errors.New("host key rejected by user")
}
return nil
case <-time.After(30 * time.Second):
return errors.New("host key prompt timeout")
case <-hub.ctx.Done():
return errors.New("session closed")
}
},
Timeout: 30 * time.Second,
}

addr := fmt.Sprintf("%s:%d", targetHost, targetPort)
if hub.log != nil {
hub.log.Debug("Dialing SSH target", "address", addr, "user", openMsg.Username)
}

client, err := ssh.Dial("tcp", addr, sshConfig)
if err != nil {
if hub.log != nil {
hub.log.Error("SSH dial failed", "address", addr, "error", err)
}
hub.sendError(fmt.Sprintf("SSH connection failed: %v", err))
hub.cancel()
return
}
defer client.Close()

session, err := client.NewSession()
if err != nil {
if hub.log != nil {
hub.log.Error("SSH session creation failed", "error", err)
}
hub.sendError(fmt.Sprintf("SSH session creation failed: %v", err))
hub.cancel()
return
}
defer session.Close()

modes := ssh.TerminalModes{
ssh.ECHO:          1,
ssh.TTY_OP_ISPEED: 14400,
ssh.TTY_OP_OSPEED: 14400,
}

cols := openMsg.Cols
if cols == 0 {
cols = 80
}
rows := openMsg.Rows
if rows == 0 {
rows = 24
}

if err := session.RequestPty("xterm-256color", rows, cols, modes); err != nil {
if hub.log != nil {
hub.log.Error("SSH Request PTY failed", "error", err)
}
hub.sendError(fmt.Sprintf("Request PTY failed: %v", err))
hub.cancel()
return
}

stdinPipe, err := session.StdinPipe()
if err != nil {
hub.sendError(fmt.Sprintf("Stdin pipe error: %v", err))
hub.cancel()
return
}
hub.setStdinPipe(stdinPipe)
defer func() {
hub.setStdinPipe(nil)
stdinPipe.Close()
}()

stdoutPipe, err := session.StdoutPipe()
if err != nil {
hub.sendError(fmt.Sprintf("Stdout pipe error: %v", err))
hub.cancel()
return
}

stderrPipe, err := session.StderrPipe()
if err != nil {
hub.sendError(fmt.Sprintf("Stderr pipe error: %v", err))
hub.cancel()
return
}

if err := session.Shell(); err != nil {
hub.sendError(fmt.Sprintf("Failed to start shell: %v", err))
hub.cancel()
return
}

if hub.log != nil {
hub.log.Debug("SSH shell ready, sending ready message...")
}
_ = hub.sendJSON(SshServerMsg{T: "ready"})

// Handle window resize requests
go func() {
for {
select {
case sz := <-hub.resizeCh:
_ = session.WindowChange(sz[0], sz[1])
case <-hub.ctx.Done():
return
}
}
}()

// Pipe SSH stdout to WebSocket binary messages
go func() {
buf := make([]byte, 4096)
for {
n, err := stdoutPipe.Read(buf)
if n > 0 {
_ = hub.sendBinary(buf[:n])
}
if err != nil {
break
}
}
}()

// Pipe SSH stderr to WebSocket binary messages
go func() {
buf := make([]byte, 4096)
for {
n, err := stderrPipe.Read(buf)
if n > 0 {
_ = hub.sendBinary(buf[:n])
}
if err != nil {
break
}
}
}()

// Wait for session exit or context cancel
_ = session.Wait()

code := 0
_ = hub.sendJSON(SshServerMsg{T: "exit", Code: &code})
}

type browserSigner struct {
hub    *sessionHub
pubKey ssh.PublicKey
}

func (s *browserSigner) PublicKey() ssh.PublicKey {
return s.pubKey
}

func (s *browserSigner) Sign(rand io.Reader, data []byte) (*ssh.Signature, error) {
reqID := s.hub.getNextID()
sigCh := make(chan string, 1)
errCh := make(chan string, 1)

s.hub.reqMu.Lock()
s.hub.signCh[reqID] = sigCh
s.hub.errCh[reqID] = errCh
s.hub.reqMu.Unlock()

defer func() {
s.hub.reqMu.Lock()
delete(s.hub.signCh, reqID)
delete(s.hub.errCh, reqID)
s.hub.reqMu.Unlock()
}()

pubKeyB64 := base64.StdEncoding.EncodeToString(s.pubKey.Marshal())
dataB64 := base64.StdEncoding.EncodeToString(data)
hashAlgo := "sha512"

if err := s.hub.sendJSON(SshServerMsg{
T:         "agent-sign",
ID:        reqID,
PublicKey: pubKeyB64,
Data:      dataB64,
Hash:      &hashAlgo,
}); err != nil {
return nil, err
}

select {
case sigB64 := <-sigCh:
sigBytes, err := base64.StdEncoding.DecodeString(sigB64)
if err != nil {
return nil, err
}
return &ssh.Signature{
Format: s.pubKey.Type(),
Blob:   sigBytes,
}, nil
case errMsg := <-errCh:
return nil, errors.New(errMsg)
case <-time.After(30 * time.Second):
return nil, errors.New("agent sign timeout")
case <-s.hub.ctx.Done():
return nil, errors.New("session closed")
}
}

// SignWithAlgorithm implements ssh.AlgorithmSigner
func (s *browserSigner) SignWithAlgorithm(rand io.Reader, data []byte, algorithm string) (*ssh.Signature, error) {
reqID := s.hub.getNextID()
sigCh := make(chan string, 1)
errCh := make(chan string, 1)

s.hub.reqMu.Lock()
s.hub.signCh[reqID] = sigCh
s.hub.errCh[reqID] = errCh
s.hub.reqMu.Unlock()

defer func() {
s.hub.reqMu.Lock()
delete(s.hub.signCh, reqID)
delete(s.hub.errCh, reqID)
s.hub.reqMu.Unlock()
}()

pubKeyB64 := base64.StdEncoding.EncodeToString(s.pubKey.Marshal())
dataB64 := base64.StdEncoding.EncodeToString(data)
hashAlgo := algorithm

if err := s.hub.sendJSON(SshServerMsg{
T:         "agent-sign",
ID:        reqID,
PublicKey: pubKeyB64,
Data:      dataB64,
Hash:      &hashAlgo,
}); err != nil {
return nil, err
}

select {
case sigB64 := <-sigCh:
sigBytes, err := base64.StdEncoding.DecodeString(sigB64)
if err != nil {
return nil, err
}
return &ssh.Signature{
Format: algorithm,
Blob:   sigBytes,
}, nil
case errMsg := <-errCh:
return nil, errors.New(errMsg)
case <-time.After(30 * time.Second):
return nil, errors.New("agent sign timeout")
case <-s.hub.ctx.Done():
return nil, errors.New("session closed")
}
}
