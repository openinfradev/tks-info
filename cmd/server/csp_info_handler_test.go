package main

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	pb "github.com/openinfradev/tks-proto/tks_pb"
)

var (
	createdCspInfoId string
	requestCreateCSPInfo *pb.CreateCSPInfoRequest
)

func init() {
	requestCreateCSPInfo = randomCreateCSPInfoRequest()
}

// Tests

func TestCreateCSPInfo(t *testing.T) {
	testCases := []struct {
		name          string
		in            *pb.CreateCSPInfoRequest
		checkResponse func(req *pb.CreateCSPInfoRequest, res *pb.IDResponse, err error)
	}{
		{
			name: "OK",
			in:   requestCreateCSPInfo,
			checkResponse: func(req *pb.CreateCSPInfoRequest, res *pb.IDResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)

				_createdCspInfoId, err := uuid.Parse(res.Id)
				require.NoError(t, err)

				createdCspInfoId = _createdCspInfoId.String()
				t.Logf("createdCspInfoId : %s", createdCspInfoId)
			},
		},
		{
			name: "INVALID_CONTRACT_ID",
			in:   &pb.CreateCSPInfoRequest{
				ContractId : "THIS_IS_NOT_UUID",
			},
			checkResponse: func(req *pb.CreateCSPInfoRequest, res *pb.IDResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "OK_DUPLICATE",
			in:   requestCreateCSPInfo,
			checkResponse: func(req *pb.CreateCSPInfoRequest, res *pb.IDResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			s := CspInfoServer{}
			res, err := s.CreateCSPInfo(ctx, tc.in)
			tc.checkResponse(tc.in, res, err)
		})
	}
}

func TestGetCSPInfo(t *testing.T) {
	testCases := []struct {
		name          string
		in            *pb.IDRequest
		checkResponse func(req *pb.IDRequest, res *pb.GetCSPInfoResponse, err error)
	}{
		{
			name: "OK",
			in:   &pb.IDRequest{
				Id: createdCspInfoId,
			},
			checkResponse: func(req *pb.IDRequest, res *pb.GetCSPInfoResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)

				require.Equal(t, requestCreateCSPInfo.GetContractId(), res.GetContractId())
				require.Equal(t, requestCreateCSPInfo.GetCspName(), res.GetCspName())
				require.Equal(t, requestCreateCSPInfo.GetAuth(), res.GetAuth())
				require.Equal(t, requestCreateCSPInfo.GetCspType(), res.GetCspType())
			},
		},
		{
			name: "INVALID_CSP_ID",
			in:   &pb.IDRequest{
				Id : "THIS_IS_NOT_UUID",
			},
			checkResponse: func(req *pb.IDRequest, res *pb.GetCSPInfoResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "NOT_EXIST_CSP_ID",
			in:   &pb.IDRequest{
				Id : uuid.New().String(),
			},
			checkResponse: func(req *pb.IDRequest, res *pb.GetCSPInfoResponse, err error) {
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

			s := CspInfoServer{}
			res, err := s.GetCSPInfo(ctx, tc.in)
			tc.checkResponse(tc.in, res, err)
		})
	}
}

func TestGetCSPIDsByContractID(t *testing.T) {
	testCases := []struct {
		name          string
		in            *pb.IDRequest
		checkResponse func(req *pb.IDRequest, res *pb.IDsResponse, err error)
	}{
		{
			name: "OK",
			in:   &pb.IDRequest{
				Id: requestCreateCSPInfo.GetContractId(),
			},
			checkResponse: func(req *pb.IDRequest, res *pb.IDsResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)

				require.True(t, len(res.Ids) == 2)
			},
		},
		{
			name: "INVALID_CONTRACT_ID",
			in:   &pb.IDRequest{
				Id : "THIS_IS_NOT_UUID",
			},
			checkResponse: func(req *pb.IDRequest, res *pb.IDsResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "NOT_EXIST_CONTRACT_ID",
			in:   &pb.IDRequest{
				Id : uuid.New().String(),
			},
			checkResponse: func(req *pb.IDRequest, res *pb.IDsResponse, err error) {
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

			s := CspInfoServer{}
			res, err := s.GetCSPIDsByContractID(ctx, tc.in)
			tc.checkResponse(tc.in, res, err)
		})
	}
}

func TestUpdateCSPAuth(t *testing.T) {
	testCases := []struct {
		name          string
		in            *pb.UpdateCSPAuthRequest
		checkResponse func(req *pb.UpdateCSPAuthRequest, res *pb.SimpleResponse, err error)
	}{
		{
			name: "OK",
			in:   &pb.UpdateCSPAuthRequest{
				CspId: createdCspInfoId,
				Auth: "updated_pem",
			},
			checkResponse: func(req *pb.UpdateCSPAuthRequest, res *pb.SimpleResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)

				requestCreateCSPInfo.Auth = "updated_pem"
			},
		},
		{
			name: "INVALID_CSP_ID",
			in:   &pb.UpdateCSPAuthRequest{
				CspId : "THIS_IS_NOT_UUID",
			},
			checkResponse: func(req *pb.UpdateCSPAuthRequest, res *pb.SimpleResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "NOT_EXIST_CSP_ID",
			in:   &pb.UpdateCSPAuthRequest{
				CspId : uuid.New().String(),
			},
			checkResponse: func(req *pb.UpdateCSPAuthRequest, res *pb.SimpleResponse, err error) {
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

			s := CspInfoServer{}
			res, err := s.UpdateCSPAuth(ctx, tc.in)
			tc.checkResponse(tc.in, res, err)
		})
	}
}

func TestGetCSPAuth(t *testing.T) {
	testCases := []struct {
		name          string
		in            *pb.IDRequest
		checkResponse func(req *pb.IDRequest, res *pb.GetCSPAuthResponse, err error)
	}{
		{
			name: "OK",
			in:   &pb.IDRequest{
				Id: createdCspInfoId,
			},
			checkResponse: func(req *pb.IDRequest, res *pb.GetCSPAuthResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)

				require.Equal(t, requestCreateCSPInfo.GetAuth(), res.GetAuth() )
			},
		},
		{
			name: "INVALID_CSP_ID",
			in:   &pb.IDRequest{
				Id : "THIS_IS_NOT_UUID",
			},
			checkResponse: func(req *pb.IDRequest, res *pb.GetCSPAuthResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "NOT_EXIST_CSP_ID",
			in:   &pb.IDRequest{
				Id : uuid.New().String(),
			},
			checkResponse: func(req *pb.IDRequest, res *pb.GetCSPAuthResponse, err error) {
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

			s := CspInfoServer{}
			res, err := s.GetCSPAuth(ctx, tc.in)
			tc.checkResponse(tc.in, res, err)
		})
	}
}


// Helpers

func randomCreateCSPInfoRequest() *pb.CreateCSPInfoRequest {
	return &pb.CreateCSPInfoRequest{
		ContractId : uuid.New().String(),
		CspName : randomString("cspname"),
		Auth : "premkey",
		CspType : pb.CspType_AWS,
	}
}
