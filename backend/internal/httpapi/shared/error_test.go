package shared

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIErrorConstructorsAndUnwrap(t *testing.T) {
	rootErr := errors.New("database connection failed")

	badReq := NewBadRequestError("invalid payload", rootErr)
	if badReq.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", badReq.StatusCode)
	}
	if !errors.Is(badReq, rootErr) {
		t.Errorf("expected badReq to unwrap to rootErr")
	}

	notLoaded := NewNotFoundError("User", nil)
	if notLoaded.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", notLoaded.StatusCode)
	}
	if notLoaded.Message != "User not found" {
		t.Errorf("expected 'User not found', got %s", notLoaded.Message)
	}

	validation := NewValidationError("validation failed", map[string]string{"field": "email required"})
	if validation.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", validation.StatusCode)
	}
	if validation.Details == nil {
		t.Errorf("expected details to be populated")
	}
}

func TestSendAPIErrorAndWriteJSONError(t *testing.T) {
	// Test SendAPIError
	rec := httptest.NewRecorder()
	apiErr := NewAPIError(http.StatusConflict, "USER_EXISTS", "User already exists", errors.New("duplicate key"))
	SendAPIError(rec, apiErr, nil)

	if rec.Code != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}

	if resp.StatusCode != http.StatusConflict {
		t.Errorf("expected statusCode 409, got %d", resp.StatusCode)
	}
	if resp.Code != "USER_EXISTS" {
		t.Errorf("expected code 'USER_EXISTS', got %s", resp.Code)
	}
	if resp.Message != "User already exists" {
		t.Errorf("expected message 'User already exists', got %s", resp.Message)
	}
	if resp.Error != "User already exists" {
		t.Errorf("expected error 'User already exists', got %s", resp.Error)
	}

	// Test SendError
	rec2 := httptest.NewRecorder()
	SendError(rec2, http.StatusBadRequest, "Invalid UUID format", nil, nil)
	if rec2.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec2.Code)
	}

	var resp2 ErrorResponse
	if err := json.Unmarshal(rec2.Body.Bytes(), &resp2); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}
	if resp2.Message != "Invalid UUID format" || resp2.Error != "Invalid UUID format" {
		t.Errorf("unexpected message in resp2: %+v", resp2)
	}
}
