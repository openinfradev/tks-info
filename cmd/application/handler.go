package application

import (
	"context"

	"github.com/sktelecom/tks-contract/pkg/log"
	pb "github.com/sktelecom/tks-proto/pbgo"
)

type Server struct {
	pb.UnimplementedAppInfoServiceServer
}

func (*Server) GetAppID(ctx context.Context, req *pb.GetAppRequest) (*pb.IDsResponse, error) {
	clusterId := req.GetClusterId()
	log.Info("clusterId: ", clusterId)

	ids := []string{"111", "222"}
	res := &pb.IDsResponse{
		Ids: ids,
	}

	log.Info("AppIds: ", ids)
	return res, nil
}

func (s *Server) GetAllApps(ctx context.Context, req *pb.IDRequest) (*pb.GetAppsResponse, error) {
	clusterId := req.GetId()
	log.Info("clusterId: ", clusterId)

	res := &pb.GetAppsResponse{}
	return res, nil
}

func (*Server) GetApps(ctx context.Context, req *pb.GetAppsRequest) (*pb.GetAppResponse, error) {
	appId := req.GetAppId()
	log.Info("appId: ", appId)

	res := &pb.GetAppResponse{}
	return res, nil
}

func (*Server) UpdateAppStatus(ctx context.Context, req *pb.UpdateAppStatusRequest) (*pb.SimpleResponse, error) {
	clusterId := req.GetClusterId()
	log.Info("clusterId: ", clusterId)

	appId := req.GetAppId()
	log.Info("appId: ", appId)

	res := &pb.SimpleResponse{}
	return res, nil
}

func (*Server) UpdateEndpoints(ctx context.Context, req *pb.UpdateEndpointsRequest) (*pb.SimpleResponse, error) {
	clusterId := req.GetClusterId()
	log.Info("clusterId: ", clusterId)

	appId := req.GetAppId()
	log.Info("appId: ", appId)

	res := &pb.SimpleResponse{}
	return res, nil
}
