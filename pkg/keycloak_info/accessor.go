package keycloak_info

import (
  "github.com/google/uuid"
  "gorm.io/gorm"
  "fmt"

  model "github.com/openinfradev/tks-info/pkg/keycloak_info/model"
  pb "github.com/openinfradev/tks-proto/tks_pb"
)

type KeycloakInfoAccessor struct {
  db *gorm.DB
}

func New(db *gorm.DB) *KeycloakInfoAccessor {
  return &KeycloakInfoAccessor{
    db: db,
  }
}

func (x *KeycloakInfoAccessor) Create(clusterId uuid.UUID, realm string, clientId string, secret string, privateKey string ) (uuid.UUID, error) {
  keycloackInfo := model.KeycloakInfo{ClusterId: clusterId, Realm: realm, ClientId: clientId, Secret: secret, PrivateKey: privateKey}

  res := x.db.Create(&keycloackInfo)
  if res.Error != nil {
    nilId, _ := uuid.Parse("")
    return nilId, res.Error
  }

  return keycloackInfo.Id, nil
}

func (x *KeycloakInfoAccessor) GetKeycloakInfos(clusterId uuid.UUID) ([]*pb.KeycloakInfo, error) {
  var keycloakInfos []model.KeycloakInfo

  res := x.db.Find("cluster_id").Find(&keycloakInfos, "cluster_id = ?", clusterId)

  if res.RowsAffected == 0 || res.Error != nil {
    return []*pb.KeycloakInfo{}, fmt.Errorf("Could not find KeycloakInfo with cluster ID: %s", clusterId)
  }

  pbKeycloakInfos := []*pb.KeycloakInfo{}
  for _, item := range keycloakInfos {
    pbKeycloakInfos = append(pbKeycloakInfos, ConvertToPbKeycloakInfo(item))
  }
  return pbKeycloakInfos, nil
}

func ConvertToPbKeycloakInfo(keycloakInfo model.KeycloakInfo) *pb.KeycloakInfo {
  return &pb.KeycloakInfo{
    ClusterId: keycloakInfo.ClusterId.String(),
    Realm:  keycloakInfo.Realm,
    ClientId:  keycloakInfo.ClientId,
    Secret:     keycloakInfo.Secret,
    PrivateKey: keycloakInfo.PrivateKey,
  }
}
