package cluster

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	pb "github.com/sktelecom/tks-proto/pbgo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	MAX_RETRY_COUNT = 10
)

// Accessor accesses to csp info in-memory data.
type Accessor struct {
	clusters map[ID]Cluster
}

// NewClusterAccessor returns new Accessor to access clusters.
func NewClusterAccessor() *Accessor {
	return &Accessor{
		clusters: map[ID]Cluster{},
	}
}

// Get returns a Cluster if it exists.
func (c Accessor) Get(id ID) (Cluster, error) {
	cluster, exists := c.clusters[id]
	if !exists {
		return Cluster{}, fmt.Errorf("Cluster ID %s does not exist.", id)
	}
	return cluster, nil
}

// GetClusterIDsByContractID returns a list of cluster ID by contract ID if it exists.
func (c Accessor) GetClusterIDsByContractID(id ID) ([]ID, error) {
	res := []ID{}
	for _, cluster := range c.clusters {
		if cluster.ContractID == id {
			res = append(res, cluster.ID)
		}
	}
	if len(res) == 0 {
		return res, fmt.Errorf("Cluster for contract id %s does not exist.", id)
	}
	return res, nil
}

// GetClusterIDsByCspID returns a list of cluster ID by Csp ID if it exists.
func (c Accessor) GetClustersByCspID(cspId ID) ([]Cluster, error) {
	res := []Cluster{}
	for _, cluster := range c.clusters {
		if cluster.CspID == cspId {
			res = append(res, cluster)
		}
	}
	if len(res) == 0 {
		return res, fmt.Errorf("Cluster for Csp id %s does not exist.", cspId)
	}
	return res, nil
}

// List returns a list of clusters in array.
func (c Accessor) List() []Cluster {
	res := []Cluster{}

	for _, t := range c.clusters {
		res = append(res, t)
	}
	return res
}

// Create creates new cluster with contract ID, csp ID, name.
func (c *Accessor) Create(contractID ID, cspID ID, name string, conf *pb.ClusterConf) (ID, error) {
	newID, err := c.GenerateNewClusterID()
	if err != nil {
		return newID, err
	}
	c.clusters[newID] = Cluster{
		ID:            newID,
		ContractID:    contractID,
		CspID:         cspID,
		Name:          name,
		Status:        pb.ClusterStatus_UNSPECIFIED,
		CreatedTs:     time.Now(),
		LastUpdatedTs: time.Now(),
	}
	return newID, nil
}

// UpdateStatus updates an status of cluster for Cluster.
func (c *Accessor) UpdateStatus(id ID, status pb.ClusterStatus) error {
	if _, exists := c.clusters[id]; !exists {
		return fmt.Errorf("Cluster ID %s does not exist.", id)
	}
	cluster := c.clusters[id]
	c.clusters[id] = Cluster{
		ID:            cluster.ID,
		ContractID:    cluster.ContractID,
		CspID:         cluster.CspID,
		Name:          cluster.Name,
		Status:        status,
		CreatedTs:     cluster.CreatedTs,
		LastUpdatedTs: time.Now(),
	}
	return nil
}

func (c Accessor) ClustertToPbCluster(cluster Cluster) *pb.Cluster {
	return &pb.Cluster{
		Id:         string(cluster.ID),
		Name:       cluster.Name,
		CreatedAt:  timestamppb.New(cluster.CreatedTs),
		UpdatedAt:  timestamppb.New(cluster.CreatedTs),
		Status:     cluster.Status,
		ContractId: string(cluster.ContractID),
		Conf:       cluster.Conf,
		Kubeconfig: cluster.Kubeconfig,
	}
}

// GenerateNewClusterID returns unique ID for cluster.
func (c Accessor) GenerateNewClusterID() (ID, error) {
	for i := 0; i < MAX_RETRY_COUNT; i++ {
		newID := ID(uuid.New().String())
		if _, exists := c.clusters[newID]; !exists {
			return newID, nil
		}
	}
	return ID(""), fmt.Errorf("Failed to generate new cluster ID")
}
