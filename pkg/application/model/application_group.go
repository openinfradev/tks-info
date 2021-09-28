package model

import (
	"time"

	uuid "github.com/google/uuid"
	pb "github.com/openinfradev/tks-proto/tks_pb"
	"gorm.io/gorm"
)

// ApplicationGroup represents an application group data in Database.
type ApplicationGroup struct {
	ID            uuid.UUID `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	Name          string
	Type          pb.AppGroupType
	Status        pb.AppGroupStatus
	ClusterId     uuid.UUID `gorm:"type:uuid;"`
	ExternalLabel string
	UpdatedAt     time.Time
	CreatedAt     time.Time
}

func (c *ApplicationGroup) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return nil
}
