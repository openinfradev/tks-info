package main

import (
  "context"
  "fmt"
  "gorm.io/gorm"
  "errors"

  "github.com/google/uuid"
  "github.com/sktelecom/tks-contract/pkg/log"
  "github.com/sktelecom/tks-info/pkg/cluster"
  "github.com/sktelecom/tks-info/pkg/csp_info"
  pb "github.com/sktelecom/tks-proto/pbgo"
)

var (
  clusterAccessor *cluster.ClusterAccessor
  cspInfoAccessor *csp_info.CspInfoAccessor
)

type ClusterInfoServer struct {
  pb.UnimplementedClusterInfoServiceServer
}

func InitClusterInfoHandler(db *gorm.DB) {
  clusterAccessor = cluster.New(db)
}


// AddClusterInfo add newly created cluster with csp id
func (s *ClusterInfoServer) AddClusterInfo(ctx context.Context, in *pb.AddClusterInfoRequest) (*pb.IDResponse, error) {
  log.Info("request AddClusterInfo for Contract ID", in.GetContractId())

  cspId, err := uuid.Parse(in.GetCspId())
  if err != nil {
    res := pb.IDResponse{
      Code: pb.Code_INVALID_ARGUMENT,
      Error: &pb.Error{
        Msg: fmt.Sprintf("Invalid CSP ID %s", in.GetCspId()),
      },
    }
    return &res, err
  }

  contractId, err := uuid.Parse(in.GetContractId())
  if err != nil {
    res := pb.IDResponse{
      Code: pb.Code_INVALID_ARGUMENT,
      Error: &pb.Error{
        Msg: fmt.Sprintf("Invalid contract ID %s", in.GetContractId()),
      },
    }
    return &res, err
  }

  // return an error if csp id does not exist.
  // Robert: Is this necessary? If so, what about contractID check?
  if _, err := cspInfoAccessor.GetCSPInfo(cspId); err != nil {
    return &pb.IDResponse{
      Code: pb.Code_NOT_FOUND,
      Error: &pb.Error{
        Msg: err.Error(),
      },
    }, err
  }

  // Create cluster record
  cID, err := clusterAccessor.CreateClusterInfo(contractId, cspId, in.GetName(), *(in.GetConf()))
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
    Id:    cID.String(),
  }, nil
}

// GetCluster get cluster for the id of the cluster
func (s *ClusterInfoServer) GetCluster(ctx context.Context, in *pb.GetClusterRequest) (*pb.GetClusterResponse, error) {
  clusterId, err := uuid.Parse(in.GetClusterId())
  if err != nil {
    res := pb.GetClusterResponse{
      Code: pb.Code_INVALID_ARGUMENT,
      Error: &pb.Error{
        Msg: fmt.Sprintf("Invalid cluster ID %s", in.GetClusterId()),
      },
    }
    return &res, err
  }

  cluster, err := clusterAccessor.GetCluster(clusterId)
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
    Cluster: cluster,
  }, nil
}

// GetClusters get every clusters by csp id
func (s *ClusterInfoServer) GetClusters(ctx context.Context, in *pb.GetClustersRequest) (*pb.GetClustersResponse, error) {
  contractId := in.GetContractId()
  cspId := in.GetCspId()

  if contractId == "" && cspId == "" {
    err := errors.New("Wrong parameter")
    res := pb.GetClustersResponse{
      Code: pb.Code_INVALID_ARGUMENT,
      Error: &pb.Error{
        //Msg: fmt.Sprintf("Either contractID or cspId must be provided"),
        Msg: fmt.Sprintf("Neither contractID or cspId was provided. Exactly one of those must be provided."),
      },
      Clusters: nil,
    }
    return &res, err
  } else if contractId != "" && cspId != "" {
    err := errors.New("Wrong parameter")
    res := pb.GetClustersResponse{
      Code: pb.Code_INVALID_ARGUMENT,
      Error: &pb.Error{
        Msg: fmt.Sprintf("Both contractID and cspId was provided. Exactly one of those must be provided."),
      },
      Clusters: nil,
    }
    return &res, err
  } else if contractId != "" && cspId == "" {
    /*****************************
    * Get clusters by contractID *
    *****************************/
    conIdParsed, err := uuid.Parse(contractId)
    if err != nil {
      return &pb.GetClustersResponse{
        Code: pb.Code_INVALID_ARGUMENT,
        Error: &pb.Error{
          Msg: fmt.Sprintf("Invalid Contract ID %s", contractId),
        },
        Clusters: nil,
      }, err
    }

    clusters, err := clusterAccessor.GetClustersByContractID(conIdParsed)
    if err != nil {
      return &pb.GetClustersResponse{
        Code: pb.Code_NOT_FOUND,
        Error: &pb.Error{
          Msg: err.Error(),
        },
        Clusters: nil,
      }, err
    }

    // Successfully return GetClustersResponse
    return &pb.GetClustersResponse{
      Code:     pb.Code_OK_UNSPECIFIED,
      Error:    nil,
      Clusters: clusters,
    }, nil
  } else {
    /************************
    * Get clusters by cspID *
    ************************/
    cspIdParsed, err := uuid.Parse(cspId)
    if err != nil {
      return &pb.GetClustersResponse{
        Code: pb.Code_INVALID_ARGUMENT,
        Error: &pb.Error{
          Msg: fmt.Sprintf("Invalid CSP ID %s", cspId),
        },
        Clusters: nil,
      }, err
    }

    clusters, err := clusterAccessor.GetClustersByCspID(cspIdParsed)
    if err != nil {
      return &pb.GetClustersResponse{
        Code: pb.Code_NOT_FOUND,
        Error: &pb.Error{
          Msg: err.Error(),
        },
        Clusters: nil,
      }, err
    }

    // Successfully return GetClustersResponse
    return &pb.GetClustersResponse{
      Code:     pb.Code_OK_UNSPECIFIED,
      Error:    nil,
      Clusters: clusters,
    }, nil
  }
}

// UpdateClusterStatus update Status of the Cluster
func (s *ClusterInfoServer) UpdateClusterStatus(ctx context.Context, in *pb.UpdateClusterStatusRequest) (*pb.SimpleResponse, error) {
  clusterId, err := uuid.Parse(in.GetClusterId())
  if err != nil {
    res := pb.SimpleResponse{
      Code: pb.Code_INVALID_ARGUMENT,
      Error: &pb.Error{
        Msg: fmt.Sprintf("Invalid Cluster ID %s", in.GetClusterId()),
      },
    }
    return &res, err
  }

  err = clusterAccessor.UpdateStatus(clusterId, in.GetStatus())
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
