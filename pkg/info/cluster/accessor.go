package cluster

import (
	"fmt"
	"time"
  uuid "github.com/google/uuid"
  "gorm.io/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"

  "github.com/sktelecom/tks-contract/pkg/log"
  model "github.com/sktelecom/tks-contract/pkg/contract/model"
  pb "github.com/sktelecom/tks-proto/pbgo"
)

const (
	MAX_RETRY_COUNT = 10
)

// Accessor accesses to csp info in-memory data.
type Accessor struct {
  db *gorm.DB
}

// NewClusterAccessor returns new Accessor to access clusters.
func New(db *gorm.DB) *Accessor {
  return &Accessor{
    db: db,
  }
}

// Get returns a Cluster if it exists.
func (c Accessor) Get(id uuid.UUID) (*pb.Cluster, error) {
  var cluster model.Cluster
  res := x.db.First(&cluster, id)
  if res.RowsAffected == 0 || res.Error != nil {
    return &pb.Cluster{}, fmt.Errorf("Could not find Cluster with ID: %s", id)
  }

	return &cluster, nil
}

// GetClusterIDsByContractID returns a list of cluster ID by contract ID if it exists.
func (c Accessor) GetClusterIDsByContractID(contractId uuid.UUID) ([]uuid.UUID, error) {
  var cluster model.Cluster
  // This is possible, too
  // cluster := model.Cluster{}

  res := x.db.Select("id").Find(&cluster, "contract_id = ?", contractId)

  if res.RowsAffected == 0 || res.Error != nil {
    return &model.Cluster{}, fmt.Errorf("Could not find cluster with contractID: %s", contractId)
  }

  // how to return array of IDs?
  clusterIdArr := []uuid.UUID
  //clusterIdArr := make([]uuid.UUID, len??)

	for _, item := range cluster {
		clusterIdArr = append(clusterIdArr, item)
	}

	return clusterIdArr, nil
}

// GetClusterIDsByCspID returns a list of clusters by CspID if it exists.
// Robert: model and pb are different. what should I return?
func (c Accessor) GetClustersByCspID(cspId ID) (*[]pb.Cluster, error) {
  var cluster model.Cluster

  res := x.db.Find(&cluster, "csp_id = ?", cspId)

  if res.RowsAffected == 0 || res.Error != nil {
    return &pb.Cluster{}, fmt.Errorf("Could not find cluster with cspID: %s", cspId)
  }

  //resultContracts = append(resultContracts, reflectToPbContract(contract, &pbQuota))
  //return ConvertToPbCluster(cluster)

	return &cluster, nil
}

// List returns a list of clusters in array.
func (c Accessor) List() *[]pb.Cluster {
  var cluster model.Cluster

  res := x.db.Find(&cluster)

  //resultContracts = append(resultContracts, reflectToPbContract(contract, &pbQuota))

  return cluster
}

// Create creates new cluster with contract ID, csp ID, name.
func (c *Accessor) Create(contractID ID, cspID uuid.UUID, name string, conf *pb.ClusterConf) (uuid.UUID, error) {
  cluster := model.cluster{
    ContractID: contractID,
    CspID: cspID,
    Name: name,
    Status: pb.ClusterStatus_UNSPECIFIED
  }
  err := x.db.Transaction(func(tx *gorm.DB) error {
    res := tx.Create(&cluster)
    if res.Error != nil {
      return res.Error
    }
  }

  return cluster.ID, nil
}

// Robert: Done until this line

// UpdateStatus updates an status of cluster for Cluster.
func (c *Accessor) UpdateStatus(id uuid.UUID, status pb.ClusterStatus) error {
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

func (c Accessor) ConvertToPbCluster(cluster Cluster) *pb.Cluster {
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
