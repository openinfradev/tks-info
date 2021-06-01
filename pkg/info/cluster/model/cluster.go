package model

import (
	"time"

	pb "github.com/sktelecom/tks-proto/pbgo"
	uuid "github.com/google/uuid"
	"gorm.io/gorm"
)

// Cluster represents a kubernetes cluster information.
type Cluster struct {
  ID                uuid.UUID      `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	Name              string
	ContractID        string
	CspID             string
	Status            pb.ClusterStatus
	Conf              *pb.ClusterConf
	Kubeconfig        string
  UpdatedAt         time.Time
  CreatedAt         time.Time
}

func (c *Cluster) BeforeCreate(tx *gorm.DB) (err error) {
  c.ID = uuid.New()
  return nil
}
