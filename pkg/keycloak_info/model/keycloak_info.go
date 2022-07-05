package model

import (
	"time"

	uuid "github.com/google/uuid"
	"gorm.io/gorm"
)

// KeycloakInfo represents a KeycloakInfo data in Database.
type KeycloakInfo struct {
	Id         uuid.UUID `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	ClusterId  string
	Realm      string
	ClientId   string
	Secret     string
	PrivateKey string
	UpdatedAt  time.Time
	CreatedAt  time.Time
}

func (c *KeycloakInfo) BeforeCreate(tx *gorm.DB) (err error) {
	c.Id = uuid.New()
	return nil
}
