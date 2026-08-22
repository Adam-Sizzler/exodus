//go:build !linux && !darwin

package netx

import "syscall"

// bindToDevice is a no-op on platforms without a per-socket interface binding
// option; the dialer's source-address pin still applies.
func bindToDevice(string) func(network, address string, c syscall.RawConn) error {
	return nil
}
