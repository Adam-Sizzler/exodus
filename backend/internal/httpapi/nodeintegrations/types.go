package nodeintegrations

import (
	"encoding/json"
	"time"
)

type NodeIntegrationAPI struct {
	UUID        string          `json:"uuid"`
	Name        string          `json:"name"`
	Description *string         `json:"description"`
	Config      json.RawMessage `json:"config"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

type CreateNodeIntegrationRequest struct {
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Config      json.RawMessage `json:"config"`
}

type UpdateNodeIntegrationRequest struct {
	UUID         string           `json:"uuid"`
	Name         *string          `json:"name,omitempty"`
	Description  *string          `json:"description,omitempty"`
	Config       *json.RawMessage `json:"config,omitempty"`
	RestartNodes *bool            `json:"restartNodes,omitempty"`
}
