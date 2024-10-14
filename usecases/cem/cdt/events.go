package cdt

import (
	"github.com/enbility/eebus-go/features/client"
	"github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// HandleEvent handles events for the CDT use case.
func (e *CDT) HandleEvent(payload spineapi.EventPayload) {
	if !e.IsCompatibleEntityType(payload.Entity) {
		return
	}

	if internal.IsEntityConnected(payload) {
		e.dhwCircuitconnected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.SetpointDescriptionListDataType:
		logging.Log().Debug("SetpointDescriptionListDataType")

	case *model.SetpointConstraintsListDataType:
		e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateSetpointConstraints)

	case *model.SetpointListDataType:
		e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateSetpoints)

	case *model.HvacOperationModeDescriptionListDataType,
		*model.HvacSystemFunctionSetpointRelationListDataType:
		e.resolveOpModeToSetpointMapping(payload)
	}
}

// resolveOpModeToSetpointMapping resolves the mapping between operation modes and setpoints.
func (e *CDT) resolveOpModeToSetpointMapping(payload spineapi.EventPayload) {
	hvac, err := client.NewHvac(e.LocalEntity, payload.Entity)
	if err != nil {
		logging.Log().Debug(err)
	}

	// We need both operation mode descriptions and relations to resolve the mapping
	opModeDescriptions, _ := hvac.GetHvacOperationModeDescriptions()
	relations, _ := hvac.GetHvacSystemFunctionOperationModeRelations()
	if len(opModeDescriptions) == 0 || len(relations) == 0 {
		return
	}

	clear(e.operationModeToSetpoint)

	// Create a mapping between operation mode IDs and operation modes
	operationModeIdToOperationMode := make(map[model.HvacOperationModeIdType]model.HvacOperationModeTypeType)
	for _, opModeDescription := range opModeDescriptions {
		modeId := opModeDescription.OperationModeId
		mode := opModeDescription.OperationModeType
		operationModeIdToOperationMode[*modeId] = *mode
	}

	// Create a mapping between operation modes and setpoint IDs
	operationModeToSetpoint := make(map[model.HvacOperationModeTypeType][]model.SetpointIdType)
	for _, relation := range relations {
		mode := operationModeIdToOperationMode[*relation.OperationModeId]
		operationModeToSetpoint[mode] = append(operationModeToSetpoint[mode], relation.SetpointId...)
	}

	for mode, setpointIDs := range operationModeToSetpoint {
		if len(setpointIDs) != 1 {
			if mode == model.HvacOperationModeTypeTypeAuto {
				// For the "auto" operation mode, multiple setpoints (up to four) are allowed as per the specification
				logging.Log().Debugf("Operation mode %s cycles between %d setpoints", mode, len(setpointIDs))
			} else {
				// For other operation modes, having multiple setpoints is not allowed
				// but not explicitly considered an error according to the specification
				logging.Log().Errorf("Operation mode %s has %d setpoint IDs", mode, len(setpointIDs))
			}
			continue
		}

		// Save the unique 1:1 mapping between the operation mode and its corresponding setpoint ID
		e.operationModeToSetpoint[mode] = setpointIDs[0]
	}

	// Now that we have resolved the mapping between operation modes and setpoints,
	// we can request the setpoint descriptions, setpoints, and setpoint constraints
	if setPoint, err := client.NewSetPoint(e.LocalEntity, payload.Entity); err == nil {
		if _, err := setPoint.RequestSetPointDescriptions(nil, nil); err != nil {
			logging.Log().Debug(err)
		}

		if _, err := setPoint.RequestSetPoints(nil, nil); err != nil {
			logging.Log().Debug(err)
		}

		if _, err := setPoint.RequestSetPointConstraints(nil, nil); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// dhwCircuitconnected processes required steps when a DHW Circuit is connected.
func (e *CDT) dhwCircuitconnected(entity spineapi.EntityRemoteInterface) {
	// Request operation mode descriptions and setpoint relationships to map modes to setpoints.
	if hvac, err := client.NewHvac(e.LocalEntity, entity); err == nil {
		if !hvac.HasSubscription() {
			if _, err := hvac.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}
		}

		if _, err := hvac.RequestHvacOperationModeDescriptions(nil, nil); err != nil {
			logging.Log().Debug(err)
		}

		if _, err := hvac.RequestHvacSystemFunctionSetPointRelations(nil, nil); err != nil {
			logging.Log().Debug(err)
		}
	}
}
