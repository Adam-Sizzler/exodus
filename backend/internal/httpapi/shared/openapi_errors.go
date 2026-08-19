package shared

// ErrorResponse represents a standard error response returned by the Exodus API.
// swagger:model
type ErrorResponse struct {
	StatusCode int    `json:"statusCode,omitempty" example:"400"`
	Message    string `json:"message" example:"invalid input"`
	Error      string `json:"error" example:"invalid input"`
	Code       string `json:"code,omitempty" example:"A001"`
	Details    any    `json:"details,omitempty"`
}
