package ohpcf

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Scenario 1 - Monitor heat pump compressor's power consumption flexibility

// Read the current Smart Energy Management Data
//
// parameters:
//   - entity: the entity of the e.g. HVAC
//
// return values:
//   - limit: load limit data
//
// possible errors:
//   - ErrDataNotAvailable if no such limit is (yet) available
//   - and others
func (e *OHPCF) SmartEnergyManagementData(entity spineapi.EntityRemoteInterface) (
	mode *model.SmartEnergyManagementPsDataType, resultErr error) {

	resultErr = api.ErrNoCompatibleEntity
	if !e.IsCompatibleEntityType(entity) {
		return
	}

	resultErr = api.ErrDataNotAvailable
	smartEnergyManagement, err := client.NewSmartEnergyManagementPs(e.LocalEntity, entity)
	if err != nil || smartEnergyManagement == nil {
		return
	}

	smartEnergyManagementData, err := smartEnergyManagement.GetData()
	if err != nil {
		return
	}

	resultErr = nil

	return smartEnergyManagementData, resultErr
}

// Scenario 2 - Control heat pump compressor's power consumption flexibility

// Write the Smart Energy Management Data
//
// parameters:
//   - entity: the entity of the heatpump compressor
//   - value: the new limit in W
func (e *OHPCF) WriteSmartEnergyManagementData(entity spineapi.EntityRemoteInterface,
	data *model.SmartEnergyManagementPsDataType) (*model.MsgCounterType, error) {

	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	smartEnergyManagement, err := client.NewSmartEnergyManagementPs(e.LocalEntity, entity)
	if err != nil || smartEnergyManagement == nil {
		return nil, api.ErrDataNotAvailable
	}

	msgCounter, err := smartEnergyManagement.WriteData(data)
	if err != nil {
		return nil, err
	}

	return msgCounter, nil
}
