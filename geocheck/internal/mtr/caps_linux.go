//go:build linux

package mtr

// datagramSeesTimeExceeded reports whether an unprivileged ICMP datagram
// ("ping") socket receives the TimeExceeded errors that traceroute depends on.
//
// On Linux it does not: ping_rcv() only accepts echo replies, and ICMP errors
// are queued on the socket error queue instead, reachable only through
// IP_RECVERR and MSG_ERRQUEUE. Without a raw socket the path stays invisible,
// so we say so rather than reporting a route of stars.
const datagramSeesTimeExceeded = false

// datagramKeepsICMPID is false because ping_v4_sendmsg overwrites the echo
// identifier with the socket's bound port, so a reply's ID says nothing about
// which probe it answers. In exchange the kernel demultiplexes replies to the
// right socket for us.
const datagramKeepsICMPID = false

const privilegeHint = "run as root, or grant the binary the capability once with " +
	"`sudo setcap cap_net_raw+p " + BinaryPathPlaceholder + "`"
