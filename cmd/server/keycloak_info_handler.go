package main

import (
  "context"
  "gorm.io/gorm"

  "github.com/openinfradev/tks-contract/pkg/log"
  "github.com/openinfradev/tks-info/pkg/keycloak_info"
  pb "github.com/openinfradev/tks-proto/pbgo"
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
  log.Warn("Not Implemented gRPC API: 'CreateKeycloakInfo'")
  return &pb.IDResponse{
    Code:  pb.Code_UNIMPLEMENTED,
    Error: nil,
  }, nil
}

func (s *KeycloakInfoServer) GetKeycloakInfos(ctx context.Context, in *pb.IDRequest) (*pb.GetKeycloakInfosResponse, error) {
  log.Debug("Request 'GetKeycloakInfos' ")
  log.Warn("Not Implemented gRPC API: 'GetKeycloakInfos'")
  
  return &pb.GetKeycloakInfosResponse{
    Infos: nil,
  }, nil

}

func (s *KeycloakInfoServer) GetKeycloakInfo(ctx context.Context, in *pb.GetKeycloakInfoRequest) (*pb.GetKeycloakInfoResponse, error) {
  log.Debug("Request 'GetKeycloakInfo' ")
  log.Warn("Not Implemented gRPC API: 'GetKeycloakInfo'")
  return &pb.GetKeycloakInfoResponse{
    Code:  pb.Code_OK_UNSPECIFIED,
    Error: nil,
    ClusterId: "",
    Realm: "",
    ClientId: "",
    Secret: "",
    PrivateKey: "",
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

