package model

import (
	"time"

	"github.com/openinfradev/tks-common/pkg/helper"
	pb "github.com/openinfradev/tks-proto/tks_pb"
	"gorm.io/gorm"
)

// ApplicationGroup represents an application group data in Database.
type ApplicationGroup struct {
	ID            string `gorm:"primarykey"`
	Name          string
	Type          pb.AppGroupType
	WorkflowId    string
	Status        pb.AppGroupStatus
	StatusDesc    string
	ClusterId     string
	ExternalLabel string
	UpdatedAt     time.Time
	CreatedAt     time.Time
}

func (c *ApplicationGroup) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = helper.GenerateApplicaionGroupId()
	return nil
}
