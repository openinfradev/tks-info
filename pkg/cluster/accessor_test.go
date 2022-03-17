package cluster_test

import (
  "fmt"
  "os"
  "testing"

  "gorm.io/driver/postgres"
  "gorm.io/gorm"

  "github.com/google/uuid"
  "github.com/stretchr/testify/assert"

  "github.com/openinfradev/tks-common/pkg/helper"
  "github.com/openinfradev/tks-common/pkg/log"

  "github.com/openinfradev/tks-info/pkg/cluster"
  "github.com/openinfradev/tks-info/pkg/cluster/model"
  pb "github.com/openinfradev/tks-proto/tks_pb"
)

var (
  clusterId       uuid.UUID
  cspId           uuid.UUID
  contractId      uuid.UUID
  clusterAccessor *cluster.ClusterAccessor
  clusterName     string
)

var (
  testDBHost string
  testDBPort string
  err error
)

func init() {
  contractId = uuid.New()
  cspId = uuid.New()
  clusterName = "testCluster"

  log.Disable()
}

func getAccessor() (*cluster.ClusterAccessor, error) {
  dsn := fmt.Sprintf(
    "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
    testDBHost, "postgres", "password", "tks", testDBPort)
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil {
    return nil, err
  }

  db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

  if err := db.AutoMigrate(&model.Cluster{}); err != nil {
    return nil, err
  }

  return cluster.New(db), nil
}

func TestMain(m *testing.M) {
  pool, resource, err := helper.CreatePostgres()
  if err != nil {
    fmt.Printf("Could not create postgres: %s", err)
    os.Exit(-1)
  }
  testDBHost, testDBPort = helper.GetHostAndPort(resource)
  clusterAccessor, _ = getAccessor()

  code := m.Run()

  if err := helper.RemovePostgres(pool, resource); err != nil {
    fmt.Printf("Could not remove postgres: %s", err)
    os.Exit(-1)
  }
  os.Exit(code)
}

// Create creates new cluster with contract ID, csp ID, name.
func TestCreateClusterInfo(t *testing.T) {
  dummyConf := pb.ClusterConf{
    SshKeyName:     "tks-seoul",
    Region:         "ap-northeast-2",
    NumOfAz:        3,
    MachineType:    "t3.large",
    MinSizePerAz:   1,
    MaxSizePerAz:   5,
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
