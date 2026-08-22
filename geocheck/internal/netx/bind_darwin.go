//go:build darwin

package netx

import (
	"net"
	"syscall"

	"golang.org/x/sys/unix"
)

// bindToDevice pins a socket to an interface with IP_BOUND_IF / IPV6_BOUND_IF.
// Failures are ignored for the same reason as on Linux: the source-address pin
// is the primary mechanism and this is a belt-and-braces addition.
func bindToDevice(iface string) func(network, address string, c syscall.RawConn) error {
	if iface == "" {
		return nil
	}
	return func(network, _ string, c syscall.RawConn) error {
		ifi, err := net.InterfaceByName(iface)
		if err != nil {
			return err
		}
		level, opt := unix.IPPROTO_IP, unix.IP_BOUND_IF
		if isV6Network(network) {
			level, opt = unix.IPPROTO_IPV6, unix.IPV6_BOUND_IF
		}
		var inner error
		if err := c.Control(func(fd uintptr) {
			inner = unix.SetsockoptInt(int(fd), level, opt, ifi.Index)
		}); err != nil {
			return err
		}
		if inner == unix.EPERM || inner == unix.EACCES {
			return nil
		}
		return inner
	}
}

func isV6Network(network string) bool {
	return len(network) > 0 && network[len(network)-1] == '6'
}
