package csp_info_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/openinfradev/tks-common/pkg/helper"
	"github.com/openinfradev/tks-common/pkg/log"

	"github.com/openinfradev/tks-info/pkg/csp_info"
	"github.com/openinfradev/tks-info/pkg/csp_info/model"
)

var (
	cspId           uuid.UUID
	contractId      uuid.UUID
	cspInfoAccessor *csp_info.CspInfoAccessor
)

var (
	testDBHost string
	testDBPort string
	err error
)

func init() {
	contractId = uuid.New()

	log.Disable()
}

func getAccessor() (*csp_info.CspInfoAccessor, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
		testDBHost, "postgres", "password", "tks", testDBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	if err := db.AutoMigrate(&model.CSPInfo{}); err != nil {
		return nil, err
	}

	return csp_info.New(db), nil
}

func TestMain(m *testing.M) {
	pool, resource, err := helper.CreatePostgres()
	if err != nil {
		fmt.Printf("Could not create postgres: %s", err)
		os.Exit(-1)
	}
	testDBHost, testDBPort = helper.GetHostAndPort(resource)
	cspInfoAccessor, _ = getAccessor()

	code := m.Run()

	if err := helper.RemovePostgres(pool, resource); err != nil {
		fmt.Printf("Could not remove postgres: %s", err)
		os.Exit(-1)
	}
	os.Exit(code)
}

func TestCreateCSPInfo(t *testing.T) {
	cspId, err = cspInfoAccessor.Create(contractId, "dummy", "DUMMYAUTH", 0)
	if err != nil {
		t.Errorf("An error occurred while creating new cspInfo. Err: %s", err)
	}
}

func TestGetCSPIDsByContractID(t *testing.T) {
	ids, err := cspInfoAccessor.GetCSPIDsByContractID(contractId)
	if err != nil {
		t.Errorf("An error occurred while getting CSP IDs. Err: %s", err)
	}

	for idx, id := range ids {
		t.Logf("%d) CSP id: %s", idx+1, id)
	}
}

func TestUpdateCSPAuth(t *testing.T) {
	err := cspInfoAccessor.UpdateCSPAuth(cspId, "NEWDUMMYAUTH")
	if err != nil {
		t.Errorf("An error occurred while updating CSP auth. Err: %s", err)
	}
}
