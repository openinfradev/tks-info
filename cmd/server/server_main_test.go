package main

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/openinfradev/tks-common/pkg/helper"
	"github.com/openinfradev/tks-common/pkg/log"

	modelApplication "github.com/openinfradev/tks-info/pkg/application/model"
	modelCluster "github.com/openinfradev/tks-info/pkg/cluster/model"
	modelCspInfoInfo "github.com/openinfradev/tks-info/pkg/csp_info/model"
	modelKeyCloackInfo "github.com/openinfradev/tks-info/pkg/keycloak_info/model"
)

var (
	err error
	db  *gorm.DB
)

func init() {
	log.Disable()
}

func TestMain(m *testing.M) {
	pool, resource, err := helper.CreatePostgres()
	if err != nil {
		fmt.Printf("Could not create postgres: %s", err)
		os.Exit(-1)
	}
	testDBHost, testDBPort := helper.GetHostAndPort(resource)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
		testDBHost, "postgres", "password", "tks", testDBPort)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		os.Exit(-1)
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	err = db.AutoMigrate(&modelApplication.Application{})
	err = db.AutoMigrate(&modelApplication.ApplicationGroup{})
	err = db.AutoMigrate(&modelCluster.Cluster{})
	err = db.AutoMigrate(&modelCspInfoInfo.CSPInfo{})
	err = db.AutoMigrate(&modelKeyCloackInfo.KeycloakInfo{})
	if err != nil {
		os.Exit(-1)
	}

	InitAppInfoHandler(db)
	InitKeycloakInfoHandler(db)
	InitClusterInfoHandler(db)
	InitCspInfoHandler(db)

	code := m.Run()

	_ = pool != nil

	if err := helper.RemovePostgres(pool, resource); err != nil {
		fmt.Printf("Could not remove postgres: %s", err)
		os.Exit(-1)
	}

	os.Exit(code)
}

// Helpers
func randomString(prefix string) string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return fmt.Sprintf("%s-%d", prefix, r.Int31n(1000000000))
}
