package model

import (
	"time"

	uuid "github.com/google/uuid"
	pb "github.com/openinfradev/tks-proto/tks_pb"
	"gorm.io/gorm"
)

// Cluster represents a kubernetes cluster information.
type Cluster struct {
	ID             uuid.UUID `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	Name           string
	ContractID     uuid.UUID
	CspID          uuid.UUID
	Status         pb.ClusterStatus
	MasterFlavor   string
	MasterReplicas int32
	MasterRootSize int64
	WorkerFlavor   string
	WorkerReplicas int32
	WorkerRootSize int64
	K8sVersion     string
	Kubeconfig     string
	UpdatedAt      time.Time
	CreatedAt      time.Time
}

func (c *Cluster) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return nil
}
