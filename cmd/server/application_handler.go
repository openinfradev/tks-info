package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/openinfradev/tks-common/pkg/log"
	"github.com/openinfradev/tks-info/pkg/application"
	app "github.com/openinfradev/tks-info/pkg/application"
	pb "github.com/openinfradev/tks-proto/tks_pb"
)

var acc *app.Accessor

type AppInfoServer struct {
	pb.UnimplementedAppInfoServiceServer
}

func InitAppInfoHandler(db *gorm.DB) {
	acc = application.New(db)
}

func (s *AppInfoServer) CreateAppGroup(ctx context.Context, in *pb.CreateAppGroupRequest) (*pb.IDResponse, error) {
	clusterID, err := uuid.Parse(in.GetClusterId())
	if err != nil {
		res := pb.IDResponse{
			Code: pb.Code_INVALID_ARGUMENT,
			Error: &pb.Error{
				Msg: fmt.Sprintf("invalid cluster ID %s", in.GetClusterId()),
			},
		}
		return &res, err
	}
	log.Info("Request 'CreateAppGroup' for cluster id ", clusterID)
	appGroup := in.GetAppGroup()

	id, err := acc.Create(clusterID, appGroup)
	if err != nil {
		return &pb.IDResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}

	res := &pb.IDResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
		Id:    id.String(),
	}
	return res, nil
}

func (s *AppInfoServer) GetAppGroupsByClusterID(ctx context.Context, in *pb.IDRequest) (*pb.GetAppGroupsResponse, error) {
	clusterID, err := uuid.Parse(in.GetId())
	if err != nil {
		return &pb.GetAppGroupsResponse{
			Code: pb.Code_INVALID_ARGUMENT,
			Error: &pb.Error{
				Msg: fmt.Sprintf("invalid app group ID %s", in.GetId()),
			},
		}, err
	}
	log.Info("GetAppGroupsByClusterID request for clusterId: ", clusterID)

	appGroups, err := acc.GetAppGroupsByClusterID(clusterID, 0, 10)
	if err != nil {
		return &pb.GetAppGroupsResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}

	return &pb.GetAppGroupsResponse{
		Code:      pb.Code_OK_UNSPECIFIED,
		Error:     nil,
		AppGroups: appGroups,
	}, err
}

func (s *AppInfoServer) GetAppGroups(ctx context.Context, in *pb.GetAppGroupsRequest) (*pb.GetAppGroupsResponse, error) {
	if in.GetAppGroupName() == "" && in.GetType() == pb.AppGroupType_APP_TYPE_UNSPECIFIED {
		err := fmt.Errorf("not efficient conditions to query app group.")
		return &pb.GetAppGroupsResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}
	log.Info("GetAppGroups request for app name: ", in.GetAppGroupName())

	appGroups, err := acc.GetAppGroups(in.GetAppGroupName(), in.GetType())
	if err != nil {
		return &pb.GetAppGroupsResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}

	res := &pb.GetAppGroupsResponse{
		Code:      pb.Code_OK_UNSPECIFIED,
		Error:     nil,
		AppGroups: appGroups,
	}
	return res, nil
}

func (s *AppInfoServer) GetAppGroup(ctx context.Context, in *pb.GetAppGroupRequest) (*pb.GetAppGroupResponse, error) {
	appGroupID, err := uuid.Parse(in.GetAppGroupId())
	if err != nil {
		return &pb.GetAppGroupResponse{
			Code: pb.Code_INVALID_ARGUMENT,
			Error: &pb.Error{
				Msg: fmt.Sprintf("invalid app group ID %s", in.GetAppGroupId()),
			},
		}, err
	}
	log.Info("GetAppGroup request for app group ID: ", appGroupID)
	appGroup, err := acc.GetAppGroup(appGroupID)
	if err != nil {
		return &pb.GetAppGroupResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}

	return &pb.GetAppGroupResponse{
		Code:     pb.Code_OK_UNSPECIFIED,
		Error:    nil,
		AppGroup: appGroup,
	}, nil

}

func (*AppInfoServer) UpdateAppGroupStatus(ctx context.Context, in *pb.UpdateAppGroupStatusRequest) (*pb.SimpleResponse, error) {
	appGroupID, err := uuid.Parse(in.GetAppGroupId())
	if err != nil {
		return &pb.SimpleResponse{
			Code: pb.Code_INVALID_ARGUMENT,
			Error: &pb.Error{
				Msg: fmt.Sprintf("invalid app group ID %s", in.GetAppGroupId()),
			},
		}, err
	}
	log.Info("UpdateAppGroupStatus request for app group ID: ", appGroupID)
	if err = acc.UpdateAppGroupStatus(appGroupID, in.GetStatus()); err != nil {
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

func (s *AppInfoServer) DeleteAppGroup(ctx context.Context, in *pb.DeleteAppGroupRequest) (*pb.SimpleResponse, error) {
	appGroupID, err := uuid.Parse(in.GetAppGroupId())
	if err != nil {
		return &pb.SimpleResponse{
			Code: pb.Code_INVALID_ARGUMENT,
			Error: &pb.Error{
				Msg: fmt.Sprintf("invalid app group ID %s", in.GetAppGroupId()),
			},
		}, err
	}
	log.Info("DeleteAppGroup request for app group ID: ", appGroupID)
	if err = acc.DeleteAppGroup(appGroupID); err != nil {
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

func (*AppInfoServer) GetAppsByAppGroupID(ctx context.Context, in *pb.IDRequest) (*pb.GetAppsResponse, error) {
	appGroupID, err := uuid.Parse(in.GetId())
	if err != nil {
		return &pb.GetAppsResponse{
			Code: pb.Code_INVALID_ARGUMENT,
			Error: &pb.Error{
				Msg: fmt.Sprintf("invalid app group ID %s", in.GetId()),
			},
		}, err
	}
	log.Info("GetAppsByAppGroupID request for app group ID: ", appGroupID)
	apps, err := acc.GetAppsByAppGroupID(appGroupID)
	if err != nil {
		return &pb.GetAppsResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}
	return &pb.GetAppsResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
		Apps:  apps,
	}, nil
}

func (*AppInfoServer) GetApps(ctx context.Context, in *pb.GetAppsRequest) (*pb.GetAppsResponse, error) {
	appGroupID, err := uuid.Parse(in.GetAppGroupId())
	if err != nil {
		return &pb.GetAppsResponse{
			Code: pb.Code_INVALID_ARGUMENT,
			Error: &pb.Error{
				Msg: fmt.Sprintf("invalid app group ID %s", in.GetAppGroupId()),
			},
		}, err
	}
	log.Info("GetApps request for app group ID: ", appGroupID)
	apps, err := acc.GetApps(appGroupID, in.GetType())
	if err != nil {
		return &pb.GetAppsResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}
	return &pb.GetAppsResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
		Apps:  apps,
	}, nil
}

func (*AppInfoServer) UpdateApp(ctx context.Context, in *pb.UpdateAppRequest) (*pb.SimpleResponse, error) {
	appGroupID, err := uuid.Parse(in.GetAppGroupId())
	if err != nil {
		return &pb.SimpleResponse{
			Code: pb.Code_INVALID_ARGUMENT,
			Error: &pb.Error{
				Msg: fmt.Sprintf("invalid app group ID %s", in.GetAppGroupId()),
			},
		}, err
	}
	log.Info("UpdateApp request for app group ID: ", appGroupID)
	log.Info(">>> endpoint: ", in.GetEndpoint())
	if err = acc.UpdateApp(appGroupID, in.GetAppType(), in.GetEndpoint(), in.GetMetadata()); err != nil {
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
