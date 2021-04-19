package info

import (
	"context"

	"github.com/sktelecom/tks-contract/pkg/log"
	"github.com/sktelecom/tks-info/pkg/info/csp"
	pb "github.com/sktelecom/tks-proto/pbgo"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	cspAccessor *csp.Accessor
)

type Server struct {
	pb.UnimplementedInfoServiceServer
}

func init() {
	cspAccessor = csp.NewCSPAccessor()
}

// CreateCSPInfo create new CSP Info for the contract id.
func (s *Server) CreateCSPInfo(ctx context.Context, in *pb.CreateCSPInfoRequest) (*pb.IDResponse, error) {
	log.Debug("request CreateCSPInfo for contractID ", in.GetContractId())
	id, err := cspAccessor.Create(csp.ID(in.GetContractId()), in.GetAuth())
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
		Id:    string(id),
	}, nil
}

// GetCSPIDs returns all CSP ids.
func (s *Server) GetCSPIDs(ctx context.Context, empty *emptypb.Empty) (*pb.IDsResponse, error) {
	log.Debug("request GetCSPIDs")
	cspinfos := cspAccessor.List()
	ids := []string{}

	for _, csp := range cspinfos {
		ids = append(ids, string(csp.ID))
	}

	return &pb.IDsResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
		Ids:   ids,
	}, nil
}

// GetCSPIDsByContractID returns the CSP ids by the contract id.
func (s *Server) GetCSPIDsByContractID(ctx context.Context, in *pb.IDRequest) (*pb.IDsResponse, error) {
	log.Debug("request GetCSPIDsByContractID for contract ID ", in.GetId())
	ids, err := cspAccessor.GetCSPIDsByContractID(csp.ID(in.GetId()))
	if err != nil {
		return &pb.IDsResponse{
			Code: pb.Code_NOT_FOUND,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}

	res := []string{}
	for _, id := range ids {
		res = append(res, string(id))
	}
	return &pb.IDsResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
		Ids:   res,
	}, nil
}

// UpdateCSPInfo updates an authentication config for CSP.
func (s *Server) UpdateCSPInfo(ctx context.Context, in *pb.UpdateCSPInfoRequest) (*pb.SimpleResponse, error) {
	log.Debug("request UpdateCSPInfo for CSP ID ", in.GetCspId())
	if err := cspAccessor.Update(csp.ID(in.GetCspId()), in.GetAuth()); err != nil {
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
func (s *Server) GetCSPAuth(ctx context.Context, in *pb.IDRequest) (*pb.GetCSPAuthResponse, error) {
	log.Debug("request GetCSPAuth for CSP ID ", in.GetId())
	csp, err := cspAccessor.Get(csp.ID(in.GetId()))
	if err != nil {
		return &pb.GetCSPAuthResponse{
			Code: pb.Code_NOT_FOUND,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}
	return &pb.GetCSPAuthResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
		Auth:  csp.Auth,
	}, nil
}

// AddClusterInfo add newly created cluster with csp id
func (s *Server) AddClusterInfo(ctx context.Context, in *pb.AddClusterInfoRequest) (*pb.IDResponse, error) {
	return nil, nil
}

// GetClusterget cluster for the id of the cluster
func (s *Server) GetCluster(ctx context.Context, in *pb.GetClusterRequest) (*pb.GetClusterResponse, error) {
	return nil, nil
}

// GetClusters get every clusters on the mutlcluster
func (s *Server) GetClusters(ctx context.Context, in *pb.GetClustersRequest) (*pb.GetClustersResponse, error) {
	return nil, nil
}

// UpdateClusterStatus update Status of the Cluster
func (s *Server) UpdateClusterStatus(ctx context.Context, in *pb.UpdateClusterStatusRequest) (*pb.SimpleResponse, error) {
	return nil, nil
}

// ValidateLabelUniqueness check uniqueness of the label
func (s *Server) ValidateLabelUniqueness(ctx context.Context, in *pb.ValidateLabelUniquenessRequest) (*pb.ValidateLabelUniquenessResponse, error) {
	return nil, nil
}
