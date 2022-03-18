package main

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	pb "github.com/openinfradev/tks-proto/tks_pb"
)

var (
	createdClusterId      string
	requestAddClusterInfo *pb.AddClusterInfoRequest
)

func init() {
	requestAddClusterInfo = randomAddClusterInfoRequest()
}

// Tests

func TestAddClusterInfo(t *testing.T) {
	contractId, err := uuid.Parse(requestAddClusterInfo.ContractId)
	require.NoError(t, err)
	cspId, err := cspInfoAccessor.Create(contractId, "csp", "auth", 1)
	require.NoError(t, err)
	requestAddClusterInfo.CspId = cspId.String()

	testCases := []struct {
		name          string
		in            *pb.AddClusterInfoRequest
		checkResponse func(req *pb.AddClusterInfoRequest, res *pb.IDResponse, err error)
	}{
		{
			name: "OK",
			in:   requestAddClusterInfo,
			checkResponse: func(req *pb.AddClusterInfoRequest, res *pb.IDResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)

				_createdClusterId, err := uuid.Parse(res.Id)
				require.NoError(t, err)

				createdClusterId = _createdClusterId.String()
				t.Logf("createdClusterId : %s", createdClusterId)
			},
		},
		{
			name: "INVALID_CONTRACT_ID",
			in: &pb.AddClusterInfoRequest{
				ContractId: "NO_UUID_STRING",
			},
			checkResponse: func(req *pb.AddClusterInfoRequest, res *pb.IDResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "INVALID_CSP_ID",
			in: &pb.AddClusterInfoRequest{
				ContractId: "NO_UUID_STRING",
			},
			checkResponse: func(req *pb.AddClusterInfoRequest, res *pb.IDResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "NOT_EXISTED_CSP_ID",
			in: &pb.AddClusterInfoRequest{
				ContractId: requestAddClusterInfo.ContractId,
				CspId:      uuid.New().String(),
			},
			checkResponse: func(req *pb.AddClusterInfoRequest, res *pb.IDResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_NOT_FOUND)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			s := ClusterInfoServer{}
			res, err := s.AddClusterInfo(ctx, tc.in)
			tc.checkResponse(tc.in, res, err)
		})
	}
}

func TestGetCluster(t *testing.T) {
	testCases := []struct {
		name          string
		in            *pb.GetClusterRequest
		checkResponse func(req *pb.GetClusterRequest, res *pb.GetClusterResponse, err error)
	}{
		{
			name: "OK",
			in: &pb.GetClusterRequest{
				ClusterId: createdClusterId,
			},
			checkResponse: func(req *pb.GetClusterRequest, res *pb.GetClusterResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)

				require.Equal(t, res.Cluster.Id, createdClusterId)
				require.Equal(t, res.Cluster.ContractId, requestAddClusterInfo.ContractId)
				require.Equal(t, res.Cluster.Name, requestAddClusterInfo.Name)
				require.Equal(t, res.Cluster.Status, pb.ClusterStatus_UNSPECIFIED)
			},
		},
		{
			name: "INVALID_CLUSTER_ID",
			in: &pb.GetClusterRequest{
				ClusterId: "NO_UUID_STRING",
			},
			checkResponse: func(req *pb.GetClusterRequest, res *pb.GetClusterResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "NOT_EXISTED_CLUSTER",
			in: &pb.GetClusterRequest{
				ClusterId: uuid.New().String(),
			},
			checkResponse: func(req *pb.GetClusterRequest, res *pb.GetClusterResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_NOT_FOUND)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			s := ClusterInfoServer{}
			res, err := s.GetCluster(ctx, tc.in)
			tc.checkResponse(tc.in, res, err)
		})
	}
}

func TestGetClusters(t *testing.T) {
	testCases := []struct {
		name          string
		in            *pb.GetClustersRequest
		checkResponse func(req *pb.GetClustersRequest, res *pb.GetClustersResponse, err error)
	}{
		{
			name: "OK_BY_CONTRACT_ID",
			in: &pb.GetClustersRequest{
				ContractId: requestAddClusterInfo.ContractId,
			},
			checkResponse: func(req *pb.GetClustersRequest, res *pb.GetClustersResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)

				require.True(t, len(res.Clusters) == 1)
				require.Equal(t, res.Clusters[0].Id, createdClusterId)
			},
		},
		{
			name: "OK_BY_CSP_ID",
			in: &pb.GetClustersRequest{
				CspId: requestAddClusterInfo.CspId,
			},
			checkResponse: func(req *pb.GetClustersRequest, res *pb.GetClustersResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)

				require.True(t, len(res.Clusters) == 1)
				require.Equal(t, res.Clusters[0].Id, createdClusterId)
			},
		},
		{
			name: "INVALID_CONTRACT_ID",
			in: &pb.GetClustersRequest{
				ContractId: "NO_UUID_STRING",
			},
			checkResponse: func(req *pb.GetClustersRequest, res *pb.GetClustersResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "INVALID_CSP_ID",
			in: &pb.GetClustersRequest{
				ContractId: "NO_UUID_STRING",
			},
			checkResponse: func(req *pb.GetClustersRequest, res *pb.GetClustersResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "CONTRACT_ID_AND_CSP_ID_PROVIDED",
			in: &pb.GetClustersRequest{
				ContractId: uuid.New().String(),
				CspId:      uuid.New().String(),
			},
			checkResponse: func(req *pb.GetClustersRequest, res *pb.GetClustersResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "NOT_EXISTED_CLUSTERS_BY_CONTRACT_ID",
			in: &pb.GetClustersRequest{
				ContractId: uuid.New().String(),
			},
			checkResponse: func(req *pb.GetClustersRequest, res *pb.GetClustersResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_NOT_FOUND)
			},
		},
		{
			name: "NOT_EXISTED_CLUSTERS_BY_CSP_ID",
			in: &pb.GetClustersRequest{
				CspId: uuid.New().String(),
			},
			checkResponse: func(req *pb.GetClustersRequest, res *pb.GetClustersResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_NOT_FOUND)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			s := ClusterInfoServer{}
			res, err := s.GetClusters(ctx, tc.in)
			tc.checkResponse(tc.in, res, err)
		})
	}
}

func TestUpdateStatus(t *testing.T) {
	testCases := []struct {
		name          string
		in            *pb.UpdateClusterStatusRequest
		checkResponse func(req *pb.UpdateClusterStatusRequest, res *pb.SimpleResponse, err error)
	}{
		{
			name: "OK",
			in: &pb.UpdateClusterStatusRequest{
				ClusterId: createdClusterId,
				Status:    pb.ClusterStatus_INSTALLING,
			},
			checkResponse: func(req *pb.UpdateClusterStatusRequest, res *pb.SimpleResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)

				clusterId, err := uuid.Parse(createdClusterId)
				require.NoError(t, err)

				cluster, err := clusterAccessor.GetCluster(clusterId)
				require.NoError(t, err)

				require.Equal(t, cluster.Id, createdClusterId)
				require.Equal(t, cluster.Status, pb.ClusterStatus_INSTALLING)
			},
		},
		{
			name: "INVALID_CLUSTER_ID",
			in: &pb.UpdateClusterStatusRequest{
				ClusterId: "NO_UUID_STRING",
				Status:    pb.ClusterStatus_INSTALLING,
			},
			checkResponse: func(req *pb.UpdateClusterStatusRequest, res *pb.SimpleResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "NOT_EXISTED_CLUSTER",
			in: &pb.UpdateClusterStatusRequest{
				ClusterId: uuid.New().String(),
				Status:    pb.ClusterStatus_INSTALLING,
			},
			checkResponse: func(req *pb.UpdateClusterStatusRequest, res *pb.SimpleResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INTERNAL)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			s := ClusterInfoServer{}
			res, err := s.UpdateClusterStatus(ctx, tc.in)
			tc.checkResponse(tc.in, res, err)
		})
	}
}

// Helpers

func randomAddClusterInfoRequest() *pb.AddClusterInfoRequest {
	return &pb.AddClusterInfoRequest{
		ContractId: uuid.New().String(),
		CspId:      uuid.New().String(),
		Name:       randomString("Name"),
		Conf: &pb.ClusterConf{
			SshKeyName:   randomString("SSHKEYNAME"),
			Region:       randomString("REGION"),
			NumOfAz:      3,
			MachineType:  randomString("MACHINETYPE"),
			MinSizePerAz: 1,
			MaxSizePerAz: 5,
		},
	}
}
