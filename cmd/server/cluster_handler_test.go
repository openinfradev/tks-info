package main

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/openinfradev/tks-common/pkg/helper"
	pb "github.com/openinfradev/tks-proto/tks_pb"
	mocktks "github.com/openinfradev/tks-proto/tks_pb/mock"
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
	contractId := requestAddClusterInfo.ContractId
	require.True(t, helper.ValidateContractId(contractId))

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

				t.Logf("res.Id : %s", res.Id)
				_createdClusterId := res.Id
				require.True(t, helper.ValidateClusterId(_createdClusterId))

				createdClusterId = _createdClusterId
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
				ClusterId: helper.GenerateClusterId(),
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
		buildStubs    func(mockContractClient *mocktks.MockContractServiceClient)
		checkResponse func(req *pb.GetClustersRequest, res *pb.GetClustersResponse, err error)
	}{
		{
			name: "OK_BY_CONTRACT_ID",
			in: &pb.GetClustersRequest{
				ContractId: requestAddClusterInfo.ContractId,
			},
			buildStubs: func(mockContractClient *mocktks.MockContractServiceClient) {},
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
			buildStubs: func(mockContractClient *mocktks.MockContractServiceClient) {},
			checkResponse: func(req *pb.GetClustersRequest, res *pb.GetClustersResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)

				require.True(t, len(res.Clusters) == 1)
				require.Equal(t, res.Clusters[0].Id, createdClusterId)
			},
		},
		{
			name: "USE_DEFAULT_CONTRACT_AND_NOT_EXISTS_CLUSTERS",
			in: &pb.GetClustersRequest{
				ContractId: "",
				CspId:      "",
			},
			buildStubs: func(mockContractClient *mocktks.MockContractServiceClient) {
				mockContractClient.EXPECT().GetDefaultContract(gomock.Any(), gomock.Any()).Times(1).
					Return(
						&pb.GetContractResponse{
							Code:  pb.Code_OK_UNSPECIFIED,
							Error: nil,
							Contract: &pb.Contract{
								ContractId: helper.GenerateContractId(),
								CspId:      uuid.New().String(),
							},
						}, nil)
			},
			checkResponse: func(req *pb.GetClustersRequest, res *pb.GetClustersResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)
			},
		},
		{
			name: "FAILED_TO_GET_DEFAULT_CONTRACT",
			in: &pb.GetClustersRequest{
				ContractId: "",
				CspId:      "",
			},
			buildStubs: func(mockContractClient *mocktks.MockContractServiceClient) {
				mockContractClient.EXPECT().GetDefaultContract(gomock.Any(), gomock.Any()).Times(1).
					Return(
						&pb.GetContractResponse{
							Code:  pb.Code_OK_UNSPECIFIED,
							Error: nil,
						}, errors.New("FAILED_TO_GET_DEFAULT_CONTRACT"))
			},
			checkResponse: func(req *pb.GetClustersRequest, res *pb.GetClustersResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_NOT_FOUND)
			},
		},
		{
			name: "INVALID_CONTRACT_ID",
			in: &pb.GetClustersRequest{
				ContractId: "NO_ID_STRING",
			},
			buildStubs: func(mockContractClient *mocktks.MockContractServiceClient) {},
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
			buildStubs: func(mockContractClient *mocktks.MockContractServiceClient) {},
			checkResponse: func(req *pb.GetClustersRequest, res *pb.GetClustersResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "CONTRACT_ID_AND_CSP_ID_PROVIDED",
			in: &pb.GetClustersRequest{
				ContractId: helper.GenerateContractId(),
				CspId:      uuid.New().String(),
			},
			buildStubs: func(mockContractClient *mocktks.MockContractServiceClient) {},
			checkResponse: func(req *pb.GetClustersRequest, res *pb.GetClustersResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "NOT_EXISTED_CLUSTERS_BY_CONTRACT_ID",
			in: &pb.GetClustersRequest{
				ContractId: helper.GenerateContractId(),
			},
			buildStubs: func(mockContractClient *mocktks.MockContractServiceClient) {},
			checkResponse: func(req *pb.GetClustersRequest, res *pb.GetClustersResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)
			},
		},
		{
			name: "NOT_EXISTED_CLUSTERS_BY_CSP_ID",
			in: &pb.GetClustersRequest{
				CspId: uuid.New().String(),
			},
			buildStubs: func(mockContractClient *mocktks.MockContractServiceClient) {},
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

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockContarctClient := mocktks.NewMockContractServiceClient(ctrl)
			contractClient = mockContarctClient

			tc.buildStubs(mockContarctClient)

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

				clusterId := createdClusterId
				require.True(t, helper.ValidateClusterId(clusterId))

				cluster, err := clusterAccessor.GetCluster(clusterId)
				require.NoError(t, err)

				require.Equal(t, cluster.Id, createdClusterId)
				require.Equal(t, cluster.Status, pb.ClusterStatus_INSTALLING)
			},
		},
		{
			name: "INVALID_CLUSTER_ID",
			in: &pb.UpdateClusterStatusRequest{
				ClusterId: "NO_CID_STRING",
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
				ClusterId: helper.GenerateClusterId(),
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
		ContractId: helper.GenerateContractId(),
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
