package cluster

import (
  "fmt"
  _ "time"
  uuid "github.com/google/uuid"
  "gorm.io/gorm"
  "google.golang.org/protobuf/types/known/timestamppb"

  _ "github.com/sktelecom/tks-contract/pkg/log"
  model "github.com/sktelecom/tks-info/pkg/cluster/model"
  pb "github.com/sktelecom/tks-proto/pbgo"
)

// Accessor accesses cluster info in DB.
type ClusterAccessor struct {
  db *gorm.DB
}

// NewClusterAccessor returns new Accessor to access clusters.
func New(db *gorm.DB) *ClusterAccessor {
  return &ClusterAccessor{
    db: db,
  }
}


// Get returns a Cluster if it exists.
func (x *ClusterAccessor) GetCluster(id uuid.UUID) (*pb.Cluster, error) {
  var cluster model.Cluster
  res := x.db.First(&cluster, id)
  if res.RowsAffected == 0 || res.Error != nil {
    return &pb.Cluster{}, fmt.Errorf("Could not find Cluster with ID: %s", id)
  }

  pbCluster := ConvertToPbCluster(cluster)
  return pbCluster, nil
}

// GetClusterIDsByContractID returns a list of clusters by ContractID if it exists.
func (x *ClusterAccessor) GetClustersByContractID(contractId uuid.UUID) ([]*pb.Cluster, error) {
  var clusters []model.Cluster

  res := x.db.Find(&clusters, "contract_id = ?", contractId)

  if res.RowsAffected == 0 || res.Error != nil {
    // Robert: Is this correct?
    return []*pb.Cluster{}, fmt.Errorf("Could not find clusters with contractID: %s", contractId)
  }

  pbClusters := []*pb.Cluster{}
  for _, cluster := range clusters {
    pbClusters = append(pbClusters, ConvertToPbCluster(cluster))
  }

  return pbClusters, nil
}

// GetClusterIDsByCspID returns a list of clusters by CspID if it exists.
func (x *ClusterAccessor) GetClustersByCspID(cspId uuid.UUID) ([]*pb.Cluster, error) {
  var clusters []model.Cluster

  res := x.db.Find(&clusters, "csp_id = ?", cspId)

  if res.RowsAffected == 0 || res.Error != nil {
    return []*pb.Cluster{}, fmt.Errorf("Could not find clusters with cspID: %s", cspId)
  }

  pbClusters := []*pb.Cluster{}
  for _, cluster := range clusters {
    pbClusters = append(pbClusters, ConvertToPbCluster(cluster))
  }

  return pbClusters, nil
}

// Create creates new cluster with contract ID, csp ID, name.
func (x *ClusterAccessor) CreateClusterInfo(contractId uuid.UUID, cspId uuid.UUID, name string, conf *pb.ClusterConf) (uuid.UUID, error) {
  cluster := model.Cluster{
    ContractID: contractId,
    CspID: cspId,
    Name: name,
    Status: pb.ClusterStatus_UNSPECIFIED,
    MasterFlavor: conf.MasterFlavor,
    MasterReplicas: conf.MasterReplicas,
    MasterRootSize: conf.MasterRootSize,
    WorkerFlavor: conf.WorkerFlavor,
    WorkerReplicas: conf.WorkerReplicas,
    WorkerRootSize: conf.WorkerRootSize,
    K8sVersion: conf.K8SVersion,
  }

  res := x.db.Create(&cluster)
  if res.Error != nil {
    nilId, _ := uuid.Parse("")
    return nilId, res.Error
  }

  return cluster.ID, nil
}

// UpdateStatus updates an status of cluster for Cluster.
func (x *ClusterAccessor) UpdateStatus(id uuid.UUID, status pb.ClusterStatus) error {
  res := x.db.Model(&model.Cluster{}).
    Where("ID = ?", id).
    Update("Status", status)

  if res.Error != nil || res.RowsAffected == 0 {
    return fmt.Errorf("nothing updated in cluster with id %s", id.String())
  }

  return nil
}

func ConvertToPbCluster(cluster model.Cluster) *pb.Cluster {
  tempConf := pb.ClusterConf{
    MasterFlavor: cluster.MasterFlavor,
    MasterReplicas: cluster.MasterReplicas,
    MasterRootSize: cluster.MasterRootSize,
    WorkerFlavor: cluster.WorkerFlavor,
    WorkerReplicas: cluster.WorkerReplicas,
    WorkerRootSize: cluster.WorkerRootSize,
    K8SVersion: cluster.K8sVersion,
  }

  return &pb.Cluster{
    Id:         cluster.ID.String(),
    Name:       cluster.Name,
    CreatedAt:  timestamppb.New(cluster.CreatedAt),
    UpdatedAt:  timestamppb.New(cluster.CreatedAt),
    Status:     cluster.Status,
    ContractId: cluster.ContractID.String(),
    Kubeconfig: cluster.Kubeconfig,
    Conf:       &tempConf,
  }
}
