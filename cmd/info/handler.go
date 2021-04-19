package info

import (
	"context"

	"github.com/sktelecom/tks-contract/pkg/log"
	"github.com/sktelecom/tks-info/pkg/info/cluster"
	"github.com/sktelecom/tks-info/pkg/info/csp"
	pb "github.com/sktelecom/tks-proto/pbgo"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	cspAccessor     *csp.Accessor
	clusterAccessor *cluster.Accessor
)

type Server struct {
	pb.UnimplementedInfoServiceServer
}

func init() {
	cspAccessor = csp.NewCSPAccessor()
	clusterAccessor = cluster.NewClusterAccessor()
}

// CreateCSPInfo create new CSP Info for the contract id.
func (s *Server) CreateCSPInfo(ctx context.Context, in *pb.CreateCSPInfoRequest) (*pb.IDResponse, error) {
	log.Info("request CreateCSPInfo for contractID ", in.GetContractId())
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
	log.Info("request AddClusterInfo for Contract ID", in.GetContractId())
	// return an error if csp id does not exist.
	if _, err := cspAccessor.Get(csp.ID(in.GetCspId())); err != nil {
		return &pb.IDResponse{
			Code: pb.Code_NOT_FOUND,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}
	cID, err := clusterAccessor.Create(cluster.ID(in.GetContractId()), cluster.ID(in.GetCspId()), in.GetName(), in.GetConf())
	if err != nil {
		return &pb.IDResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}
	log.Info("Created new cluster id:", cID)
	return &pb.IDResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
		Id:    string(cID),
	}, nil
}

// GetCluster get cluster for the id of the cluster
func (s *Server) GetCluster(ctx context.Context, in *pb.GetClusterRequest) (*pb.GetClusterResponse, error) {
	cluster, err := clusterAccessor.Get(cluster.ID(in.GetClusterId()))
	if err != nil {
		return &pb.GetClusterResponse{
			Code: pb.Code_NOT_FOUND,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}
	return &pb.GetClusterResponse{
		Code:    pb.Code_OK_UNSPECIFIED,
		Error:   nil,
		Cluster: clusterAccessor.ClustertToPbCluster(cluster),
	}, nil
}

// GetClusters get every clusters by csp id
func (s *Server) GetClusters(ctx context.Context, in *pb.GetClustersRequest) (*pb.GetClustersResponse, error) {
	clusters, err := clusterAccessor.GetClustersByCspID(cluster.ID(in.GetCspId()))
	if err != nil {
		return &pb.GetClustersResponse{
			Code: pb.Code_NOT_FOUND,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}
	res := []*pb.Cluster{}
	for _, cluster := range clusters {
		res = append(res, clusterAccessor.ClustertToPbCluster(cluster))
	}
	return &pb.GetClustersResponse{
		Code:     pb.Code_OK_UNSPECIFIED,
		Error:    nil,
		Clusters: res,
	}, nil
}

// UpdateClusterStatus update Status of the Cluster
func (s *Server) UpdateClusterStatus(ctx context.Context, in *pb.UpdateClusterStatusRequest) (*pb.SimpleResponse, error) {
	err := clusterAccessor.UpdateStatus(cluster.ID(in.GetClusterId()), in.GetStatus())
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

// ValidateLabelUniqueness check uniqueness of the label
func (s *Server) ValidateLabelUniqueness(ctx context.Context, in *pb.ValidateLabelUniquenessRequest) (*pb.ValidateLabelUniquenessResponse, error) {
	return nil, nil
}
