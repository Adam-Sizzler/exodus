package shared

// ErrorResponse represents a standard error response returned by the Exodus API.
// swagger:model
type ErrorResponse struct {
	Message string `json:"message" example:"invalid input"`
	Error   string `json:"error" example:"invalid input"`
}
