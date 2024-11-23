package mdsf

import (
	"github.com/enbility/eebus-go/features/client"
	internal "github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	"github.com/enbility/ship-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *MDSF) HandleEvent(payload spineapi.EventPayload) {
	if !e.IsCompatibleEntityType(payload.Entity) {
		return
	}

	if internal.IsEntityConnected(payload) {
		e.dhwCircuitConnected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.HvacSystemFunctionDescriptionListDataType:
		e.systemFunctionsDescriptionsUpdate(payload)

	case *model.HvacSystemFunctionOperationModeRelationListDataType:
		e.systemFunctionOperationModeRelationsUpdate(payload)

	case *model.HvacOperationModeDescriptionListDataType,
		*model.HvacSystemFunctionListDataType:
		e.updateOperationModes(payload)

	case *model.HvacOverrunDescriptionListDataType:
		e.updateOverrunDescriptions(payload)

	case *model.HvacOverrunListDataType:
		e.updateOverruns(payload)
	}
}

func (e *MDSF) updateOverruns(payload spineapi.EventPayload) {
	if e.oneTimeDhwOverrunId == nil {
		logging.Log().Error("OneTimeDhw overrun ID is not set")
		return
	}

	overruns, _ := payload.Data.(*model.HvacOverrunListDataType)
	for _, overrun := range overruns.HvacOverrunData {
		if overrun.OverrunId == nil || *overrun.OverrunId == *e.oneTimeDhwOverrunId {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateOverrunStatus)
		}
	}
}

func (e *MDSF) updateOverrunDescriptions(payload spineapi.EventPayload) {
	hvac, err := client.NewHvac(e.LocalEntity, payload.Entity)
	if err != nil {
		logging.Log().Debug(err)
		return
	}

	filter := model.HvacOverrunDescriptionDataType{
		OverrunType: util.Ptr(model.HvacOverrunTypeTypeOneTimeDhw),
	}
	descriptions, err := hvac.GetHvacOverrunDescriptionsForFilter(filter)
	if err != nil {
		logging.Log().Debug(err)
		return
	}
	if len(descriptions) != 1 {
		logging.Log().Errorf("Expected exactly one overrun description for oneTimeDhw overrun type. Got %d", len(descriptions))
		return
	}

	// Store the overrun ID for the "oneTimeDhw" overrun type.
	e.oneTimeDhwOverrunId = descriptions[0].OverrunId

	selector := &model.HvacOverrunListDataSelectorsType{
		OverrunId: e.oneTimeDhwOverrunId,
	}
	if _, err := hvac.RequestHvacOverruns(selector, nil); err != nil {
		logging.Log().Debug(err)
	}
}

func (e *MDSF) updateOperationModes(payload spineapi.EventPayload) {
	hvac, err := client.NewHvac(e.LocalEntity, payload.Entity)
	if err != nil {
		logging.Log().Debug(err)
		return
	}

	relations, _ := hvac.GetHvacSystemFunctionOperationModeRelations()
	descriptions, _ := hvac.GetHvacOperationModeDescriptions()
	systemFunctions, _ := hvac.GetHvacSystemFunctions()

	if len(descriptions) == 0 || len(systemFunctions) == 0 {
		return
	}

	clear(e.operationModeForOperationModeId)

	opModeForOpModeId := make(map[model.HvacOperationModeIdType]model.HvacOperationModeTypeType)
	for _, description := range descriptions {
		opModeForOpModeId[*description.OperationModeId] = *description.OperationModeType
	}

	for _, relation := range relations {
		for _, operationModeId := range relation.OperationModeId {
			if operationMode, found := opModeForOpModeId[operationModeId]; found {
				e.operationModeForOperationModeId[operationModeId] = operationMode
			}
		}
	}

	e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateOperationMode)
}

func (e *MDSF) systemFunctionOperationModeRelationsUpdate(payload spineapi.EventPayload) {
	hvac, err := client.NewHvac(e.LocalEntity, payload.Entity)
	if err != nil {
		logging.Log().Debug(err)
		return
	}

	relations, _ := payload.Data.(*model.HvacSystemFunctionOperationModeRelationListDataType)
	for _, relation := range relations.HvacSystemFunctionOperationModeRelationData {
		for _, operationModeId := range relation.OperationModeId {
			selector := &model.HvacOperationModeDescriptionListDataSelectorsType{
				OperationModeId: &operationModeId,
			}
			if _, err := hvac.RequestHvacOperationModeDescriptions(selector, nil); err != nil {
				logging.Log().Debug(err)
			}
		}
	}
}

func (e *MDSF) systemFunctionsDescriptionsUpdate(payload spineapi.EventPayload) {
	if hvac, err := client.NewHvac(e.LocalEntity, payload.Entity); err == nil {
		filter := model.HvacSystemFunctionDescriptionDataType{
			SystemFunctionType: util.Ptr(model.HvacSystemFunctionTypeTypeDhw),
		}
		descriptions, err := hvac.GetHvacSystemFunctionDescriptionsForFilter(filter)
		if err != nil {
			logging.Log().Debug(err)
			return
		}

		systemFunctionIds := make([]model.HvacSystemFunctionIdType, 0)
		for _, description := range descriptions {
			systemFunctionIds = append(systemFunctionIds, *description.SystemFunctionId)
		}
		if len(systemFunctionIds) != 1 {
			logging.Log().Errorf("Expected exactly one DHW system function description. Got %d", len(systemFunctionIds))
			return
		}

		// Store the system function ID for the DHW system function.
		e.dhwSystemFunctionId = &systemFunctionIds[0]

		// According to the usecase specification, we need to perform a partial read of hvacSystemFunctionOperationModeRelationListData
		// using systemFunctionId derived from hvacSystemFunctionDescriptionListData.
		systemFunctionOperationModeRelationsSelector := &model.HvacSystemFunctionOperationModeRelationListDataSelectorsType{
			SystemFunctionId: e.dhwSystemFunctionId,
		}
		if _, err := hvac.RequestHvacSystemFunctionOperationModeRelations(systemFunctionOperationModeRelationsSelector, nil); err != nil {
			logging.Log().Debug(err)
		}

		// According to the specification, we need to perform a partial read of hvacSystemFunctionListData
		// using systemFunctionId from hvacSystemFunctionDescriptionListData.
		systemFunctionsSelector := &model.HvacSystemFunctionListDataSelectorsType{
			SystemFunctionId: e.dhwSystemFunctionId,
		}
		if _, err := hvac.RequestHvacSystemFunctions(systemFunctionsSelector, nil); err != nil {
			logging.Log().Debug(err)
		}

		// The hvacOverrunDescriptionListData read SHOULD be a "full" read operation.
		_, err = hvac.RequestHvacOverrunDescriptions(nil, nil)
		if err != nil {
			logging.Log().Debug(err)
		}
	}
}

func (e *MDSF) dhwCircuitConnected(entity spineapi.EntityRemoteInterface) {
	if hvac, err := client.NewHvac(e.LocalEntity, entity); err == nil {
		if !hvac.HasSubscription() {
			if _, err := hvac.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}
		}

		if _, err := hvac.RequestHvacSystemFunctionDescriptions(nil, nil); err != nil {
			logging.Log().Debug(err)
		}
	}
}
