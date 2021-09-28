package csp_info

import (
  "fmt"
  uuid "github.com/google/uuid"
  "gorm.io/gorm"

  model "github.com/openinfradev/tks-info/pkg/csp_info/model"
  pb "github.com/openinfradev/tks-proto/tks_pb"
)

// Accessor accesses to csp info in-memory data.
type CspInfoAccessor struct {
  db *gorm.DB
}

// NewCspInfoAccessor returns new Accessor to access csp info.
func New(db *gorm.DB) *CspInfoAccessor {
  return &CspInfoAccessor{
    db: db,
  }
}

// Get returns a CSP Info if it exists.
func (x *CspInfoAccessor) GetCSPInfo(id uuid.UUID) (model.CSPInfo, error) {
  var cspInfo model.CSPInfo
  res := x.db.First(&cspInfo, id)
  if res.RowsAffected == 0 || res.Error != nil {
    return model.CSPInfo{}, fmt.Errorf("Could not find CSPInfo with ID: %s", id.String())
  }

  return cspInfo, nil
}

// GetCSPIDsByContractID returns a list of CSP ID by contract ID if it exists.
func (x *CspInfoAccessor) GetCSPIDsByContractID(contractId uuid.UUID) ([]string, error) {
  var cspInfos []model.CSPInfo

  res := x.db.Select("id").Find(&cspInfos, "contract_id = ?", contractId)

  if res.RowsAffected == 0 || res.Error != nil {
    return []string{}, fmt.Errorf("Could not find CSPInfo with contract ID: %s", contractId)
  }

  var idArr []string
  for _, item := range cspInfos {
    idArr = append(idArr, item.ID.String())
  }
  return idArr, nil
}

// Create creates new CSP info with contractID and auth.
func (x *CspInfoAccessor) Create(contractId uuid.UUID, name string, auth string, cspType pb.CspType ) (uuid.UUID, error) {
  cspInfo := model.CSPInfo{ContractID: contractId, Name: name, Auth: auth, CspType: cspType}

  res := x.db.Create(&cspInfo)
  if res.Error != nil {
    nilId, _ := uuid.Parse("")
    return nilId, res.Error
  }

  return cspInfo.ID, nil
}

// Update updates an authentication info for CSP.
func (x *CspInfoAccessor) UpdateCSPAuth(id uuid.UUID, auth string) error {
  res := x.db.Model(&model.CSPInfo{}).
    Where("ID = ?", id).
    Update("Auth", auth)

  if res.Error != nil || res.RowsAffected == 0 {
    return fmt.Errorf("nothing updated in cspInfo for id %s", id.String())
  }

  return nil
}
