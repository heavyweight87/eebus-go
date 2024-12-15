package mrt

import (
	"github.com/enbility/eebus-go/api"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

// MRT - Monitoring of Room Temperature Use Case
type MRT struct {
	*usecase.UseCaseBase
}

var _ ucapi.MaMRTInterface = (*MRT)(nil)

// Create a new Monitoring of Room Temperature Use Case
func NewMRT(
	localEntity spineapi.EntityLocalInterface,
	eventCB api.EntityEventCallback,
) *MRT {
	validActorTypes := []model.UseCaseActorType{
		model.UseCaseActorTypeHVACRoom,
	}
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeHvacRoom,
	}
	useCaseScenarios := []api.UseCaseScenario{
		{
			Scenario:  model.UseCaseScenarioSupportType(1),
			Mandatory: true,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeMeasurement,
			},
		},
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeMonitoringAppliance,
		model.UseCaseNameTypeMonitoringOfRoomTemperature,
		"1.0.0",
		"release",
		useCaseScenarios,
		eventCB,
		UseCaseSupportUpdate,
		validActorTypes,
		validEntityTypes,
	)

	uc := &MRT{
		UseCaseBase: usecase,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *MRT) AddFeatures() {
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeMeasurement,
	}

	for _, feature := range clientFeatures {
		_ = e.LocalEntity.GetOrAddFeature(feature, model.RoleTypeClient)
	}
}
