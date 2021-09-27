package keycloak_info

import (
  "github.com/google/uuid"
  "gorm.io/gorm"

  model "github.com/openinfradev/tks-info/pkg/keycloak_info/model"
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

/*
func (x *KeycloakInfoAccessor) GetKeycloakInfos(clusterId uuid.UUID) ([]KeycloakInfo, error) {
  var keycloakInfos []model.KeycloakInfo

  res := x.db.Select("id").Find(&keycloakInfos, "cluster_id = ?", clusterId)

  if res.RowsAffected == 0 || res.Error != nil {
    return []string{}, fmt.Errorf("Could not find KeycloakInfo with cluster ID: %s", clusterId)
  }

  var keycloakInfos []KeycloakInfo
  for _, item := range keycloakInfos {
    idArr = append(keycloakInfos, item.KeycloakInfo)
  }
  return idArr, nil
}
*/
