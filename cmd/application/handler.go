package application

import (
	"context"

	"github.com/sktelecom/tks-contract/pkg/log"
	pb "github.com/sktelecom/tks-proto/pbgo"
)

type Server struct {
	pb.UnimplementedAppInfoServiceServer
}

func (s *Server) GetAppIDs(ctx context.Context, req *pb.IDRequest) (*pb.IDsResponse, error) {
	clusterId := req.GetId()
	log.Info("clusterId: ", clusterId)

	ids := []string{"111", "222"}
	res := &pb.IDsResponse{
		Ids: ids,
	}

	log.Info("AppIds: ", ids)
	return res, nil
}

func (s *Server) GetAllAppsByClusterID(ctx context.Context, req *pb.IDRequest) (*pb.GetAppsResponse, error) {
	clusterId := req.GetId()
	log.Info("clusterId: ", clusterId)

	res := &pb.GetAppsResponse{}
	return res, nil
}

func (s *Server) GetAppsByName(ctx context.Context, req *pb.GetAppsRequest) (*pb.GetAppsResponse, error) {
	appName := req.GetAppName()
	log.Info("app name: ", appName)

	res := &pb.GetAppsResponse{}
	return res, nil
}

func (s *Server) GetAppsByType(ctx context.Context, req *pb.GetAppsRequest) (*pb.GetAppsResponse, error) {
	appType := req.GetType()
	log.Info("app type: ", appType)

	res := &pb.GetAppsResponse{}
	return res, nil
}

func (*Server) GetApp(ctx context.Context, req *pb.GetAppRequest) (*pb.GetAppResponse, error) {
	appId := req.GetAppId()
	log.Info("appId: ", appId)

	res := &pb.GetAppResponse{}
	return res, nil
}

func (s *Server) UpdateAppStatus(ctx context.Context, req *pb.UpdateAppStatusRequest) (*pb.SimpleResponse, error) {
	clusterId := req.GetClusterId()
	log.Info("clusterId: ", clusterId)

	appId := req.GetAppId()
	log.Info("appId: ", appId)

	res := &pb.SimpleResponse{}
	return res, nil
}

func (s *Server) UpdateEndpoints(ctx context.Context, req *pb.UpdateEndpointsRequest) (*pb.SimpleResponse, error) {
	clusterId := req.GetClusterId()
	log.Info("clusterId: ", clusterId)

	appId := req.GetAppId()
	log.Info("appId: ", appId)

	res := &pb.SimpleResponse{}
	return res, nil
}
