package csp

import (
	"time"
)

// CSPInfo represents an IaaS information used to deploy K8S clusters.
type CSPInfo struct {
	ID            ID        `json:"id"`
	ContractID    ID        `json:"contract_id"`
	Auth          string    `json:"auth"`
	CreatedTs     time.Time `json:"created_ts"`
	LastUpdatedTs time.Time `json:"last_updated_ts"`
}

// ID is a global unique ID.
type ID string
