package application_test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/openinfradev/tks-common/pkg/helper"
	"github.com/openinfradev/tks-common/pkg/log"

	"github.com/openinfradev/tks-info/pkg/application"
	"github.com/openinfradev/tks-info/pkg/application/model"
	pb "github.com/openinfradev/tks-proto/tks_pb"
)

var (
	clusterID  uuid.UUID
	appGroupID uuid.UUID
	appName    string
	accessor   *application.Accessor
)

var (
	testDBHost string
	testDBPort string
)

func init() {
	clusterID = uuid.New()

	log.Disable()
}

func getAccessor() (*application.Accessor, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
		testDBHost, "postgres", "password", "tks", testDBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	if err := db.AutoMigrate(&model.Application{}); err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&model.ApplicationGroup{}); err != nil {
		return nil, err
	}

	return application.New(db), nil
}

func TestMain(m *testing.M) {
	pool, resource, err := helper.CreatePostgres()
	if err != nil {
		fmt.Printf("Could not create postgres: %s", err)
		os.Exit(-1)
	}
	testDBHost, testDBPort = helper.GetHostAndPort(resource)
	accessor, _ = getAccessor()

	code := m.Run()

	if err := helper.RemovePostgres(pool, resource); err != nil {
		fmt.Printf("Could not remove postgres: %s", err)
		os.Exit(-1)
	}
	os.Exit(code)
}

func TestCreateApplicationGroup(t *testing.T) {
	var err error
	appName = getRandomString("gotest")
	appGroup1 := pb.AppGroup{
		AppGroupName:  appName,
		Type:          pb.AppGroupType_LMA,
		Subtype:       pb.AppGroupSubtype_LOKI,
		Status:        pb.AppGroupStatus_APP_GROUP_INSTALLING,
		ExternalLabel: "test_env",
	}
	t.Logf("new cluster ID %s", clusterID)
	appGroupID, err = accessor.Create(clusterID, &appGroup1)
	if err != nil {
		t.Errorf("an error was unexpected while creating new application group: %s", err)
	}
	appGroup2 := pb.AppGroup{
		AppGroupName:  appName,
		Type:          pb.AppGroupType_SERVICE_MESH,
		Status:        pb.AppGroupStatus_APP_GROUP_INSTALLING,
		ExternalLabel: "",
	}

	appGroupID2, err := accessor.Create(clusterID, &appGroup2)
	if err != nil {
		t.Errorf("an error was unexpected while creating new application group: %s", err)
	}
	t.Logf("new app group id: %s, %s", appGroupID, appGroupID2)
}
func TestGetAppGroupsByClusterID(t *testing.T) {
	appGroups, err := accessor.GetAppGroupsByClusterID(clusterID, 0, 10)
	if err != nil {
		t.Errorf("an error was unexpected while creating new application group: %s", err)
	}
	for idx, appGroup := range appGroups {
		t.Logf("%d) matching app group id: %s, type: %d", idx+1, appGroup.AppGroupId, appGroup.Type)
	}
}
func TestGetAppGroups(t *testing.T) {
	appGroups, err := accessor.GetAppGroups(appName, pb.AppGroupType_APP_TYPE_UNSPECIFIED)
	if err != nil {
		t.Errorf("an error was unexpected while creating new application group: %s", err)
	}
	for idx := range appGroups {
		t.Logf("%d) matching app group id: %s, type: %d", idx+1, appGroups[idx].AppGroupId, appGroups[idx].Type)
	}
}
func TestGetAppGroup(t *testing.T) {
	appGroup, err := accessor.GetAppGroup(appGroupID)
	if err != nil {
		t.Errorf("an error was unexpected while get application group: %s", err)
	}
	t.Logf("matching app group name: %s", appGroup.AppGroupName)
}
func TestUpdateAppGroupStatus(t *testing.T) {
	if err := accessor.UpdateAppGroupStatus(appGroupID, pb.AppGroupStatus_APP_GROUP_RUNNING); err != nil {
		t.Errorf("an error was unexpected while update application group: %s", err)
	}

	appGroup, err := accessor.GetAppGroup(appGroupID)
	if err != nil {
		return
	}
	if appGroup.Status != pb.AppGroupStatus_APP_GROUP_RUNNING {
		t.Errorf("app group status was not updated, status: %d", appGroup.Status)
	}
}

func TestUpdateApp(t *testing.T) {
	if err := accessor.UpdateApp(appGroupID, pb.AppType_PROMETHEUS,
		"http://localhost:9090", "{\"metadata\":\"no_data\"}"); err != nil {
		t.Errorf("an error was unexpected while update prometheus: %s", err)
	}
	if err := accessor.UpdateApp(appGroupID, pb.AppType_KIALI,
		"http://localhost:20001", "{\"metadata\":\"no_data\"}"); err != nil {
		t.Errorf("an error was unexpected while update kiali: %s", err)
	}
}

func TestGetAppsByAppGroupID(t *testing.T) {
	apps, err := accessor.GetAppsByAppGroupID(appGroupID)
	if err != nil {
		t.Errorf("an error was unexpected while get applications: %s", err)
	}

	t.Logf(">>> Get apps by app_group_id result:")
	for _, app := range apps {
		t.Logf("app id: %s, app type: %d", app.AppId, app.Type)
	}
}
func TestGetApps(t *testing.T) {
	apps, err := accessor.GetApps(appGroupID, pb.AppType_PROMETHEUS)
	if err != nil {
		t.Errorf("an error was unexpected while get applications: %s", err)
	}

	t.Logf(">>> Get apps by app_group_id result:")
	for _, app := range apps {
		t.Logf("app id: %s, app type: %d", app.AppId, app.Type)
	}
}

func TestDeleteAppGroup(t *testing.T) {
	if err := accessor.DeleteAppGroup(appGroupID); err != nil {
		t.Errorf("an error was unexpected while delete application group: %s", err)
	}

	_, err := accessor.GetAppGroup(appGroupID)
	expectedErr := fmt.Errorf("could not find application group for app_group_id %s", appGroupID)
	if err.Error() == expectedErr.Error() {
		return
	}
}

func getRandomString(prefix string) string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return fmt.Sprintf("%s-%d", prefix, r.Int31n(1000000000))
}
