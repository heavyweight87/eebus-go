package cdsf

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	"github.com/enbility/ship-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Scenario 1

// return the current operation mode
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
func (e *OHPCF) OperationMode(entity spineapi.EntityRemoteInterface) (
	mode model.HvacOperationModeTypeType, resultErr error) {

	resultErr = api.ErrNoCompatibleEntity
	if !e.IsCompatibleEntityType(entity) {
		return
	}

	resultErr = api.ErrDataNotAvailable
	hvac, err := client.NewHvac(e.LocalEntity, entity)
	if err != nil || hvac == nil {
		return
	}

	filter := model.HvacOperationModeDescriptionDataType{
		OperationModeType: util.Ptr(model.HvacOperationModeTypeTypeOff),
	}
	limitDescriptions, err := hvac.GetOperatingModeDescriptionsForFilter(filter)
	if err != nil || len(limitDescriptions) != 1 {
		return
	}

	resultErr = nil

	return
}
