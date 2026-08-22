package server

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"exodus-node/config"
)

type InstanceLock struct {
	file *os.File
}

// AcquireInstanceLock checks if another instance of exodus-node is already running on the same host.
// It returns an InstanceLock and logs a warning if a duplicate instance is detected.
func AcquireInstanceLock(cfg *config.NodeConfig) *InstanceLock {
	log := cfg.LoggerFor("Bootstrap")

	lockPaths := []string{
		"/var/run/exodus-node.lock",
		"/tmp/exodus-node.lock",
	}

	var lockFile *os.File

	for _, path := range lockPaths {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			continue
		}
		f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			continue
		}

		err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err != nil {
			log.Warn("Multiple node instances detected running on this host! Check active Docker containers to avoid port and resource collisions.", "lock_file", path, "error", err)
			_ = f.Close()
			return nil
		}

		// Write current PID to lock file
		_ = f.Truncate(0)
		_, _ = f.Seek(0, 0)
		_, _ = fmt.Fprintf(f, "%d\n", os.Getpid())

		lockFile = f
		break
	}

	return &InstanceLock{file: lockFile}
}

// Release closes and releases the file lock.
func (l *InstanceLock) Release() {
	if l == nil || l.file == nil {
		return
	}
	_ = syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
	_ = l.file.Close()
}
