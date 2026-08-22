//go:build !linux

package mtr

// datagramSeesTimeExceeded reports whether an unprivileged ICMP datagram
// socket receives TimeExceeded errors. Darwin and the BSDs deliver them on the
// normal receive path, so unprivileged tracing works as-is.
const datagramSeesTimeExceeded = true

// datagramKeepsICMPID is true on the BSDs: the identifier we set is the one
// that goes out and comes back. It matters here because these kernels hand a
// copy of every inbound ICMP packet to every ICMP socket on the host, so the
// identifier is part of telling our replies from another process's.
const datagramKeepsICMPID = true

const privilegeHint = "run with sudo for a raw socket"
