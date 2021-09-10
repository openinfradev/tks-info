package keycloak_info_test

import (
  "testing"
  uuid "github.com/google/uuid"
  "gorm.io/gorm"
  "gorm.io/driver/postgres"
  "github.com/openinfradev/tks-info/pkg/keycloak_info"
)

var (
  Id  uuid.UUID
  clusterId uuid.UUID
  keycloakInfoAccessor   *keycloak_info.KeycloakInfoAccessor
  err error
)

func init() {
  dsn := "host=localhost user=postgres password=password dbname=tks port=5432 sslmode=disable TimeZone=Asia/Seoul"
  db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  keycloakInfoAccessor = keycloak_info.New(db)

  // Create contract in advance for test cases
  clsuterId = uuid.New()
}

func TestCreateKeycloakInfo(t *testing.T) {
  id, err = keycloakInfoAccessor.Create(clusterId, "dummy", "DUMMYAUTH")
  if err != nil {
    t.Errorf("An error occurred while creating new cspInfo. Err: %s", err)
  }
}



