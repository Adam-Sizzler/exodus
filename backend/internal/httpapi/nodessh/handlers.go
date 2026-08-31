package nodessh

import (
"crypto/rand"
"database/sql"
"encoding/base64"
"encoding/json"
"errors"
"net/http"
"strings"
"sync"
"time"

"exodus/internal/config"
"exodus/internal/httpapi/auth"
"exodus/internal/httpapi/shared"

"github.com/google/uuid"
)

type TicketInfo struct {
AdminUUID string
NodeUUID  string
ExpiresAt time.Time
}

var (
ticketLock sync.RWMutex
ticketMap  = make(map[string]TicketInfo)
)

func storeTicket(ticket, adminUUID, nodeUUID string, ttl time.Duration) {
ticketLock.Lock()
defer ticketLock.Unlock()
ticketMap[ticket] = TicketInfo{
AdminUUID: adminUUID,
NodeUUID:  nodeUUID,
ExpiresAt: time.Now().Add(ttl),
}
}

func consumeTicket(ticket string) (TicketInfo, bool) {
ticketLock.Lock()
defer ticketLock.Unlock()
info, ok := ticketMap[ticket]
if !ok {
return TicketInfo{}, false
}
delete(ticketMap, ticket)
if time.Now().After(info.ExpiresAt) {
return TicketInfo{}, false
}
return info, true
}

func extractNodeUUID(path string) string {
path = strings.Trim(path, "/")
// Case 1: /api/node-ssh/{uuid}/ticket
if strings.HasSuffix(path, "/ticket") {
trimmed := strings.TrimSuffix(path, "/ticket")
parts := strings.Split(trimmed, "/")
if len(parts) > 0 {
return parts[len(parts)-1]
}
}
// Case 2: /api/node-ssh/tickets/{uuid}
if strings.Contains(path, "tickets/") {
parts := strings.Split(path, "tickets/")
if len(parts) > 1 {
return strings.Trim(parts[1], "/")
}
}
// Case 3: {uuid}
parts := strings.Split(path, "/")
if len(parts) > 0 {
last := parts[len(parts)-1]
if _, err := uuid.Parse(last); err == nil {
return last
}
}
return ""
}

// NodeSSHTicketHandler handles issuing tickets for node SSH web terminal (POST /api/node-ssh/:uuid/ticket or /api/node-ssh/tickets/:uuid)
func NodeSSHTicketHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
return auth.RequireAdminRole(func(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
return
}

nodeUUID := extractNodeUUID(r.URL.Path)
if nodeUUID == "" {
var body struct {
NodeUUID string `json:"nodeUuid"`
}
_ = json.NewDecoder(r.Body).Decode(&body)
nodeUUID = body.NodeUUID
}

if _, err := uuid.Parse(nodeUUID); err != nil {
shared.SendError(w, http.StatusBadRequest, "invalid node UUID format", nil, cfg)
return
}

// Verify node exists
var exists int
err := db.QueryRowContext(r.Context(), `SELECT 1 FROM nodes WHERE uuid = $1`, nodeUUID).Scan(&exists)
if err != nil {
if errors.Is(err, sql.ErrNoRows) {
shared.SendAPIError(w, shared.ErrNodeNotFound, cfg)
return
}
shared.SendError(w, http.StatusInternalServerError, "failed to query node", err, cfg)
return
}

// Generate random 32-byte url-safe ticket
buf := make([]byte, 32)
if _, err := rand.Read(buf); err != nil {
shared.SendError(w, http.StatusInternalServerError, "failed to generate ticket", err, cfg)
return
}
ticket := base64.RawURLEncoding.EncodeToString(buf)

adminUUID := ""
if principal, ok := auth.CurrentAuthPrincipal(r.Context()); ok && principal != nil {
adminUUID = principal.AdminUUID
}

ttl := 15 * time.Second
storeTicket(ticket, adminUUID, nodeUUID, ttl)

shared.WriteJSON(w, http.StatusCreated, map[string]any{
"response": CreateSshTicketResponse{
Ticket:           ticket,
Path:             "/api/node-ssh/ws",
ExpiresInSeconds: 15,
},
})
})
}

// NodeSSHVaultEvaluateHandler handles OPRF evaluation of blinded elements (POST /api/node-ssh/vault/evaluate)
func NodeSSHVaultEvaluateHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
return auth.RequireAdminRole(func(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
return
}

var req EvaluateVaultRequestBody
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
return
}

if strings.TrimSpace(req.Blinded) == "" {
shared.SendError(w, http.StatusBadRequest, "blinded cannot be empty", nil, cfg)
return
}

secret := cfg.JWT.AuthSecret
if secret == "" {
secret = "default-exodus-secret"
}

evaluated, err := EvaluateBlindedElement(secret, req.Blinded)
if err != nil {
shared.SendError(w, http.StatusBadRequest, "vault evaluation failed", err, cfg)
return
}

shared.WriteJSON(w, http.StatusOK, map[string]any{
"response": EvaluateVaultResponse{
Evaluated: evaluated,
},
})
})
}

// NodeSSHDispatcherHandler dispatches all /api/node-ssh/ requests
func NodeSSHDispatcherHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
ticketHandler := NodeSSHTicketHandler(db, cfg)
vaultHandler := NodeSSHVaultEvaluateHandler(db, cfg)
wsHandler := NodeSSHWSHandler(db, cfg)

return func(w http.ResponseWriter, r *http.Request) {
path := r.URL.Path
if cfg != nil && cfg.Backend.IsCustom() {
path = strings.TrimPrefix(path, cfg.Backend.Trimmed())
}

if path == "/api/node-ssh/vault/evaluate" {
vaultHandler.ServeHTTP(w, r)
return
}
if path == "/api/node-ssh/ws" {
wsHandler.ServeHTTP(w, r)
return
}
if strings.HasSuffix(path, "/ticket") || strings.Contains(path, "/tickets") {
ticketHandler.ServeHTTP(w, r)
return
}

shared.WriteJSONError(w, http.StatusNotFound, "route not found")
}
}
