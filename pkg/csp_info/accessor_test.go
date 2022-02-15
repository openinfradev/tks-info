package csp_info_test

import (
	uuid "github.com/google/uuid"
	"github.com/openinfradev/tks-info/pkg/csp_info"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

var (
	cspId           uuid.UUID
	contractId      uuid.UUID
	cspInfoAccessor *csp_info.CspInfoAccessor
	err             error
)

func init() {
	dsn := "host=localhost user=postgres password=password dbname=tks port=5432 sslmode=disable TimeZone=Asia/Seoul"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	cspInfoAccessor = csp_info.New(db)

	// Create contract in advance for test cases
	contractId = uuid.New()
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
