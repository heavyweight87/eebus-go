package ohpcf

import (
	"github.com/enbility/eebus-go/features/client"
	"github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
)

// handle SPINE events
func (e *OHPCF) HandleEvent(payload spineapi.EventPayload) {
	// only about events from an EV entity or device changes for this remote device

	if !e.IsCompatibleEntityType(payload.Entity) {
		return
	}

	if internal.IsEntityConnected(payload) {
		// get the smart energy management data from the remote entity
		smartEnergyManagement, err := client.NewSmartEnergyManagementPs(e.LocalEntity, payload.Entity)
		if err != nil || smartEnergyManagement == nil {
			return
		}
		if _, err := smartEnergyManagement.RequestData(); err != nil {
			logging.Log().Debug(err)
		}
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	default:
		return
	}
}
