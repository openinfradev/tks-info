package model

import (
	"time"

	uuid "github.com/google/uuid"
	pb "github.com/openinfradev/tks-proto/tks_pb"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Application contains endpoints and metadata of each application.
type Application struct {
	ID         uuid.UUID `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	Endpoint   string
	Metadata   datatypes.JSON
	Type       pb.AppType
	AppGroupId string
	UpdatedAt  time.Time
	CreatedAt  time.Time
}

func (c *Application) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return nil
}
