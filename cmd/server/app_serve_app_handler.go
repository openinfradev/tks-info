package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/openinfradev/tks-common/pkg/log"
	asa "github.com/openinfradev/tks-info/pkg/app_serve_app"
	pb "github.com/openinfradev/tks-proto/tks_pb"
)

var asaAccessor *asa.AsaAccessor

type AppServeAppServer struct {
	pb.UnimplementedAppServeAppServiceServer
}

func InitAppServeAppHandler(db *gorm.DB) {
	asaAccessor = asa.New(db)
}

func (s *AppServeAppServer) CreateAppServeApp(ctx context.Context, in *pb.CreateAppServeAppRequest) (*pb.IDResponse, error) {
	appServeApp := in.GetApp()
	contractID := appServeApp.GetContractId()
	log.Info("Handling request 'CreateAppServeApp' for contract id ", contractID)

	id, err := asaAccessor.Create(contractID, appServeApp)
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
		Id:    id.String(),
	}
	return res, nil
}

func (s *AppServeAppServer) UpdateAppServeAppStatus(ctx context.Context, in *pb.UpdateAppServeAppStatusRequest) (*pb.SimpleResponse, error) {
	appServeAppId, err := uuid.Parse(in.GetAppServeAppId())
	if err != nil {
		return &pb.SimpleResponse{
			Code: pb.Code_INVALID_ARGUMENT,
			Error: &pb.Error{
				Msg: fmt.Sprintf("invalid appServeApp ID %s", in.GetAppServeAppId()),
			},
		}, err
	}

	err = asaAccessor.UpdateStatus(appServeAppId, in.GetStatus(), in.GetOutput())
	if err != nil {
		return &pb.SimpleResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}
	return &pb.SimpleResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
	}, nil
}

func (s *AppServeAppServer) UpdateAppServeAppEndpoint(ctx context.Context, in *pb.UpdateAppServeAppEndpointRequest) (*pb.SimpleResponse, error) {
	appServeAppId, err := uuid.Parse(in.GetAppServeAppId())
	if err != nil {
		return &pb.SimpleResponse{
			Code: pb.Code_INVALID_ARGUMENT,
			Error: &pb.Error{
				Msg: fmt.Sprintf("invalid appServeApp ID %s", in.GetAppServeAppId()),
			},
		}, err
	}

	err = asaAccessor.UpdateEndpoint(appServeAppId, in.GetEndpoint())
	if err != nil {
		return &pb.SimpleResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}
	return &pb.SimpleResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
	}, nil
}

func (s *AppServeAppServer) GetAppServeApps(ctx context.Context, in *pb.GetAppServeAppsRequest) (*pb.GetAppServeAppsResponse, error) {
	contractId := in.GetContractId()
	log.Info("GetAppServeApps request for contractID: ", contractId)

	appServeApps, err := asaAccessor.GetAppServeApps(contractId)
	if err != nil {
		return &pb.GetAppServeAppsResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}

	return &pb.GetAppServeAppsResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
		Apps:  appServeApps,
	}, nil
}

func (s *AppServeAppServer) GetAppServeApp(ctx context.Context, in *pb.GetAppServeAppRequest) (*pb.GetAppServeAppResponse, error) {
	id, err := uuid.Parse(in.GetAppServeAppId())
	if err != nil {
		return &pb.GetAppServeAppResponse{
			Code: pb.Code_INVALID_ARGUMENT,
			Error: &pb.Error{
				Msg: fmt.Sprintf("Invalid appServeApp ID: %s", in.GetAppServeAppId()),
			},
		}, err
	}
	log.Info("Received GetAppServeApp request for ID: ", id)

	appServeApp, err := asaAccessor.GetAppServeApp(id, in.GetContractId())
	if err != nil {
		return &pb.GetAppServeAppResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}

	return &pb.GetAppServeAppResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
		App:   appServeApp,
	}, nil

}
