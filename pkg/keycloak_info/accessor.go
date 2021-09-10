package keycloak_info

import (
  uuid "github.com/google/uuid"
  "gorm.io/gorm"

  model "github.com/openinfradev/tks-info/pkg/keycloak_info/model"
)

// Accessor accesses to keycloak info in-memory data.
type KeycloakInfoAccessor struct {
  db *gorm.DB
}

// NewKeycloakInfoAccessor returns new Accessor to access csp info.
func New(db *gorm.DB) *KeycloakInfoAccessor {
  return &KeycloakInfoAccessor{
    db: db,
  }
}

// Create creates new Keycloak info
func (x *KeycloakInfoAccessor) Create(clusterId uuid.UUID, realm string, clientId string, secret string, privateKey string ) (uuid.UUID, error) {
  keycloackInfo := model.KeycloakInfo{ClusterId: clusterId, Realm: realm, ClientId: clientId, Secret: secret, PrivateKey: privateKey}

  res := x.db.Create(&keycloackInfo)
  if res.Error != nil {
    nilId, _ := uuid.Parse("")
    return nilId, res.Error
  }

  return keycloackInfo.Id, nil
}

