package application

import (
	"time"

	pb "github.com/sktelecom/tks-proto/pbgo"
)

type Accessor interface {
	AddApp(appInfo *AppInfo) (ID, error)
	DeleteApp(clusterID ID, appID ID) error
	GetAppIDs(clusterID ID) ([]ID, error)
	GetAllAppsByClusterID(clusterID ID) ([]*AppInfo, error)
	GetAppsByName(clusterID ID, appName string) ([]*AppInfo, error)
	GetAppsByType(clusterID ID, appType pb.AppType) ([]*AppInfo, error)
	GetApp(clusterID ID, appID ID) (*AppInfo, error)
	UpdateApp(appInfo *AppInfo) error
	UpdateAppStatus(appID ID, appStatus pb.AppStatus) error
	UpdateEndpoints(appID ID, appEndpoints []*pb.Endpoint) error
}

// ID is a global unique ID.
type ID string

// AppInfo represents an ServiceApp information which is deployed K8S clusters.
type AppInfo struct {
	AppID         ID             `json:"app_id"`
	AppName       string         `json:"auth"`
	AppType       pb.AppType     `json:"app_type"`
	Owner         ID             `json:"owner"`
	AppStatus     pb.AppStatus   `json:"app_status"`
	EndPoints     []*pb.Endpoint `json:"endpoints"`
	ExternalLabel string         `json:"external_label"`
	CreatedTs     time.Time      `json:"created_ts"`
	LastUpdatedTs time.Time      `json:"last_updated_ts"`
}
