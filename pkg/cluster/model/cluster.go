package model

import (
	"time"

	uuid "github.com/google/uuid"
	"github.com/openinfradev/tks-common/pkg/helper"
	pb "github.com/openinfradev/tks-proto/tks_pb"
	"gorm.io/gorm"
)

// Cluster represents a kubernetes cluster information.
type Cluster struct {
	ID           string `gorm:"primarykey"`
	Name         string
	ContractID   string
	CspID        uuid.UUID
	WorkflowId   string
	Status       pb.ClusterStatus
	StatusDesc   string
	SshKeyName   string
	Region       string
	NumOfAz      int32
	MachineType  string
	MinSizePerAz int32
	MaxSizePerAz int32
	Kubeconfig   string
	UpdatedAt    time.Time
	CreatedAt    time.Time
}

func (c *Cluster) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = helper.GenerateClusterId()
	return nil
}
