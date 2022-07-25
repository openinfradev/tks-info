package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/openinfradev/tks-common/pkg/helper"
	"github.com/openinfradev/tks-common/pkg/log"
	"github.com/openinfradev/tks-info/pkg/csp_info"
	pb "github.com/openinfradev/tks-proto/tks_pb"
)

var (
	cspInfoAccessor *csp_info.CspInfoAccessor
)

type CspInfoServer struct {
	pb.UnimplementedCspInfoServiceServer
}

func InitCspInfoHandler(db *gorm.DB) {
	cspInfoAccessor = csp_info.New(db)
}

// CreateCSPInfo create new CSP Info for the contract id.
func (s *CspInfoServer) CreateCSPInfo(ctx context.Context, in *pb.CreateCSPInfoRequest) (*pb.IDResponse, error) {
	log.Info("Request CreateCSPInfo for contractID ", in.GetContractId())

	contractId := in.GetContractId()
	if !helper.ValidateContractId(contractId) {
		return &pb.IDResponse{
			Code: pb.Code_INVALID_ARGUMENT,
			Error: &pb.Error{
				Msg: fmt.Sprintf("invalid contract ID %s", contractId),
			},
		}, fmt.Errorf("invalid contract ID %s", contractId)
	}

	id, err := cspInfoAccessor.Create(contractId, in.GetCspName(), in.GetAuth(), in.GetCspType())
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

// GetCSPInfo is used to get CSP Info by id.
func (s *CspInfoServer) GetCSPInfo(ctx context.Context, in *pb.IDRequest) (*pb.GetCSPInfoResponse, error) {
	log.Info("request GetCSPInfo for CSP ID ", in.GetId())

	cspId, err := uuid.Parse(in.GetId())
	if err != nil {
		return &pb.GetCSPInfoResponse{
			Code: pb.Code_INVALID_ARGUMENT,
			Error: &pb.Error{
				Msg: fmt.Sprintf("invalid csp ID %s", in.GetId()),
			},
		}, err
	}

	cspInfo, err2 := cspInfoAccessor.GetCSPInfo(cspId)
	if err2 != nil {
		return &pb.GetCSPInfoResponse{
			Code: pb.Code_NOT_FOUND,
			Error: &pb.Error{
				Msg: err2.Error(),
			},
		}, err2
	}

	return &pb.GetCSPInfoResponse{
		Code:       pb.Code_OK_UNSPECIFIED,
		Error:      nil,
		ContractId: cspInfo.ContractID,
		CspName:    cspInfo.Name,
		Auth:       cspInfo.Auth,
		CspType:    cspInfo.CspType,
	}, nil
}

// GetCSPIDsByContractID returns the CSP ids by the contract id.
func (s *CspInfoServer) GetCSPIDsByContractID(ctx context.Context, in *pb.IDRequest) (*pb.IDsResponse, error) {
	log.Info("request GetCSPIDsByContractID for contract ID ", in.GetId())

	contractId := in.GetId()
	if !helper.ValidateContractId(contractId) {
		return &pb.IDsResponse{
			Code: pb.Code_INVALID_ARGUMENT,
			Error: &pb.Error{
				Msg: fmt.Sprintf("invalid contract ID %s", contractId),
			},
			Ids: nil,
		}, fmt.Errorf("invalid contract ID %s", contractId)
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
	log.Info("request UpdateCSPAuth for CSP ID ", in.GetCspId())

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
	log.Info("request GetCSPAuth for CSP ID ", in.GetId())

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
