package keycloak_info_test

import (
	uuid "github.com/google/uuid"
	"github.com/openinfradev/tks-info/pkg/keycloak_info"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

var (
	Id                   uuid.UUID
	clusterId            uuid.UUID
	keycloakInfoAccessor *keycloak_info.KeycloakInfoAccessor
	err                  error
)

func init() {
	dsn := "host=localhost user=postgres password=password dbname=tks port=5432 sslmode=disable TimeZone=Asia/Seoul"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	keycloakInfoAccessor = keycloak_info.New(db)

	// Create contract in advance for test cases
	clusterId = uuid.New()
}

func TestCreateKeycloakInfo(t *testing.T) {
	Id, err = keycloakInfoAccessor.Create(clusterId, "realm", "clientId", "secret", "privatekey")
	if err != nil {
		t.Errorf("An error occurred while creating new cspInfo. Err: %s", err)
	}
}
