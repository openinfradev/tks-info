package main

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	pb "github.com/openinfradev/tks-proto/tks_pb"
)

var (
	createdKeycloakInfoId     string
	requestCreateKeycloakInfo *pb.CreateKeycloakInfoRequest
)

func init() {
	requestCreateKeycloakInfo = randomCreateKeycloakRequest()
}

// Tests

func TestCreateKeycloakInfo(t *testing.T) {
	testCases := []struct {
		name          string
		in            *pb.CreateKeycloakInfoRequest
		checkResponse func(req *pb.CreateKeycloakInfoRequest, res *pb.IDResponse, err error)
	}{
		{
			name: "OK",
			in:   requestCreateKeycloakInfo,
			checkResponse: func(req *pb.CreateKeycloakInfoRequest, res *pb.IDResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)

				keycloakInfoId, err := uuid.Parse(res.Id)
				require.NoError(t, err)

				createdKeycloakInfoId = keycloakInfoId.String()
				t.Logf("createdKeycloakInfoId : %s", createdKeycloakInfoId)
			},
		},
		{
			name: "INVALID_ARGUMENT_CLUSTERID",
			in: &pb.CreateKeycloakInfoRequest{
				ClusterId: randomString("NOT_UUID"),
			},
			checkResponse: func(req *pb.CreateKeycloakInfoRequest, res *pb.IDResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "INVALID_ARGUMENT_EMPTY_REQUEST",
			in:   &pb.CreateKeycloakInfoRequest{},
			checkResponse: func(req *pb.CreateKeycloakInfoRequest, res *pb.IDResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			s := KeycloakInfoServer{}
			res, err := s.CreateKeycloakInfo(ctx, tc.in)
			tc.checkResponse(tc.in, res, err)
		})
	}
}

func TestGetKeycloakInfoByClusterId(t *testing.T) {

	testCases := []struct {
		name          string
		in            *pb.IDRequest
		checkResponse func(req *pb.IDRequest, res *pb.GetKeycloakInfoResponse, err error)
	}{
		{
			name: "OK",
			in: &pb.IDRequest{
				Id: requestCreateKeycloakInfo.GetClusterId(),
			},
			checkResponse: func(req *pb.IDRequest, res *pb.GetKeycloakInfoResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)
				require.True(t, len(res.GetKeycloakInfos()) > 0)
			},
		},
		{
			name: "INVALID_ARGUMENT_CLUSTERID",
			in: &pb.IDRequest{
				Id: randomString("NOT_UUID"),
			},
			checkResponse: func(req *pb.IDRequest, res *pb.GetKeycloakInfoResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "INTERNAL_NO_KEYCLOAK_BY_CLUSTER_ID",
			in: &pb.IDRequest{
				Id: uuid.New().String(),
			},
			checkResponse: func(req *pb.IDRequest, res *pb.GetKeycloakInfoResponse, err error) {
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

			s := KeycloakInfoServer{}
			res, err := s.GetKeycloakInfoByClusterId(ctx, tc.in)
			tc.checkResponse(tc.in, res, err)
		})
	}

}

// Helpers

func randomCreateKeycloakRequest() *pb.CreateKeycloakInfoRequest {
	return &pb.CreateKeycloakInfoRequest{
		ClusterId:  uuid.New().String(),
		Realm:      randomString("Realm"),
		ClientId:   randomString("ClientId"),
		Secret:     randomString("Secret"),
		PrivateKey: randomString("PrivateKey"),
	}
}
