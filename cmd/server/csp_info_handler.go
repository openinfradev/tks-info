package main

import (
  "context"
  "fmt"
  "gorm.io/gorm"

  "github.com/google/uuid"
  "github.com/sktelecom/tks-contract/pkg/log"
  "github.com/sktelecom/tks-info/pkg/csp_info"
  pb "github.com/sktelecom/tks-proto/pbgo"
)

//var (
//  cspInfoAccessor     *csp_info.Accessor
//)

type CspInfoServer struct {
  pb.UnimplementedCspInfoServiceServer
}

func InitCspInfoHandler(db *gorm.DB) {
  cspInfoAccessor = csp_info.New(db)
}

// CreateCSPInfo create new CSP Info for the contract id.
func (s *CspInfoServer) CreateCSPInfo(ctx context.Context, in *pb.CreateCSPInfoRequest) (*pb.IDResponse, error) {
  log.Info("Request CreateCSPInfo for contractID ", in.GetContractId())

  contractId, err := uuid.Parse(in.GetContractId())
  if err != nil {
    res := pb.IDResponse{
      Code: pb.Code_INVALID_ARGUMENT,
      Error: &pb.Error{
        Msg: fmt.Sprintf("invalid contract ID %s", in.GetContractId()),
      },
    }
    return &res, err
  }

  id, err := cspInfoAccessor.Create(contractId, in.GetCspName(), in.GetAuth())
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

// GetCSPIDsByContractID returns the CSP ids by the contract id.
func (s *CspInfoServer) GetCSPIDsByContractID(ctx context.Context, in *pb.IDRequest) (*pb.IDsResponse, error) {
  log.Debug("request GetCSPIDsByContractID for contract ID ", in.GetId())

  contractId, err := uuid.Parse(in.GetId())
  if err != nil {
    res := pb.IDsResponse{
      Code: pb.Code_INVALID_ARGUMENT,
      Error: &pb.Error{
        Msg: fmt.Sprintf("invalid contract ID %s", in.GetId()),
      },
      Ids: nil,
    }
    return &res, err
  }

  ids, err := cspInfoAccessor.GetCSPIDsByContractID(contractId)
  if err != nil {
    return &pb.IDsResponse{
      Code: pb.Code_NOT_FOUND,
      Error: &pb.Error{
        Msg: err.Error(),
      },
      Ids: nil,
    }, err
  }

  return &pb.IDsResponse{
    Code:  pb.Code_OK_UNSPECIFIED,
    Error: nil,
    Ids:   ids,
  }, nil
}

// UpdateCSPAuth updates an authentication config for CSP.
func (s *CspInfoServer) UpdateCSPAuth(ctx context.Context, in *pb.UpdateCSPAuthRequest) (*pb.SimpleResponse, error) {
  log.Debug("request UpdateCSPAuth for CSP ID ", in.GetCspId())

  cspId, err := uuid.Parse(in.GetCspId())
  if err != nil {
    res := pb.SimpleResponse{
      Code: pb.Code_INVALID_ARGUMENT,
      Error: &pb.Error{
        Msg: fmt.Sprintf("invalid csp ID %s", in.GetCspId()),
      },
    }
    return &res, err
  }

  if err := cspInfoAccessor.UpdateCSPAuth(cspId, in.GetAuth()); err != nil {
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

// GetCSPAuth returns an authentication info by csp id.
func (s *CspInfoServer) GetCSPAuth(ctx context.Context, in *pb.IDRequest) (*pb.GetCSPAuthResponse, error) {
  log.Debug("request GetCSPAuth for CSP ID ", in.GetId())

  cspId, err := uuid.Parse(in.GetId())
  if err != nil {
    res := pb.GetCSPAuthResponse{
      Code: pb.Code_INVALID_ARGUMENT,
      Error: &pb.Error{
        Msg: fmt.Sprintf("invalid csp ID %s", in.GetId()),
      },
    }
    return &res, err
  }

  cspInfo, err2 := cspInfoAccessor.GetCSPInfo(cspId)
  if err2 != nil {
    res := pb.GetCSPAuthResponse{
      Code: pb.Code_NOT_FOUND,
      Error: &pb.Error{
        Msg: err2.Error(),
      },
    }
    return &res, err2
  }

  return &pb.GetCSPAuthResponse{
    Code:  pb.Code_OK_UNSPECIFIED,
    Error: nil,
    Auth:  cspInfo.Auth,
  }, nil
}
