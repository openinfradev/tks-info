package application

import (
	"time"

	pb "github.com/sktelecom/tks-proto/pbgo"
)

type PostgreAccessor struct{}

func New() *PostgreAccessor {
	return &PostgreAccessor{}
}

func getDummyAppInfo() *AppInfo {
	return &AppInfo{
		AppID:     "111",
		AppName:   "my_service_mesh",
		AppType:   pb.AppType_SERVICE_MESH,
		Owner:     "ccc",
		AppStatus: pb.AppStatus_APP_RUNNING,
		EndPoints: []*pb.Endpoint{
			{
				Type: pb.EpType_KIALI,
				Url:  "kiali.istio-system.svc.cluster.k1",
			},
			{
				Type: pb.EpType_JAEGER,
				Url:  "jaeger.istio-system.svc.cluster.k1",
			},
		},
		ExternalLabel: "service_mesh",
		CreatedTs:     time.Now(),
		LastUpdatedTs: time.Now(),
	}
}

func getDummyAppInfos() []*AppInfo {
	return []*AppInfo{
		{
			AppID:     "111",
			AppName:   "my_service_mesh",
			AppType:   pb.AppType_SERVICE_MESH,
			Owner:     "ccc",
			AppStatus: pb.AppStatus_APP_RUNNING,
			EndPoints: []*pb.Endpoint{
				{
					Type: pb.EpType_KIALI,
					Url:  "kiali.istio-system.svc.cluster.k1",
				},
				{
					Type: pb.EpType_JAEGER,
					Url:  "jaeger.istio-system.svc.cluster.k1",
				},
			},
			ExternalLabel: "service_mesh",
			CreatedTs:     time.Now(),
			LastUpdatedTs: time.Now(),
		},
		{
			AppID:     "222",
			AppName:   "my_lma",
			AppType:   pb.AppType_LMA,
			Owner:     "ccc",
			AppStatus: pb.AppStatus_APP_RUNNING,
			EndPoints: []*pb.Endpoint{
				{
					Type: pb.EpType_KIALI,
					Url:  "kiali.istio-system.svc.cluster.k2",
				},
				{
					Type: pb.EpType_JAEGER,
					Url:  "jaeger.istio-system.svc.cluster.k2",
				},
			},
			ExternalLabel: "service_mesh",
			CreatedTs:     time.Now(),
			LastUpdatedTs: time.Now(),
		},
	}
}

func (p *PostgreAccessor) AddApp(appInfo *AppInfo) (ID, error) {
	var id ID = "111"
	return id, nil
}

func (p *PostgreAccessor) DeleteApp(clusterID ID, appID ID) error {
	return nil
}

func (p *PostgreAccessor) GetAppIDs(clusterID ID) ([]ID, error) {
	ids := []ID{"111", "222"}
	return ids, nil
}

func (p *PostgreAccessor) GetAllAppsByClusterID(clusterID ID) ([]*AppInfo, error) {
	appInfos := getDummyAppInfos()
	return appInfos, nil
}

func (p *PostgreAccessor) GetAppsByName(clusterID ID, appName string) ([]*AppInfo, error) {
	appInfos := getDummyAppInfos()
	return appInfos, nil
}

func (p *PostgreAccessor) GetAppsByType(clusterID ID, appType pb.AppType) ([]*AppInfo, error) {
	appInfos := getDummyAppInfos()
	return appInfos, nil
}

func (p *PostgreAccessor) GetApp(clusterID ID, appID ID) (*AppInfo, error) {
	appInfo := getDummyAppInfo()
	return appInfo, nil
}

func (p *PostgreAccessor) UpdateApp(appInfo *AppInfo) error {
	return nil
}

func (p *PostgreAccessor) UpdateAppStatus(appID ID, appStatus pb.AppStatus) error {
	return nil
}

func (p *PostgreAccessor) UpdateEndpoints(appID ID, appEndpoints []*pb.Endpoint) error {
	return nil
}
