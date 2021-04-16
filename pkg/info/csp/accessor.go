package csp

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sktelecom/tks-contract/pkg/log"
)

// Accessor accesses to csp info in-memory data.
type Accessor struct {
	cspinfos map[ID]CSPInfo
}

// NewCSPAccessor returns new Accessor to access csp info.
func NewCSPAccessor() *Accessor {
	return &Accessor{
		cspinfos: map[ID]CSPInfo{},
	}
}

// Get returns a CSP Info if it exists.
func (c Accessor) Get(id ID) (CSPInfo, error) {
	csp, exists := c.cspinfos[id]
	if !exists {
		return CSPInfo{}, fmt.Errorf("CSP ID %s does not exist.", id)
	}
	return csp, nil
}

// GetCSPIDsByContractID returns a list of CSP ID by contract ID if it exists.
func (c Accessor) GetCSPIDsByContractID(id ID) ([]ID, error) {
	res := []ID{}
	for _, csp := range c.cspinfos {
		if csp.ContractID == id {
			log.Info("same")
			res = append(res, csp.ID)
		}
	}
	if len(res) == 0 {
		return res, fmt.Errorf("CSP for contract id %s does not exist.", id)
	}
	return res, nil
}

// List returns a list of CSP Infos in array.
func (c Accessor) List() []CSPInfo {
	res := []CSPInfo{}

	for _, t := range c.cspinfos {
		res = append(res, t)
	}
	return res
}

// Create creates new CSP info with contractID and auth.
func (c *Accessor) Create(contractID ID, auth string) (ID, error) {
	newCSPID := ID(uuid.New().String())
	if _, exists := c.cspinfos[newCSPID]; exists {
		return "", fmt.Errorf("csp id %s does already exist.", newCSPID)
	}
	c.cspinfos[newCSPID] = CSPInfo{
		ID:            newCSPID,
		ContractID:    contractID,
		Auth:          auth,
		CreatedTs:     time.Now(),
		LastUpdatedTs: time.Now(),
	}
	return newCSPID, nil
}

// Update updates an authentication info for CSP.
func (c *Accessor) Update(id ID, auth string) error {
	if _, exists := c.cspinfos[id]; !exists {
		return fmt.Errorf("CSP ID %s does not exist.", id)
	}
	c.cspinfos[id] = CSPInfo{
		ID:            c.cspinfos[id].ID,
		ContractID:    c.cspinfos[id].ContractID,
		Auth:          auth,
		CreatedTs:     c.cspinfos[id].CreatedTs,
		LastUpdatedTs: time.Now(),
	}
	return nil
}
