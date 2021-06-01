package csp

import (
	"fmt"
	"time"
	uuid "github.com/google/uuid"
  "gorm.io/gorm"

	"github.com/sktelecom/tks-contract/pkg/log"
  model "github.com/sktelecom/tks-contract/pkg/contract/model"
  pb "github.com/sktelecom/tks-proto/pbgo"
)

// Accessor accesses to csp info in-memory data.
type Accessor struct {
  db *gorm.DB
}

// NewCSPAccessor returns new Accessor to access csp info.
//func NewCSPAccessor(db *gorm.DB) *Accessor {
func New(db *gorm.DB) *Accessor {
	return &Accessor{
		db: db,
	}
}

// Get returns a CSP Info if it exists.
// Robert: Is it okay to return CSPInfo by value??
func (c Accessor) GetCSPInfo(id uuid.UUID) (CSPInfo, error) {
  var cspInfo model.CSPInfo
	res := x.db.First(&cspInfo, id)
	if res.RowsAffected == 0 || res.Error != nil {
		return model.CSPInfo{}, fmt.Errorf("Could not find CSPInfo with ID: %s", id)
	}

	return cspInfo, nil
}

// GetCSPIDsByContractID returns a list of CSP ID by contract ID if it exists.
func (c Accessor) GetCSPIDsByContractID(contractId uuid.UUID) ([]uuid.UUID, error) {
  var cspInfo model.CSPInfo

  res := x.db.Select("id").Find(&cspInfo, "contract_id = ?", contractId)

	if res.RowsAffected == 0 || res.Error != nil {
		return &model.CSPInfo{}, fmt.Errorf("Could not find CSPInfo with contract ID: %s", contractId)
	}

  // Robert: does cspInfo contains cspIds only now? need to construct array?




  return cspInfo, nil
}


// List returns a list of CSP Infos in array.
func (c Accessor) List() []CSPInfo {
  var cspInfo model.CSPInfo

	res := x.db.Find(&cspInfo)
	return cspInfo
}

// Create creates new CSP info with contractID and auth.
func (c *Accessor) Create(contractID uuid.UUID, auth string) (uuid.UUID, error) {
	cspInfo := model.cspInfo{ContractID: contractID, Auth: auth}
	err := x.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Create(&cspInfo)
		if res.Error != nil {
			return res.Error
		}
  }

  return cspInfo.ID, nil
}

// Update updates an authentication info for CSP.
// Robert: need to return prev and current??
// Robert: need to return prev and current??
// Robert: need to return prev and current??
func (c *Accessor) Update(id ID, auth string) (prev *pb.cspInfo, curr, error) {
	res := x.db.Model(&model.cspInfo{}).
		Where("ID = ?", id).
		Updates("Auth", auth)

	if res.Error != nil || res.RowsAffected == 0 {
		return nil, nil, fmt.Errorf("nothing updated in cspInfo for id %s", ID)
	}

  return prev, curr, nil
}
