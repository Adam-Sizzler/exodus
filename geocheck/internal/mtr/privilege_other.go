//go:build !linux

package mtr

// raisePrivilege is a no-op outside Linux, which has no file capabilities:
// a raw socket there is a question of running as root or not.
func raisePrivilege() {}
