package app_serve_app

import (
	//"errors"
	"fmt"

	"github.com/google/uuid"
	//"github.com/openinfradev/tks-common/pkg/log"
	"github.com/openinfradev/tks-info/pkg/app_serve_app/model"
	pb "github.com/openinfradev/tks-proto/tks_pb"
	"google.golang.org/protobuf/types/known/timestamppb"
	//"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Accessor is an accessor to postgreSQL to query data.
type AsaAccessor struct {
	db *gorm.DB
}

// New returns new accessor's ptr.
func New(db *gorm.DB) *AsaAccessor {
	return &AsaAccessor{
		db: db,
	}
}

// Create creates a new application group in database.
func (x *AsaAccessor) Create(contractId string, app *pb.AppServeApp) (uuid.UUID, error) {
	asaModel := model.AppServeApp{
		Name:          app.GetName(),
		ContractId:    contractId,
		Version:       app.GetVersion(),
		TaskType:      app.GetTaskType(),
		Status:        app.GetStatus(),
		ArtifactUrl:   app.GetArtifactUrl(),
		ImageUrl:      app.GetImageUrl(),
		TargetClusterId: app.GetTargetClusterId(),
		Profile:       app.GetProfile(),
	}
	res := x.db.Create(&asaModel)
	if res.Error != nil {
		return uuid.Nil, res.Error
	}
	return asaModel.ID, nil
}

func (x *AsaAccessor) GetAppServeApps(contractId string) ([]*pb.AppServeApp, error) {
	var appServeApps []model.AppServeApp
	pbAppServeApps := []*pb.AppServeApp{}

	res := x.db.Find(&appServeApps, "contract_id = ?", contractId)
	if res.Error != nil {
		return nil, fmt.Errorf("Error while finding appServeApps with contractID: %s", contractId)
	}

	// If no record is found, just return empty array.
	if res.RowsAffected == 0 {
		return pbAppServeApps, nil
	}

	for _, asa := range appServeApps {
		pbAppServeApps = append(pbAppServeApps, ConvertToPbAppServeApp(asa))
	}
	return pbAppServeApps, nil
}

func (x *AsaAccessor) GetAppServeApp(id uuid.UUID, contractId string) (*pb.AppServeApp, error) {
	var appServeApp model.AppServeApp
	res := x.db.First(&appServeApp, "id = ? AND contract_id = ?", id, contractId)
	if res.RowsAffected == 0 || res.Error != nil {
		return &pb.AppServeApp{}, fmt.Errorf("Could not find AppServeApp with ID: %s", id)
	}

	pbAppServeApp := ConvertToPbAppServeApp(appServeApp)
	return pbAppServeApp, nil
}

func (x *AsaAccessor) UpdateStatus(id uuid.UUID, status string, output string) error {
	res := x.db.Model(&model.AppServeApp{}).Where("ID = ?", id).Updates(model.AppServeApp{Status: status, Output: output})

	if res.Error != nil || res.RowsAffected == 0 {
		return fmt.Errorf("UpdateStatus: nothing updated in AppServeApp with id %s", id)
	}

	return nil
}

func (x *AsaAccessor) UpdateEndpoint(id uuid.UUID, endpoint string) error {
	res := x.db.Model(&model.AppServeApp{}).Where("ID = ?", id).Update("EndpointUrl", endpoint)

	if res.Error != nil || res.RowsAffected == 0 {
		return fmt.Errorf("UpdateEndpoint: nothing updated in AppServeApp with id %s", id)
	}

	return nil
}

func ConvertToPbAppServeApp(asa model.AppServeApp) *pb.AppServeApp {
	return &pb.AppServeApp{
		Id:            asa.ID.String(),
		Name:          asa.Name,
		ContractId:    asa.ContractId,
		Version:       asa.Version,
		TaskType:      asa.TaskType,
		Status:        asa.Status,
		Output:        asa.Output,
		ImageUrl:      asa.ImageUrl,
		ArtifactUrl:   asa.ArtifactUrl,
		EndpointUrl:   asa.EndpointUrl,
		TargetClusterId: asa.TargetClusterId,
		Profile:       asa.Profile,
		CreatedAt:     timestamppb.New(asa.CreatedAt),
		UpdatedAt:     timestamppb.New(asa.UpdatedAt),
	}
}
