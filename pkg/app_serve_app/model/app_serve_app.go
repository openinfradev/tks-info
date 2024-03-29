package model

import (
	"time"

	uuid "github.com/google/uuid"
	"gorm.io/gorm"
)

// AppServeApp contains information of each AppServe application
type AppServeApp struct {
	ID                 uuid.UUID `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	Name               string
	ContractId         string
	Type               string
	AppType            string
	EndpointUrl        string
	PreviewEndpointUrl string
	TargetClusterId    string
	Status             string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (c *AppServeApp) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return nil
}
