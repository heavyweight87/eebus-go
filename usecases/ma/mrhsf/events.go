package mrhsf

import (
	"github.com/enbility/eebus-go/features/client"
	"github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	"github.com/enbility/ship-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// HandleEvent handles events for the MRHSF use case
func (e *MRHSF) HandleEvent(payload spineapi.EventPayload) {
	if !e.IsCompatibleEntityType(payload.Entity) {
		return
	}

	if internal.IsEntityConnected(payload) {
		e.connected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.HvacSystemFunctionDescriptionListDataType:
		e.updateSystemFunctionDescriptions(payload)
	case *model.HvacSystemFunctionOperationModeRelationListDataType:
		e.updateSystemFunctionOperationModeRelations(payload)
	case *model.HvacOperationModeDescriptionListDataType,
		*model.HvacSystemFunctionListDataType:
		e.updateOperationModes(payload)
	}
}

func (e *MRHSF) updateOperationModes(payload spineapi.EventPayload) {
	if e.heatingSystemFunctionID == nil {
		return
	}

	hvac, err := client.NewHvac(e.LocalEntity, payload.Entity)
	if err != nil {
		logging.Log().Debug(err)
		return
	}

	relations, _ := hvac.GetHvacSystemFunctionOperationModeRelations()
	descriptions, _ := hvac.GetHvacOperationModeDescriptions()
	functions, _ := hvac.GetHvacSystemFunctions()

	if len(descriptions) == 0 || len(functions) == 0 {
		return
	}

	clear(e.operationModeForOperationModeId)

	opModeIds := make(map[model.HvacOperationModeIdType]model.HvacOperationModeTypeType)
	for _, description := range descriptions {
		opModeIds[*description.OperationModeId] = *description.OperationModeType
	}

	for _, relation := range relations {
		for _, operationModeId := range relation.OperationModeId {
			if operationMode, found := opModeIds[operationModeId]; found {
				e.operationModeForOperationModeId[operationModeId] = operationMode
			}
		}
	}

	e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateOperationMode)
}

func (e *MRHSF) updateSystemFunctionOperationModeRelations(payload spineapi.EventPayload) {
	hvac, err := client.NewHvac(e.LocalEntity, payload.Entity)
	if err != nil {
		logging.Log().Errorf("Failed to create HVAC client: %v", err)
		return
	}

	data := payload.Data.(*model.HvacSystemFunctionOperationModeRelationListDataType)
	if data == nil || len(data.HvacSystemFunctionOperationModeRelationData) == 0 {
		return
	}

	relations := data.HvacSystemFunctionOperationModeRelationData
	if len(relations) != 1 ||
		relations[0].SystemFunctionId == nil ||
		*relations[0].SystemFunctionId != *e.heatingSystemFunctionID {
		return
	}

	modeIds := relations[0].OperationModeId
	for _, id := range modeIds {
		// Request the operation mode description by the operation mode ID for the heating function
		selector := &model.HvacOperationModeDescriptionListDataSelectorsType{
			OperationModeId: &id,
		}
		if _, err := hvac.RequestHvacOperationModeDescriptions(selector, nil); err != nil {
			logging.Log().Errorf("Failed to request HVAC operation mode descriptions: %v", err)
		}
	}

	// Request the system functions for the "heating" system function
	selector := &model.HvacSystemFunctionListDataSelectorsType{
		SystemFunctionId: e.heatingSystemFunctionID,
	}
	if _, err := hvac.RequestHvacSystemFunctions(selector, nil); err != nil {
		logging.Log().Errorf("Failed to request HVAC system functions: %v", err)
	}
}

func (e *MRHSF) updateSystemFunctionDescriptions(payload spineapi.EventPayload) {
	data := payload.Data.(*model.HvacSystemFunctionDescriptionListDataType)
	if data == nil {
		logging.Log().Errorf("Received nil HVAC system function description list data change event")
		return
	}

	hvac, err := client.NewHvac(e.LocalEntity, payload.Entity)
	if err != nil {
		logging.Log().Errorf("Failed to create HVAC client: %v", err)
		return
	}

	// Get the system function ID for the heating function
	filter := model.HvacSystemFunctionDescriptionDataType{
		SystemFunctionType: util.Ptr(model.HvacSystemFunctionTypeTypeHeating),
	}
	systemFunctions, err := hvac.GetHvacSystemFunctionDescriptionsForFilter(filter)
	if err != nil || len(systemFunctions) != 1 {
		logging.Log().Errorf("Failed to get heating function: %v", err)
		return
	}

	// Get the system function ID for the heating function
	e.heatingSystemFunctionID = systemFunctions[0].SystemFunctionId

	// Get the operation mode relations for the heating function
	selector := &model.HvacSystemFunctionOperationModeRelationListDataSelectorsType{
		SystemFunctionId: e.heatingSystemFunctionID,
	}
	if _, err := hvac.RequestHvacSystemFunctionOperationModeRelations(selector, nil); err != nil {
		logging.Log().Errorf("Failed to request HVAC system function operation mode relations: %v", err)
	}
}

func (e *MRHSF) connected(entity spineapi.EntityRemoteInterface) {
	hvac, err := client.NewHvac(e.LocalEntity, entity)
	if err != nil {
		logging.Log().Errorf("Failed to create HVAC client: %v", err)
		return
	}

	if !hvac.HasSubscription() {
		if _, err := hvac.Subscribe(); err != nil {
			logging.Log().Errorf("Failed to subscribe to HVAC: %v", err)
			return
		}
	}

	if _, err := hvac.RequestHvacSystemFunctionDescriptions(nil, nil); err != nil {
		logging.Log().Errorf("Failed to request HVAC system function descriptions: %v", err)
		return
	}
}
