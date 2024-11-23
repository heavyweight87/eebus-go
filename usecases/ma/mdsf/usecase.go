package mdsf

import (
	"github.com/enbility/eebus-go/api"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

// MDSF - Monitoring of DHW System Function Use Case
type MDSF struct {
	*usecase.UseCaseBase

	dhwSystemFunctionId             *model.HvacSystemFunctionIdType
	oneTimeDhwOverrunId             *model.HvacOverrunIdType
	operationModeForOperationModeId map[model.HvacOperationModeIdType]model.HvacOperationModeTypeType
}

var _ ucapi.MaMDSFInterface = (*MDSF)(nil)

// Create a new Monitoring of DHW System Function Use Case
func NewMDSF(
	localEntity spineapi.EntityLocalInterface,
	eventCB api.EntityEventCallback,
) *MDSF {
	validActorTypes := []model.UseCaseActorType{
		model.UseCaseActorTypeDHWCircuit,
	}
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeDHWCircuit,
	}
	useCaseScenarios := []api.UseCaseScenario{
		{
			Scenario:  model.UseCaseScenarioSupportType(1),
			Mandatory: true,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeHvac,
			},
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(2),
			Mandatory: true,
			ServerFeatures: []model.FeatureTypeType{
				model.FeatureTypeTypeHvac,
			},
		},
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeMonitoringAppliance,
		model.UseCaseNameTypeMonitoringOfDhwSystemFunction,
		"1.0.0",
		"release",
		useCaseScenarios,
		eventCB,
		UseCaseSupportUpdate,
		validActorTypes,
		validEntityTypes,
	)

	uc := &MDSF{
		UseCaseBase:                     usecase,
		operationModeForOperationModeId: make(map[model.HvacOperationModeIdType]model.HvacOperationModeTypeType),
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *MDSF) AddFeatures() {
	// client features
	var clientFeatures = []model.FeatureTypeType{}
	for _, feature := range clientFeatures {
		_ = e.LocalEntity.GetOrAddFeature(feature, model.RoleTypeClient)
	}
}
