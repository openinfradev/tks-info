package main

import (
  "context"
  "fmt"
  "gorm.io/gorm"

  "github.com/google/uuid"
  "github.com/openinfradev/tks-contract/pkg/log"
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
  log.Debug("Request 'CreateKeycloakInfo' ")

  clusterId, err := uuid.Parse(in.GetClusterId())
  if err != nil {
    res := pb.IDResponse{
      Code: pb.Code_INVALID_ARGUMENT,
      Error: &pb.Error{
        Msg: fmt.Sprintf("invalid cluster ID %s", in.GetClusterId()),
      },
    }
    return &res, err
  }

  id, err := keycloakInfoAccessor.Create(clusterId, in.GetRealm(), in.GetClientId(), in.GetSecret(), in.GetPrivateKey() )
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
  //log.Debug("Request 'GetKeycloakInfoByClusterId' clusterId ", in.GetClusterId() )

  /*
  if in.GetClusterId() == "" {
    return &pb.GetKeycloakInfoResponse {
      Code:  pb.Code_INVALID_ARGUMENT,
      Error: nil,
    }, nil
  }

  keycloakInfos, err := keycloakInfoAccessor.GetKeycloakInfos(in.GetClusterId())
  if err != nil {
    return &pb.GetKeycloakInfoResponse {
      Code: pb.Code_INTERNAL,
      Error: &pb.Error{
        Msg: fmt.Sprintf("failed to get keycloak infos. clusterId %s", in.GetClusterId() ),
      },
    }, err
  }

  return &pb.GetKeycloakInfoResponse{
    Code:  pb.Code_OK_UNSPECIFIED,
    Error: nil,
    KeycloakInfo: keyclockInfos,
  }, nil
  */
  return &pb.GetKeycloakInfoResponse{
    Code:  pb.Code_OK_UNSPECIFIED,
    Error: nil,
  }, nil
}

func (s *KeycloakInfoServer) UpdateKeycloakInfo(ctx context.Context, in *pb.IDRequest) (*pb.SimpleResponse, error) {
  log.Debug("Request 'UpdateKeycloakInfo' ")
  log.Warn("Not Implemented gRPC API: 'UpdateKeycloakInfo'")
  return &pb.SimpleResponse{
    Code:  pb.Code_UNIMPLEMENTED,
    Error: nil,
  }, nil
}

func (s *KeycloakInfoServer) DeleteKeycloakInfo(ctx context.Context, in *pb.IDRequest) (*pb.SimpleResponse, error) {
  log.Debug("Request 'DeleteKeycloakInfo' ")
  log.Warn("Not Implemented gRPC API: 'DeleteKeycloakInfo'")
  return &pb.SimpleResponse{
    Code:  pb.Code_UNIMPLEMENTED,
    Error: nil,
  }, nil
}

