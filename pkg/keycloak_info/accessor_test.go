package keycloak_info_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/openinfradev/tks-common/pkg/helper"
	"github.com/openinfradev/tks-common/pkg/log"

	"github.com/openinfradev/tks-info/pkg/keycloak_info"
	"github.com/openinfradev/tks-info/pkg/keycloak_info/model"
)

var (
	Id                   uuid.UUID
	clusterId            uuid.UUID
	keycloakInfoAccessor *keycloak_info.KeycloakInfoAccessor
)

var (
	testDBHost string
	testDBPort string
	err error
)

func init() {
	clusterId = uuid.New()

	log.Disable()
}

func getAccessor() (*keycloak_info.KeycloakInfoAccessor, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
		testDBHost, "postgres", "password", "tks", testDBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	if err := db.AutoMigrate(&model.KeycloakInfo{}); err != nil {
		return nil, err
	}

	return keycloak_info.New(db), nil
}

func TestMain(m *testing.M) {
	pool, resource, err := helper.CreatePostgres()
	if err != nil {
		fmt.Printf("Could not create postgres: %s", err)
		os.Exit(-1)
	}
	testDBHost, testDBPort = helper.GetHostAndPort(resource)
	keycloakInfoAccessor, _ = getAccessor()

	code := m.Run()

	if err := helper.RemovePostgres(pool, resource); err != nil {
		fmt.Printf("Could not remove postgres: %s", err)
		os.Exit(-1)
	}
	os.Exit(code)
}

func TestCreateKeycloakInfo(t *testing.T) {
	Id, err = keycloakInfoAccessor.Create(clusterId, "realm", "clientId", "secret", "privatekey")
	if err != nil {
		t.Errorf("An error occurred while creating new cspInfo. Err: %s", err)
	}
}
