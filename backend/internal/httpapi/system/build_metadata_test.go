package system

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"exodus/internal/config"
	"exodus/internal/logger"
)

func TestMetadataHandlerUsesBuildMetadataEnvironment(t *testing.T) {
	sha := "1234567890abcdef1234567890abcdef12345678"
	buildTime := "2026-04-21T12:34:56Z"
	repositoryURL := "git@github.com:ExampleOrg/exodus-fork.git"

	t.Setenv("EXODUS_VERSION", "1.2.3")
	t.Setenv("EXODUS_BACKEND_COMMIT", sha)
	t.Setenv("EXODUS_FRONTEND_COMMIT", sha)
	t.Setenv("EXODUS_GIT_BRANCH", "main")
	t.Setenv("EXODUS_BUILD_TIME", buildTime)
	t.Setenv("EXODUS_BUILD_NUMBER", "456")
	t.Setenv("EXODUS_REPOSITORY_URL", repositoryURL)

	testLogger, err := logger.NewLogger("none", "UTC", io.Discard)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/system/metadata", nil)
	rec := httptest.NewRecorder()

	MetadataHandler(&config.BackendConfig{Logger: testLogger})(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rec.Code)
	}

	var payload struct {
		Response struct {
			Version string `json:"version"`
			Build   struct {
				Time   string `json:"time"`
				Number string `json:"number"`
			} `json:"build"`
			Git struct {
				RepositoryURL string `json:"repositoryUrl"`
				Backend       struct {
					CommitSHA string `json:"commitSha"`
					Branch    string `json:"branch"`
					CommitURL string `json:"commitUrl"`
				} `json:"backend"`
				Frontend struct {
					CommitSHA string `json:"commitSha"`
					CommitURL string `json:"commitUrl"`
				} `json:"frontend"`
			} `json:"git"`
		} `json:"response"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expectedRepoURL := "https://github.com/ExampleOrg/exodus-fork"
	expectedCommitURL := expectedRepoURL + "/commit/" + sha

	if payload.Response.Version != "1.2.3" {
		t.Fatalf("unexpected version: %q", payload.Response.Version)
	}
	if payload.Response.Build.Time != buildTime {
		t.Fatalf("unexpected build time: %q", payload.Response.Build.Time)
	}
	if payload.Response.Build.Number != "456" {
		t.Fatalf("unexpected build number: %q", payload.Response.Build.Number)
	}
	if payload.Response.Git.RepositoryURL != expectedRepoURL {
		t.Fatalf("unexpected repository URL: %q", payload.Response.Git.RepositoryURL)
	}
	if payload.Response.Git.Backend.CommitSHA != sha {
		t.Fatalf("unexpected backend SHA: %q", payload.Response.Git.Backend.CommitSHA)
	}
	if payload.Response.Git.Backend.Branch != "main" {
		t.Fatalf("unexpected branch: %q", payload.Response.Git.Backend.Branch)
	}
	if payload.Response.Git.Backend.CommitURL != expectedCommitURL {
		t.Fatalf("unexpected backend commit URL: %q", payload.Response.Git.Backend.CommitURL)
	}
	if payload.Response.Git.Frontend.CommitSHA != sha {
		t.Fatalf("unexpected frontend SHA: %q", payload.Response.Git.Frontend.CommitSHA)
	}
	if payload.Response.Git.Frontend.CommitURL != expectedCommitURL {
		t.Fatalf("unexpected frontend commit URL: %q", payload.Response.Git.Frontend.CommitURL)
	}
}

func TestNormalizeRepositoryURL(t *testing.T) {
	tests := map[string]string{
		"git+https://github.com/TeamDominant/exodus.git": "https://github.com/TeamDominant/exodus",
		"git@github.com:Owner/repo.git":                  "https://github.com/Owner/repo",
		"ssh://git@github.com/Owner/repo.git":            "https://github.com/Owner/repo",
		"github.com/Owner/repo.git":                      "https://github.com/Owner/repo",
		"https://token@github.com/Owner/repo.git":        "https://github.com/Owner/repo",
		"unknown": "",
	}

	for raw, expected := range tests {
		if actual := normalizeRepositoryURL(raw); actual != expected {
			t.Fatalf("normalizeRepositoryURL(%q) = %q, want %q", raw, actual, expected)
		}
	}
}

func TestNormalizeVersion(t *testing.T) {
	tests := map[string]string{
		"v26.4.21.n1":               "v26.4.21.n1",
		"26.4.21":                   "26.4.21",
		"0.0.1-18-gd6f6d3bc2":       "0.0.1",
		"0.0.1-18-gd6f6d3bc2-dirty": "0.0.1",
		"latest":                    "unknown",
		"unknown":                   "unknown",
	}

	for raw, expected := range tests {
		if actual := normalizeVersion(raw); actual != expected {
			t.Fatalf("normalizeVersion(%q) = %q, want %q", raw, actual, expected)
		}
	}
}
