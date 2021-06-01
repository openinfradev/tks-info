package model

import (
	"time"

	uuid "github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// CSPInfo represents a CSPInfo data in Database.
type CSPInfo struct {
	ID                uuid.UUID      `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	ContractID        string
	Auth              string
	UpdatedAt         time.Time
	CreatedAt         time.Time
}

func (c *CSPInfo) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return nil
}
