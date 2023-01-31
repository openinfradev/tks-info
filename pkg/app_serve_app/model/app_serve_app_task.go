package model

import (
	"time"

	uuid "github.com/google/uuid"
	"gorm.io/gorm"
)

// AppServeAppTask contains information of each AppServeApp task.
type AppServeAppTask struct {
	ID             uuid.UUID `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	AppServeAppId  uuid.UUID
	Version        string
	Strategy       string
	Status         string
	Output         string
	ArtifactUrl    string
	ImageUrl       string
	ExecutablePath string
	ResourceSpec   string
	Profile        string
	AppConfig      string
	AppSecret      string
	ExtraEnv       string
	Port           string
	HelmRevision   int32
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (c *AppServeAppTask) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return nil
}
