package main

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	pb "github.com/openinfradev/tks-proto/tks_pb"
)

var (
	createdAppGroupId  string
	requestCreateAppGroup *pb.CreateAppGroupRequest
)

func init() {
	requestCreateAppGroup = randomCreateAppGroupRequest()
}


// Tests

func TestCreateAppGroup(t *testing.T) {
	testCases := []struct {
		name          string
		in            *pb.CreateAppGroupRequest
		checkResponse func(req *pb.CreateAppGroupRequest, res *pb.IDResponse, err error)
	}{
		{
			name: "OK",
			in:   requestCreateAppGroup,
			checkResponse: func(req *pb.CreateAppGroupRequest, res *pb.IDResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)

				_createdAppGroupId, err := uuid.Parse(res.Id)
				require.NoError(t, err)

				createdAppGroupId = _createdAppGroupId.String()
				requestCreateAppGroup.AppGroup.AppGroupId = createdAppGroupId
				t.Logf("createdAppGroupId : %s", createdAppGroupId)
			},
		},
		{
			name: "INVALID_CLUSTER_ID",
			in:   &pb.CreateAppGroupRequest{
				ClusterId : "NO_UUID_STRING",
			},
			checkResponse: func(req *pb.CreateAppGroupRequest, res *pb.IDResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "DUPLICATE_EXTERNAL_LABEL",
			in:   requestCreateAppGroup,
			checkResponse: func(req *pb.CreateAppGroupRequest, res *pb.IDResponse, err error) {
				require.Equal(t, res.Code, pb.Code_INTERNAL)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			s := AppInfoServer{}
			res, err := s.CreateAppGroup(ctx, tc.in)
			tc.checkResponse(tc.in, res, err)
		})
	}
}

func TestGetAppGroupsByClusterID(t *testing.T) {
	testCases := []struct {
		name          string
		in            *pb.IDRequest
		checkResponse func(req *pb.IDRequest, res *pb.GetAppGroupsResponse, err error)
	}{
		{
			name: "OK",
			in:   &pb.IDRequest{
				Id: requestCreateAppGroup.GetClusterId(),
			},
			checkResponse: func(req *pb.IDRequest, res *pb.GetAppGroupsResponse, err error) {
				t.Logf( "res %s", res )
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)

				require.True(t, len(res.GetAppGroups()) == 1 )
				require.Equal(t, res.GetAppGroups()[0].GetClusterId(), requestCreateAppGroup.GetClusterId())
				require.Equal(t, res.GetAppGroups()[0].GetAppGroupId(), requestCreateAppGroup.GetAppGroup().GetAppGroupId())
				require.Equal(t, res.GetAppGroups()[0].GetAppGroupName(), requestCreateAppGroup.GetAppGroup().GetAppGroupName())
			},
		},
		{
			name: "INVALID_CLUSTER_ID",
			in:   &pb.IDRequest{
				Id : "NO_UUID_STRING",
			},
			checkResponse: func(req *pb.IDRequest, res *pb.GetAppGroupsResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "NOTE_EXIST_APPGROUPS",
			in:   &pb.IDRequest{
				Id : uuid.New().String(),
			},
			checkResponse: func(req *pb.IDRequest, res *pb.GetAppGroupsResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)
				require.True(t, len(res.GetAppGroups()) == 0)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			s := AppInfoServer{}
			res, err := s.GetAppGroupsByClusterID(ctx, tc.in)
			tc.checkResponse(tc.in, res, err)
		})
	}
}

func TestGetAppGroups(t *testing.T) {
	testCases := []struct {
		name          string
		in            *pb.GetAppGroupsRequest
		checkResponse func(req *pb.GetAppGroupsRequest, res *pb.GetAppGroupsResponse, err error)
	}{
		{
			name: "OK",
			in:   &pb.GetAppGroupsRequest{
				AppGroupName: requestCreateAppGroup.GetAppGroup().GetAppGroupName(),
				Type: requestCreateAppGroup.GetAppGroup().GetType(),
			},
			checkResponse: func(req *pb.GetAppGroupsRequest, res *pb.GetAppGroupsResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)

				require.True(t, len(res.GetAppGroups()) == 1 )
				require.Equal(t, res.GetAppGroups()[0].GetClusterId(), requestCreateAppGroup.GetClusterId())
				require.Equal(t, res.GetAppGroups()[0].GetAppGroupId(), requestCreateAppGroup.GetAppGroup().GetAppGroupId())
				require.Equal(t, res.GetAppGroups()[0].GetAppGroupName(), requestCreateAppGroup.GetAppGroup().GetAppGroupName())
			},
		},
		{
			name: "OK_NO_TYPE",
			in:   &pb.GetAppGroupsRequest{
				AppGroupName: requestCreateAppGroup.GetAppGroup().GetAppGroupName(),
			},
			checkResponse: func(req *pb.GetAppGroupsRequest, res *pb.GetAppGroupsResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)
				require.True(t, len(res.GetAppGroups()) == 1 )
			},
		},
		{
			name: "INVALID_TYPE",
			in:   &pb.GetAppGroupsRequest{
				AppGroupName: requestCreateAppGroup.GetAppGroup().GetAppGroupName(),
				Type: pb.AppGroupType_SERVICE_MESH,
			},
			checkResponse: func(req *pb.GetAppGroupsRequest, res *pb.GetAppGroupsResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INTERNAL)
			},
		},
		{
			name: "EMPTY_NAME",
			in:   &pb.GetAppGroupsRequest{
				AppGroupName: "",
			},
			checkResponse: func(req *pb.GetAppGroupsRequest, res *pb.GetAppGroupsResponse, err error) {
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

			s := AppInfoServer{}
			res, err := s.GetAppGroups(ctx, tc.in)
			tc.checkResponse(tc.in, res, err)
		})
	}
}

func TestGetAppGroup(t *testing.T) {
	testCases := []struct {
		name          string
		in            *pb.GetAppGroupRequest
		checkResponse func(req *pb.GetAppGroupRequest, res *pb.GetAppGroupResponse, err error)
	}{
		{
			name: "OK",
			in:   &pb.GetAppGroupRequest{
				AppGroupId: createdAppGroupId,
			},
			checkResponse: func(req *pb.GetAppGroupRequest, res *pb.GetAppGroupResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)

				require.Equal(t, res.GetAppGroup().GetClusterId(), requestCreateAppGroup.GetClusterId())
				require.Equal(t, res.GetAppGroup().GetAppGroupId(), requestCreateAppGroup.GetAppGroup().GetAppGroupId())
				require.Equal(t, res.GetAppGroup().GetAppGroupName(), requestCreateAppGroup.GetAppGroup().GetAppGroupName())
				
			},
		},
		{
			name: "INVALID_APPGROUP_ID",
			in:   &pb.GetAppGroupRequest{
				AppGroupId: "NO_UUID_STRING",
			},
			checkResponse: func(req *pb.GetAppGroupRequest, res *pb.GetAppGroupResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "NO_APPGROUP_ID",
			in:   &pb.GetAppGroupRequest{
				AppGroupId: uuid.New().String(),
			},
			checkResponse: func(req *pb.GetAppGroupRequest, res *pb.GetAppGroupResponse, err error) {
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

			s := AppInfoServer{}
			res, err := s.GetAppGroup(ctx, tc.in)
			tc.checkResponse(tc.in, res, err)
		})
	}
}

func TestUpdateAppGroupStatus(t *testing.T) {
	testCases := []struct {
		name          string
		in            *pb.UpdateAppGroupStatusRequest
		checkResponse func(req *pb.UpdateAppGroupStatusRequest, res *pb.SimpleResponse, err error)
	}{
		{
			name: "OK",
			in:   &pb.UpdateAppGroupStatusRequest{
				AppGroupId: createdAppGroupId,
				Status: pb.AppGroupStatus_APP_GROUP_RUNNING,
			},
			checkResponse: func(req *pb.UpdateAppGroupStatusRequest, res *pb.SimpleResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)

				appGroupId, err := uuid.Parse(req.GetAppGroupId())
				require.NoError(t, err)

				appGroup, err := acc.GetAppGroup(appGroupId)
				require.NoError(t, err)

				require.Equal(t, appGroup.GetAppGroupId(), req.GetAppGroupId())
				require.Equal(t, appGroup.GetStatus(), req.GetStatus() )
			},
		},
		{
			name: "INVALID_APPGROUP_ID",
			in:   &pb.UpdateAppGroupStatusRequest{
				AppGroupId: "NO_UUID_STRING",
			},
			checkResponse: func(req *pb.UpdateAppGroupStatusRequest, res *pb.SimpleResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "NO_APPGROUP_ID",
			in:   &pb.UpdateAppGroupStatusRequest{
				AppGroupId: uuid.New().String(),
			},
			checkResponse: func(req *pb.UpdateAppGroupStatusRequest, res *pb.SimpleResponse, err error) {
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

			s := AppInfoServer{}
			res, err := s.UpdateAppGroupStatus(ctx, tc.in)
			tc.checkResponse(tc.in, res, err)
		})
	}
}

func TestUpdateApp(t *testing.T) {
	testCases := []struct {
		name          string
		in            *pb.UpdateAppRequest
		checkResponse func(req *pb.UpdateAppRequest, res *pb.SimpleResponse, err error)
	}{
		{
			name: "OK",
			in:   &pb.UpdateAppRequest{
				AppGroupId : createdAppGroupId,
				AppType : pb.AppType_THANOS,
				Endpoint : "endpoint",
				Metadata : "{\"metadata\":\"no_data\"}",
			},
			checkResponse: func(req *pb.UpdateAppRequest, res *pb.SimpleResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)
			},
		},
		{
			name: "INVALID_APPGROUP_ID",
			in:   &pb.UpdateAppRequest{
				AppGroupId : "THIS_IS_NOT_UUID",
			},
			checkResponse: func(req *pb.UpdateAppRequest, res *pb.SimpleResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "NO_EXIST_APPGROUP_ID",
			in:   &pb.UpdateAppRequest{
				AppGroupId : uuid.New().String(),
			},
			checkResponse: func(req *pb.UpdateAppRequest, res *pb.SimpleResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INTERNAL)
			},
		},
		{
			name: "INVALID_JSON_DATA",
			in:   &pb.UpdateAppRequest{
				AppGroupId : createdAppGroupId,
				AppType : pb.AppType_THANOS,
				Endpoint : "endpoint",
				Metadata : "THIS_IS_NOT_JOSN_DATA",
			},
			checkResponse: func(req *pb.UpdateAppRequest, res *pb.SimpleResponse, err error) {
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

			s := AppInfoServer{}
			res, err := s.UpdateApp(ctx, tc.in)
			tc.checkResponse(tc.in, res, err)
		})
	}
}

func TestGetAppsByAppGroupID(t *testing.T) {
	testCases := []struct {
		name          string
		in            *pb.IDRequest
		checkResponse func(req *pb.IDRequest, res *pb.GetAppsResponse, err error)
	}{
		{
			name: "OK",
			in:   &pb.IDRequest{
				Id: createdAppGroupId,
			},
			checkResponse: func(req *pb.IDRequest, res *pb.GetAppsResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)
			},
		},
		{
			name: "INVALID_APPGROUP_ID",
			in:   &pb.IDRequest{
				Id: "NO_UUID_STRING",
			},
			checkResponse: func(req *pb.IDRequest, res *pb.GetAppsResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "NO_APPGROUP_ID",
			in:   &pb.IDRequest{
				Id: uuid.New().String(),
			},
			checkResponse: func(req *pb.IDRequest, res *pb.GetAppsResponse, err error) {
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

			s := AppInfoServer{}
			res, err := s.GetAppsByAppGroupID(ctx, tc.in)
			tc.checkResponse(tc.in, res, err)
		})
	}
}

func TestDeleteAppGroup(t *testing.T) {
	testCases := []struct {
		name          string
		in            *pb.DeleteAppGroupRequest
		checkResponse func(req *pb.DeleteAppGroupRequest, res *pb.SimpleResponse, err error)
	}{
		{
			name: "OK",
			in:   &pb.DeleteAppGroupRequest{
				AppGroupId: createdAppGroupId,
			},
			checkResponse: func(req *pb.DeleteAppGroupRequest, res *pb.SimpleResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)

				appGroupId, err := uuid.Parse(req.GetAppGroupId())
				require.NoError(t, err)

				appGroup, err := acc.GetAppGroup(appGroupId)
				require.Error(t, err)
				require.True(t, appGroup == nil)
			},
		},
		{
			name: "INVALID_APPGROUP_ID",
			in:   &pb.DeleteAppGroupRequest{
				AppGroupId: "NO_UUID_STRING",
			},
			checkResponse: func(req *pb.DeleteAppGroupRequest, res *pb.SimpleResponse, err error) {
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "NO_APPGROUP_ID",
			in:   &pb.DeleteAppGroupRequest{
				AppGroupId: uuid.New().String(),
			},
			checkResponse: func(req *pb.DeleteAppGroupRequest, res *pb.SimpleResponse, err error) {
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

			s := AppInfoServer{}
			res, err := s.DeleteAppGroup(ctx, tc.in)
			tc.checkResponse(tc.in, res, err)
		})
	}
}

// Helpers

func randomCreateAppGroupRequest() *pb.CreateAppGroupRequest {
	return &pb.CreateAppGroupRequest {
		ClusterId : uuid.New().String(),
		AppGroup : &pb.AppGroup {
			AppGroupId : uuid.New().String(),
			AppGroupName : randomString("APPGROUPNAME"),
			Type : pb.AppGroupType_APP_TYPE_UNSPECIFIED,
			ClusterId : uuid.New().String(),
			Status : pb.AppGroupStatus_APP_GROUP_UNSPECIFIED,
			ExternalLabel: randomString("EXTERNALLABEL"),
		},
	}
}