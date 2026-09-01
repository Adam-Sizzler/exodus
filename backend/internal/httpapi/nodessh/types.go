package nodessh

// CreateSshTicketResponse is returned by NodeSSHTicketHandler on success.
type CreateSshTicketResponse struct {
	Ticket           string `json:"ticket"`
	Path             string `json:"path"`
	ExpiresInSeconds int    `json:"expiresInSeconds"`
}

// EvaluateVaultRequestBody is the request body for NodeSSHVaultEvaluateHandler.
type EvaluateVaultRequestBody struct {
	Blinded string `json:"blinded"`
}

// EvaluateVaultResponse is returned by NodeSSHVaultEvaluateHandler on success.
type EvaluateVaultResponse struct {
	Evaluated string `json:"evaluated"`
}
// WSMessage was an intermediate type from an earlier port revision — removed (#14).
