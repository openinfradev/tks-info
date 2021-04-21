package application

import (
	"context"
	"errors"

	"github.com/sktelecom/tks-contract/pkg/log"
	app "github.com/sktelecom/tks-info/pkg/application"
	pb "github.com/sktelecom/tks-proto/pbgo"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

var acc app.Accessor = app.New()

type Server struct {
	pb.UnimplementedAppInfoServiceServer
}

func (s *Server) AddApp(ctx context.Context, req *pb.AddAppRequest) (*pb.IDResponse, error) {
	clusterId := req.GetClusterId()
	serviceApp := req.GetServiceApp()
	log.Info("clusterId: ", clusterId)
	log.Info("serivceApp: ", serviceApp)

	id, err := acc.AddApp(toAppInfo(serviceApp))
	if err != nil {
		return &pb.IDResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}

	res := &pb.IDResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
		Id:    string(id),
	}
	return res, nil
}

func (s *Server) DeleteApp(ctx context.Context, req *pb.DeleteAppRequest) (*pb.SimpleResponse, error) {
	clusterId := req.GetClusterId()
	appId := req.GetAppId()
	log.Info("clusterId: ", clusterId)
	log.Info("appId: ", appId)

	err := acc.DeleteApp(app.ID(clusterId), app.ID(appId))
	if err != nil {
		return &pb.SimpleResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}

	res := &pb.SimpleResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
	}
	return res, nil
}

func (s *Server) GetAppIDs(ctx context.Context, req *pb.IDRequest) (*pb.IDsResponse, error) {
	clusterId := req.GetId()
	log.Info("clusterId: ", clusterId)

	ids, err := acc.GetAppIDs(app.ID(clusterId))
	if err != nil {
		return &pb.IDsResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}

	res := &pb.IDsResponse{
		Ids:   toStrings(ids),
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
	}

	log.Info("AppIds: ", ids)
	return res, nil
}

func (s *Server) GetAllAppsByClusterID(ctx context.Context, req *pb.IDRequest) (*pb.GetAppsResponse, error) {
	clusterId := req.GetId()
	log.Info("clusterId: ", clusterId)

	apps, err := acc.GetAllAppsByClusterID(app.ID(clusterId))
	if err != nil {
		return &pb.GetAppsResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
			Apps: nil,
		}, err
	}
	if apps == nil {
		return &pb.GetAppsResponse{
			Code: pb.Code_NOT_FOUND,
			Error: &pb.Error{
				Msg: errors.New("Data Not Found").Error(),
			},
			Apps: nil,
		}, nil
	}

	res := &pb.GetAppsResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
		Apps:  toServiceApps(apps),
	}
	return res, nil
}

func (s *Server) GetAppsByName(ctx context.Context, req *pb.GetAppsRequest) (*pb.GetAppsResponse, error) {
	clusterId := req.GetClusterId()
	appName := req.GetAppName()
	log.Info("clusterId: ", clusterId)
	log.Info("appName: ", appName)

	apps, err := acc.GetAppsByName(app.ID(clusterId), appName)
	if err != nil {
		return &pb.GetAppsResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
			Apps: nil,
		}, err
	}
	if apps == nil {
		return &pb.GetAppsResponse{
			Code: pb.Code_NOT_FOUND,
			Error: &pb.Error{
				Msg: errors.New("Data Not Found").Error(),
			},
			Apps: nil,
		}, nil
	}

	res := &pb.GetAppsResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
		Apps:  toServiceApps(apps),
	}
	return res, nil
}

func (s *Server) GetAppsByType(ctx context.Context, req *pb.GetAppsRequest) (*pb.GetAppsResponse, error) {
	clusterId := req.GetClusterId()
	appType := req.GetType()
	log.Info("clusterId: ", clusterId)
	log.Info("appType: ", appType)

	apps, err := acc.GetAppsByType(app.ID(clusterId), appType)
	if err != nil {
		return &pb.GetAppsResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
			Apps: nil,
		}, err
	}
	if apps == nil {
		return &pb.GetAppsResponse{
			Code: pb.Code_NOT_FOUND,
			Error: &pb.Error{
				Msg: errors.New("Data Not Found").Error(),
			},
			Apps: nil,
		}, nil
	}

	res := &pb.GetAppsResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
		Apps:  toServiceApps(apps),
	}
	return res, nil

}

func (*Server) GetApp(ctx context.Context, req *pb.GetAppRequest) (*pb.GetAppResponse, error) {
	clusterId := req.GetClusterId()
	appId := req.GetAppId()
	log.Info("clusterId: ", clusterId)
	log.Info("appId: ", appId)

	app, err := acc.GetApp(app.ID(clusterId), app.ID(appId))
	if err != nil {
		return &pb.GetAppResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
			App: nil,
		}, err
	}
	if app == nil {
		return &pb.GetAppResponse{
			Code: pb.Code_NOT_FOUND,
			Error: &pb.Error{
				Msg: errors.New("Data Not Found").Error(),
			},
			App: nil,
		}, nil
	}

	res := &pb.GetAppResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
		App:   toServiceApp(app),
	}
	return res, nil
}

func (s *Server) UpdateApp(ctx context.Context, req *pb.UpdateAppRequest) (*pb.SimpleResponse, error) {
	clusterId := req.GetClusterId()
	serviceApp := req.GetServiceApp()
	log.Info("clusterId: ", clusterId)
	log.Info("appId: ", serviceApp)

	err := acc.UpdateApp(toAppInfo(serviceApp))
	if err != nil {
		return &pb.SimpleResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}

	res := &pb.SimpleResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
	}
	return res, nil
}

func (s *Server) UpdateAppStatus(ctx context.Context, req *pb.UpdateAppStatusRequest) (*pb.SimpleResponse, error) {
	clusterId := req.GetClusterId()
	appId := req.GetAppId()
	appStatus := pb.AppStatus_APP_RUNNING
	log.Info("clusterId: ", clusterId)
	log.Info("appId: ", appId)
	log.Info("appStatus: ", appStatus)

	err := acc.UpdateAppStatus(app.ID(appId), appStatus)
	if err != nil {
		return &pb.SimpleResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}

	res := &pb.SimpleResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
	}
	return res, nil
}

func (s *Server) UpdateEndpoints(ctx context.Context, req *pb.UpdateEndpointsRequest) (*pb.SimpleResponse, error) {
	clusterId := req.GetClusterId()
	appId := req.GetAppId()
	endPoints := req.GetEndpoints()
	log.Info("clusterId: ", clusterId)
	log.Info("appId: ", appId)
	log.Info("endPoints: ", endPoints)

	err := acc.UpdateEndpoints(app.ID(appId), endPoints)
	if err != nil {
		return &pb.SimpleResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}

	res := &pb.SimpleResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
	}
	return res, nil
}

func toStrings(ids []app.ID) []string {
	s := make([]string, len(ids))
	for i, v := range ids {
		s[i] = string(v)
	}
	return s
}

func toAppInfo(p *pb.ServiceApp) *app.AppInfo {
	res := &app.AppInfo{
		AppID:         app.ID(p.GetAppId()),
		AppName:       p.GetAppName(),
		AppType:       p.GetType(),
		Owner:         app.ID(p.GetOwner()),
		AppStatus:     p.GetStatus(),
		EndPoints:     p.GetEndpoints(),
		ExternalLabel: p.GetExternalLabel(),
		CreatedTs:     p.GetCreatedTs().AsTime(),
		LastUpdatedTs: p.GetLastUpdatedTs().AsTime(),
	}
	return res
}

func toServiceApps(apps []*app.AppInfo) []*pb.ServiceApp {
	res := make([]*pb.ServiceApp, len(apps))
	for i, v := range apps {
		res[i] = toServiceApp(v)
	}
	return res
}

func toServiceApp(a *app.AppInfo) *pb.ServiceApp {
	res := &pb.ServiceApp{
		AppId:         string(a.AppID),
		AppName:       a.AppName,
		Type:          a.AppType,
		Owner:         string(a.Owner),
		Status:        a.AppStatus,
		Endpoints:     a.EndPoints,
		ExternalLabel: a.ExternalLabel,
		CreatedTs:     timestamppb.New(a.CreatedTs),
		LastUpdatedTs: timestamppb.New(a.LastUpdatedTs),
	}
	return res
}
