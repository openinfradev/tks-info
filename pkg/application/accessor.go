package application

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/openinfradev/tks-common/pkg/log"
	"github.com/openinfradev/tks-info/pkg/application/model"
	pb "github.com/openinfradev/tks-proto/tks_pb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Accessor is an accessor to postgreSQL to query data.
type Accessor struct {
	db *gorm.DB
}

// New returns new accessor's ptr.
func New(db *gorm.DB) *Accessor {
	return &Accessor{
		db: db,
	}
}

// Create creates a new application group in database.
func (x *Accessor) Create(clusterID string, appGroup *pb.AppGroup) (string, error) {
	existsLabel, err := x.existsExternalLabel(clusterID, appGroup.GetExternalLabel())
	if err != nil {
		return "", err
	}
	if existsLabel {
		return "",
			fmt.Errorf("can't create application group because external label %s already exists",
				appGroup.GetExternalLabel())
	}

	creator := uuid.Nil
	if appGroup.GetCreator() != "" {
		creator, err = uuid.Parse(appGroup.GetCreator())
		if err != nil {
			return "", err
		}
	}

	appGroupModel := model.ApplicationGroup{
		Name:          appGroup.GetAppGroupName(),
		ClusterId:     clusterID,
		Type:          appGroup.GetType(),
		Status:        appGroup.GetStatus(),
		ExternalLabel: appGroup.GetExternalLabel(),
		Creator:       creator,
		Description:   appGroup.GetDescription(),
	}
	res := x.db.Create(&appGroupModel)
	if res.Error != nil {
		return "", res.Error
	}
	return appGroupModel.ID, nil
}

// GetApplicatiionGroup returns an application group in database.
func (x *Accessor) GetAppGroupsByClusterID(clusterID string, offset, limit int) ([]*pb.AppGroup, error) {
	var appGroupModels []model.ApplicationGroup
	res := x.db.Offset(offset).Limit(limit).Where("cluster_id = ?", clusterID).Find(&appGroupModels)
	if res.Error != nil {
		return nil, res.Error
	}

	return reflectToPbAppGroups(appGroupModels), nil
}

// GetAppGroups returns application groups matching name and type in database.
func (x *Accessor) GetAppGroups(name string, appGroupType pb.AppGroupType) ([]*pb.AppGroup, error) {
	var (
		appGroupModels []model.ApplicationGroup
		res            *gorm.DB
	)
	if name == "" && appGroupType == pb.AppGroupType_APP_TYPE_UNSPECIFIED {
		return nil, fmt.Errorf("can't find application groups with empty name and unspecified type")
	}

	if name != "" && appGroupType != pb.AppGroupType_APP_TYPE_UNSPECIFIED {
		res = x.db.Where("name = ? AND type = ?", name, appGroupType).Find(&appGroupModels)
	} else if name == "" && appGroupType != pb.AppGroupType_APP_TYPE_UNSPECIFIED {
		res = x.db.Where("type = ?", appGroupType).Find(&appGroupModels)
	} else {
		res = x.db.Where("name = ?", name).Find(&appGroupModels)
	}
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, fmt.Errorf(
			"could not find application group for name %s, type %d", name, appGroupType)
	}
	return reflectToPbAppGroups(appGroupModels), nil
}

// GetAppGroup returns an application group by cluster_id and app_group_id.
func (x *Accessor) GetAppGroup(appGroupID string) (*pb.AppGroup, error) {
	var appGroupModel model.ApplicationGroup
	res := x.db.First(&appGroupModel, "id = ?", appGroupID)

	if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf(
			"could not find application group for app_group_id %s", appGroupID)
	} else if res.Error != nil {
		return nil, res.Error
	}
	return reflectToPbAppGroup(appGroupModel), nil
}

// UpdateAppGroupStatus updates status of application group.
func (x *Accessor) UpdateAppGroupStatus(appGroupID string, status pb.AppGroupStatus, statusDesc string, workflowId string) error {
	res := x.db.Model(&model.ApplicationGroup{}).
		Where("id = ?", appGroupID).
		Updates(map[string]interface{}{"Status": status, "StatusDesc": statusDesc, "WorkflowId": workflowId})

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("could not update application group status")
	}
	return nil
}

// DeleteAppGroup deletes an application group and applications.
func (x *Accessor) DeleteAppGroup(appGroupID string) error {
	res := x.db.Delete(&model.ApplicationGroup{}, "id = ?", appGroupID)
	log.Info("application group id ", appGroupID, " is deleted!")
	if res.Error != nil || res.RowsAffected == 0 {
		return fmt.Errorf("could not delete application group for app group id %s", appGroupID)
	}
	res = x.db.Delete(model.Application{}, "app_group_id = ?", appGroupID)
	log.Info("deleted applications count: ", res.RowsAffected)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("could not delete application for app_group_id %s", appGroupID)
	}
	return nil
}

// GetAppsByAppGroupID queies applications by app group id.
func (x *Accessor) GetAppsByAppGroupID(appGroupID string) ([]*pb.Application, error) {
	var appModels []model.Application
	res := x.db.Where("app_group_id = ?", appGroupID).Find(&appModels)
	if res.Error != nil || res.RowsAffected == 0 {
		return nil, fmt.Errorf("could not find applications for app group id %s", appGroupID)
	}
	return reflectToPbApplications(appModels), nil
}

// GetApps queies applications by app type.
func (x *Accessor) GetApps(appGroupID string, appType pb.AppType) ([]*pb.Application, error) {
	var appModels []model.Application
	res := x.db.Where("app_group_id = ? AND type = ?", appGroupID, appType).Find(&appModels)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, res.Error
	}
	return reflectToPbApplications(appModels), nil
}

// UpdateApp updates data of application in database.
func (x *Accessor) UpdateApp(appGroupID string, appType pb.AppType, endpoint, metadata string) error {
	res := x.db.Model(&model.Application{}).Where("app_group_id = ? AND type = ?", appGroupID, appType).
		Updates(map[string]interface{}{"endpoint": endpoint, "metadata": metadata})
	if res.Error != nil {
		return res.Error
	} else if res.RowsAffected == 0 {
		if err := x.createApplication(appGroupID, appType, endpoint, metadata); err != nil {
			return err
		}
	}
	return nil
}

func (x *Accessor) createApplication(appGroupID string, appType pb.AppType, endpoint, metadata string) error {
	app := model.Application{
		AppGroupId: appGroupID,
		Type:       appType,
		Endpoint:   endpoint,
		Metadata:   datatypes.JSON([]byte(metadata)),
	}
	res := x.db.Create(&app)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func reflectToPbAppGroups(models []model.ApplicationGroup) []*pb.AppGroup {
	var result []*pb.AppGroup
	for _, model := range models {
		result = append(result, reflectToPbAppGroup(model))
	}
	return result
}

func reflectToPbAppGroup(model model.ApplicationGroup) *pb.AppGroup {
	return &pb.AppGroup{
		AppGroupId:    model.ID,
		AppGroupName:  model.Name,
		Type:          model.Type,
		WorkflowId:    model.WorkflowId,
		Status:        model.Status,
		StatusDesc:    model.StatusDesc,
		ClusterId:     model.ClusterId,
		ExternalLabel: model.ExternalLabel,
		Creator:       model.Creator.String(),
		Description:   model.Description,
		CreatedAt:     timestamppb.New(model.CreatedAt),
		UpdatedAt:     timestamppb.New(model.UpdatedAt),
	}
}

func reflectToPbApplications(models []model.Application) []*pb.Application {
	var result []*pb.Application
	for _, model := range models {
		result = append(result, reflectToPbApplication(model))
	}
	return result
}

func reflectToPbApplication(model model.Application) *pb.Application {
	return &pb.Application{
		AppId:      model.ID.String(),
		AppGroupId: model.AppGroupId,
		Type:       model.Type,
		Endpoint:   model.Endpoint,
		Metadata:   model.Metadata.String(),
		CreatedAt:  timestamppb.New(model.CreatedAt),
		UpdatedAt:  timestamppb.New(model.UpdatedAt),
	}
}

func (x *Accessor) existsExternalLabel(clusterID string, label string) (bool, error) {
	if label == "" {
		return false, nil
	}
	var appGroup model.ApplicationGroup
	res := x.db.First(&appGroup, "cluster_id = ? AND external_label = ?", clusterID, label)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return false, nil
	} else if res.Error != nil {
		return false, res.Error
	}
	return true, nil
}
