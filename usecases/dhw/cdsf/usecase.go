package cdsf

import (
	"github.com/enbility/eebus-go/api"

	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

// Configuration of Domestic Water Heater System Function
type CDSF struct {
	*usecase.UseCaseBase
}

var _ ucapi.DhwCDSFInterface = (*CDSF)(nil)

func NewCDSF(
	localEntity spineapi.EntityLocalInterface,
	eventCB api.EntityEventCallback,
) *CDSF {
	validActorTypes := []model.UseCaseActorType{
		model.UseCaseActorTypeDHWCircuit,
		model.UseCaseActorTypeConfigurationAppliance,
	}
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeDHWCircuit,
	}
	useCaseScenarios := []api.UseCaseScenario{
		{
			Scenario:  model.UseCaseScenarioSupportType(1),
			Mandatory: true,
		},
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeConfigurationOfDhwSystemFunction,
		"1.0.1",
		"release",
		useCaseScenarios,
		eventCB,
		UseCaseSupportUpdate,
		validActorTypes,
		validEntityTypes,
	)

	uc := &CDSF{
		UseCaseBase: usecase,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *CDSF) AddFeatures() {
	// server features
	f := e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeHvac, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeHvacSystemFunctionListData, true, false)
	f.AddFunctionType(model.FunctionTypeHvacOperationModeDescriptionListData, true, false)

	/*	if hvac, err := server.NewHvac(e.LocalEntity, e.EventCB); err == nil {
		operationModeDesc := model.HvacOperationModeDescriptionDataType{
			OperationModeType: util.Ptr(model.HvacOperationModeTypeType),
		}
	}*/

}
