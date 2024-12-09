package ohpcf

import (
	"github.com/enbility/eebus-go/features/client"
	"github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *OHPCF) HandleEvent(payload spineapi.EventPayload) {
	// only about events from a compressor entity or device changes for this remote device

	if !e.IsCompatibleEntityType(payload.Entity) {
		return
	}

	if internal.IsEntityConnected(payload) {
		e.connected(payload.Entity)
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.SmartEnergyManagementPsDataType:
		e.SmartEnergyManagementDataUpdate(payload)
	default:
		return
	}
}

func (e *OHPCF) SmartEnergyManagementDataUpdate(payload spineapi.EventPayload) {
	smartEnergyManagementData, ok := payload.Data.(*model.SmartEnergyManagementPsDataType)
	if !ok {
		return
	}

	if smartEnergyManagementData == nil ||
		len(smartEnergyManagementData.Alternatives) == 0 ||
		len(smartEnergyManagementData.Alternatives[0].PowerSequence) == 0 {
		return
	}

	sequence := smartEnergyManagementData.Alternatives[0].PowerSequence[0]

	// Check if the state has changed and notify only if it has
	if sequence.State != nil && sequence.State.State != nil {
		if e.optionalPowerConsumptionState == nil || e.optionalPowerConsumptionState != sequence.State.State {
			e.optionalPowerConsumptionState = sequence.State.State
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateState)
		}
	}

	e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdatePowerConsumptionFlexibilitySettings)
}

func (e *OHPCF) connected(entity spineapi.EntityRemoteInterface) {
	smartEnergyManagement, err := client.NewSmartEnergyManagementPs(e.LocalEntity, entity)
	if err != nil || smartEnergyManagement == nil {
		return
	}

	if !smartEnergyManagement.HasSubscription() {
		if _, err := smartEnergyManagement.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if !smartEnergyManagement.HasBinding() {
		if _, err := smartEnergyManagement.Bind(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if _, err := smartEnergyManagement.RequestData(); err != nil {
		logging.Log().Debug(err)
	}
}
