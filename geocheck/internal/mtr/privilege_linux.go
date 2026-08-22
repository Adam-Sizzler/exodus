//go:build linux

package mtr

import "golang.org/x/sys/unix"

// capNetRaw is CAP_NET_RAW, the capability that permits raw sockets.
const capNetRaw = 13

// raisePrivilege moves CAP_NET_RAW from the permitted set into the effective
// set, which is what actually lets socket(SOCK_RAW) succeed.
//
// The binary is expected to carry the capability as `cap_net_raw+p` rather than
// `+ep`. The difference matters: with the effective bit set by the file
// capability, the kernel refuses to exec the binary at all when NET_RAW is
// outside the bounding set — so `docker run --cap-drop=ALL` would fail before
// main() runs. Carrying it permitted-only keeps the binary runnable everywhere
// and lets it ask for the privilege here, degrading quietly when it has none.
func raisePrivilege() {
	hdr := unix.CapUserHeader{Version: unix.LINUX_CAPABILITY_VERSION_3}
	var data [2]unix.CapUserData
	if err := unix.Capget(&hdr, &data[0]); err != nil {
		return
	}

	idx, mask := capNetRaw/32, uint32(1)<<(capNetRaw%32)
	if data[idx].Permitted&mask == 0 {
		// Not granted; nothing to raise. The caller falls back to a datagram
		// socket and reports the reduced capability.
		return
	}
	if data[idx].Effective&mask != 0 {
		return
	}

	data[idx].Effective |= mask
	_ = unix.Capset(&hdr, &data[0])
}
