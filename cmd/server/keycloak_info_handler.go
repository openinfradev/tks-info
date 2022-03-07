package main

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"github.com/google/uuid"

	"github.com/openinfradev/tks-common/pkg/log"
	"github.com/openinfradev/tks-info/pkg/keycloak_info"
	pb "github.com/openinfradev/tks-proto/tks_pb"
)

var (
	keycloakInfoAccessor *keycloak_info.KeycloakInfoAccessor
)

type KeycloakInfoServer struct {
	pb.UnimplementedKeycloakInfoServiceServer
}

func InitKeycloakInfoHandler(db *gorm.DB) {
	keycloakInfoAccessor = keycloak_info.New(db)
}

func (s *KeycloakInfoServer) CreateKeycloakInfo(ctx context.Context, in *pb.CreateKeycloakInfoRequest) (*pb.IDResponse, error) {
	log.Info("Request 'CreateKeycloakInfo' ")

	clusterId, err := uuid.Parse(in.GetClusterId())
	if err != nil {
		return &pb.IDResponse{
			Code: pb.Code_INVALID_ARGUMENT,
			Error: &pb.Error{
				Msg: fmt.Sprintf("invalid cluster ID %s", in.GetClusterId()),
			},
		}, err
	}

	id, err := keycloakInfoAccessor.Create(clusterId, in.GetRealm(), in.GetClientId(), in.GetSecret(), in.GetPrivateKey())
	if err != nil {
		return &pb.IDResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}

	return &pb.IDResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
		Id:    id.String(),
	}, nil
}

func (s *KeycloakInfoServer) GetKeycloakInfoByClusterId(ctx context.Context, in *pb.IDRequest) (*pb.GetKeycloakInfoResponse, error) {
	log.Info("Request 'GetKeycloakInfoByClusterId' clusterId ", in.GetId())

	clusterId, err := uuid.Parse(in.GetId())
	if err != nil {
		return &pb.GetKeycloakInfoResponse{
			Code: pb.Code_INVALID_ARGUMENT,
			Error: &pb.Error{
				Msg: fmt.Sprintf("Invalid cluster ID %s", in.GetId()),
			},
		}, err
	}

	keycloakInfos, err := keycloakInfoAccessor.GetKeycloakInfos(clusterId)
	if err != nil {
		return &pb.GetKeycloakInfoResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: fmt.Sprintf("Failed to get keycloak infos. err : %s", err.Error()),
			},
		}, err
	}

	return &pb.GetKeycloakInfoResponse{
		Code:          pb.Code_OK_UNSPECIFIED,
		Error:         nil,
		KeycloakInfos: keycloakInfos,
	}, nil
}

func (s *KeycloakInfoServer) UpdateKeycloakInfo(ctx context.Context, in *pb.IDRequest) (*pb.SimpleResponse, error) {
	log.Info("Request 'UpdateKeycloakInfo' ")
	log.Warn("Not Implemented gRPC API: 'UpdateKeycloakInfo'")
	return &pb.SimpleResponse{
		Code:  pb.Code_UNIMPLEMENTED,
		Error: nil,
	}, nil
}

func (s *KeycloakInfoServer) DeleteKeycloakInfo(ctx context.Context, in *pb.IDRequest) (*pb.SimpleResponse, error) {
	log.Info("Request 'DeleteKeycloakInfo' ")
	log.Warn("Not Implemented gRPC API: 'DeleteKeycloakInfo'")
	return &pb.SimpleResponse{
		Code:  pb.Code_UNIMPLEMENTED,
		Error: nil,
	}, nil
}
