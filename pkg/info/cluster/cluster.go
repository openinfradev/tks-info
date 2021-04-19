package cluster

import (
	"time"

	pb "github.com/sktelecom/tks-proto/pbgo"
)

// Cluster represents a kubernetes cluster information.
type Cluster struct {
	ID            ID
	Name          string
	ContractID    ID
	CspID         ID
	Status        pb.ClusterStatus
	CreatedTs     time.Time
	LastUpdatedTs time.Time
	Conf          *pb.ClusterConf
	Kubeconfig    string
}

// ID is a global unique ID.
type ID string
