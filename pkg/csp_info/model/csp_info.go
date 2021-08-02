package model

import (
  "time"

  pb "github.com/sktelecom/tks-proto/pbgo"
  uuid "github.com/google/uuid"
  "gorm.io/gorm"
)

// CSPInfo represents a CSPInfo data in Database.
type CSPInfo struct {
  ID                uuid.UUID      `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
  ContractID        uuid.UUID
  Name              string
  Auth              string
  CspType           pb.CspType
  UpdatedAt         time.Time
  CreatedAt         time.Time
}

func (c *CSPInfo) BeforeCreate(tx *gorm.DB) (err error) {
  c.ID = uuid.New()
  return nil
}
