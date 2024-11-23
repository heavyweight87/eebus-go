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
		e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateSmartEnergyManagementData)
	}
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

	if _, err := smartEnergyManagement.RequestData(); err != nil {
		logging.Log().Debug(err)
	}
}
