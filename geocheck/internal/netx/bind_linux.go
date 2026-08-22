//go:build linux

package netx

import (
	"syscall"

	"golang.org/x/sys/unix"
)

// bindToDevice pins a socket to an interface with SO_BINDTODEVICE so that
// policy routing cannot move the flow elsewhere. It needs CAP_NET_RAW, so a
// permission failure is ignored: the LocalAddr pin set by the dialer already
// selects the right source address.
func bindToDevice(iface string) func(network, address string, c syscall.RawConn) error {
	if iface == "" {
		return nil
	}
	return func(_, _ string, c syscall.RawConn) error {
		var inner error
		if err := c.Control(func(fd uintptr) {
			inner = unix.SetsockoptString(int(fd), unix.SOL_SOCKET, unix.SO_BINDTODEVICE, iface)
		}); err != nil {
			return err
		}
		if inner == unix.EPERM || inner == unix.EACCES {
			return nil
		}
		return inner
	}
}
