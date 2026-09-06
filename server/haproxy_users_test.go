package server

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildHaproxyUsersContent(t *testing.T) {
	users := []HaproxyUserEntry{
		{
			Username:       "user1",
			VLESSUUID:      "4e2de032-0775-e465-cff8-8bf4f6338f59",
			TrojanPassword: "my-trojan-pass",
			AnytlsPassword: "my-anytls-pass",
			NaivePassword:  "my-naive-pass",
		},
		{
			Username:      "user2",
			NaivePassword: "pass:with:colons",
		},
		{
			Username: "   ", // Should be ignored
		},
	}

	content := buildHaproxyUsersContent(users)
	lines := strings.Split(strings.TrimRight(content, "\n"), "\n")

	expectedToken1 := base64.StdEncoding.EncodeToString([]byte("user1:my-naive-pass"))
	expectedToken2 := base64.StdEncoding.EncodeToString([]byte("user2:pass:with:colons"))

	expectedLines := []string{
		"user1,4e2de032-0775-e465-cff8-8bf4f6338f59",
		"user1," + normalizeTrojanHash("my-trojan-pass"),
		"user1," + normalizeAnytlsHash("my-anytls-pass"),
		"user1,basic:" + expectedToken1,
		"user2,basic:" + expectedToken2,
	}

	if len(lines) != len(expectedLines) {
		t.Fatalf("expected %d lines, got %d. Content:\n%s", len(expectedLines), len(lines), content)
	}

	for i, expected := range expectedLines {
		if lines[i] != expected {
			t.Errorf("line %d: expected %q, got %q", i, expected, lines[i])
		}
		if strings.HasPrefix(lines[i], "1,") || strings.HasPrefix(lines[i], "0,") {
			t.Errorf("line %d has legacy prefix: %q", i, lines[i])
		}
	}
}

func TestApplyHaproxyModule(t *testing.T) {
	tmpDir := t.TempDir()
	origPath := haproxyUsersFilePath
	haproxyUsersFilePath = filepath.Join(tmpDir, "users.csv")
	defer func() {
		haproxyUsersFilePath = origPath
	}()

	payload := DeployModulesPayload{
		HaproxyEnabled: true,
		HaproxyUsers: []HaproxyUserEntry{
			{
				Username:      "alice",
				VLESSUUID:     "11111111-2222-3333-4444-555555555555",
				NaivePassword: "secretalice",
			},
		},
	}

	// 1. Initial write
	changed, err := applyHaproxyModule(payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Fatalf("expected changed=true on initial write")
	}

	data, err := os.ReadFile(haproxyUsersFilePath)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if !strings.Contains(string(data), "alice,11111111-2222-3333-4444-555555555555") {
		t.Errorf("missing vless line in %s", string(data))
	}
	if !strings.Contains(string(data), "alice,basic:") {
		t.Errorf("missing naive line in %s", string(data))
	}

	// 2. Second write with same content -> should return changed=false
	changed, err = applyHaproxyModule(payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changed {
		t.Fatalf("expected changed=false when content is identical")
	}

	// 3. Disable haproxy -> should remove file
	payload.HaproxyEnabled = false
	changed, err = applyHaproxyModule(payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Fatalf("expected changed=true when disabling and removing file")
	}
	if _, err := os.Stat(haproxyUsersFilePath); !os.IsNotExist(err) {
		t.Fatalf("expected file to be removed, stat err: %v", err)
	}
}
