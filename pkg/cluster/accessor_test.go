package cluster_test

import (
  "testing"
  "github.com/stretchr/testify/assert"
  uuid "github.com/google/uuid"
  "gorm.io/gorm"
  "gorm.io/driver/postgres"
  "github.com/sktelecom/tks-info/pkg/cluster"
  pb "github.com/sktelecom/tks-proto/pbgo"
)

var (
  clusterId  uuid.UUID
  cspId  uuid.UUID
  contractId uuid.UUID
  clusterAccessor   *cluster.ClusterAccessor
  clusterName string
  err error
)

func init() {
  dsn := "host=localhost user=postgres password=password dbname=tks port=5432 sslmode=disable TimeZone=Asia/Seoul"
  db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  clusterAccessor = cluster.New(db)
  contractId = uuid.New()
  cspId = uuid.New()
  clusterName = "testCluster"
}

// Create creates new cluster with contract ID, csp ID, name.
func TestCreateClusterInfo(t *testing.T) {
  dummyConf := pb.ClusterConf{
    MasterFlavor: "tiny",
    MasterReplicas: 3,
    MasterRootSize: 50,
    WorkerFlavor: "medium",
    WorkerReplicas: 5,
    WorkerRootSize: 50,
    K8SVersion: "1.18.8",
    Region: "ap-southeast-2",
  }

  clusterId, err = clusterAccessor.CreateClusterInfo(contractId, cspId, clusterName, &dummyConf)
  if err != nil {
    t.Errorf("An error occurred while creating new clusterInfo. Err: %s", err)
  }

  t.Logf("Created clusterID: %s", clusterId.String())
}

func TestGetCluster(t *testing.T) {
  cluster, err := clusterAccessor.GetCluster(clusterId)
  if err != nil {
    t.Errorf("An error occurred while getting cluster info. Err: %s", err)
  }

  t.Logf("Retrieved clusterName: %s", cluster.Name)
}

func TestGetClustersByCspID(t *testing.T) {
  clusters, err := clusterAccessor.GetClustersByCspID(cspId)
  if err != nil {
    t.Errorf("An error occurred while getting clusterInfo by cspID. Err:  %s", err)
  }

  for _, cluster := range clusters {
    //t.Logf("%d) cluster id: %s, name: %d", idx+1, cluster.ID, cluster.Name)
    assert.Equal(t, "testCluster", cluster.Name, "Not expected cluster name!")
  }
}

func TestUpdateStatus(t *testing.T) {
  err := clusterAccessor.UpdateStatus(clusterId, pb.ClusterStatus_INSTALLING)
  if err != nil {
    t.Errorf("An error occurred while updating cluster status. Err: %s", err)
  }
}
