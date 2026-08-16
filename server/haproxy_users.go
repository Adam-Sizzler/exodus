package server

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const haproxyUsersFilePath = "/opt/app/haproxy/data/users.csv"

func applyHaproxyModule(modules DeployModulesPayload) (bool, error) {
	if !modules.HaproxyEnabled {
		err := os.Remove(haproxyUsersFilePath)
		switch {
		case err == nil:
			return true, nil
		case os.IsNotExist(err):
			return false, nil
		default:
			return false, fmt.Errorf("remove haproxy users file: %w", err)
		}
	}

	lines := make([]string, 0, len(modules.HaproxyUsers)*3)
	for _, user := range modules.HaproxyUsers {
		username := strings.TrimSpace(user.Username)
		if username == "" {
			continue
		}
		if uuid := strings.TrimSpace(user.VLESSUUID); uuid != "" {
			lines = append(lines, fmt.Sprintf("1,%s,%s", username, uuid))
		}
		if trojan := normalizeTrojanHash(user.TrojanPassword); trojan != "" {
			lines = append(lines, fmt.Sprintf("1,%s,%s", username, trojan))
		}
		if anytls := normalizeAnytlsHash(user.AnytlsPassword); anytls != "" {
			lines = append(lines, fmt.Sprintf("1,%s,%s", username, anytls))
		}
	}

	if err := os.MkdirAll(filepath.Dir(haproxyUsersFilePath), 0o755); err != nil {
		return false, fmt.Errorf("create haproxy data dir: %w", err)
	}

	content := ""
	if len(lines) > 0 {
		content = strings.Join(lines, "\n") + "\n"
	}
	existing, readErr := os.ReadFile(haproxyUsersFilePath)
	if readErr == nil && bytes.Equal(existing, []byte(content)) {
		return false, nil
	}
	if readErr != nil && !os.IsNotExist(readErr) {
		return false, fmt.Errorf("read haproxy users file: %w", readErr)
	}
	if err := os.WriteFile(haproxyUsersFilePath, []byte(content), 0o644); err != nil {
		return false, fmt.Errorf("write haproxy users file: %w", err)
	}

	return true, nil
}

func normalizeTrojanHash(secret string) string {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return ""
	}
	sum := sha256.Sum224([]byte(secret))
	return hex.EncodeToString(sum[:])
}

func normalizeAnytlsHash(secret string) string {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(sum[:])
}
