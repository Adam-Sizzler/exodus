package nodessh

type CreateSshTicketResponse struct {
Ticket           string `json:"ticket"`
Path             string `json:"path"`
ExpiresInSeconds int    `json:"expiresInSeconds"`
}

type EvaluateVaultRequestBody struct {
Blinded string `json:"blinded"`
}

type EvaluateVaultResponse struct {
Evaluated string `json:"evaluated"`
}

type WSMessage struct {
Type string `json:"type"`
Cols int    `json:"cols,omitempty"`
Rows int    `json:"rows,omitempty"`
Data string `json:"data,omitempty"`
}
