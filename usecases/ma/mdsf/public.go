package mdsf

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/ship-go/logging"
	"github.com/enbility/ship-go/util"
	spineapi "github.com/enbility/spine-go/api"
)

// OperationMode returns the current operation mode.
//
// Parameters:
//   - entity: The entity to get the operation mode for.
//
// Possible errors:
//   - ErrDataNotAvailable: If the operation mode is not (yet) available.
//   - Other: Any other errors encountered during the process.
func (c *MDSF) OperationMode(
	entity spineapi.EntityRemoteInterface,
) (*ucapi.HvacOperationModeType, error) {
	hvac, err := client.NewHvac(c.LocalEntity, entity)
	if err != nil {
		return nil, err
	}

	systemFunctionData, err := hvac.GetHvacSystemFunctionForId(*c.dhwSystemFunctionId)
	if err != nil {
		logging.Log().Debug(err)
		return nil, err
	}

	currentOperationModeId := systemFunctionData.CurrentOperationModeId
	if currentOperationModeId == nil {
		return nil, api.ErrDataNotAvailable
	}

	currentOperationMode, found := c.operationModeForOperationModeId[*currentOperationModeId]
	if !found {
		return nil, api.ErrDataNotAvailable
	}

	return util.Ptr(ucapi.HvacOperationModeType(currentOperationMode)), nil
}

// Returns the overrun status.
//
// Parameters:
//   - entity: The entity to get the overrun status for.
//
// Possible errors:
//   - ErrDataNotAvailable: If the overrun status is not (yet) available.
//   - Other errors: Any other errors encountered during the process.
func (c *MDSF) OverrunStatus(
	entity spineapi.EntityRemoteInterface,
) (*ucapi.HvacOverrunStatusType, error) {
	hvac, err := client.NewHvac(c.LocalEntity, entity)
	if err != nil {
		return nil, err
	}

	if c.oneTimeDhwOverrunId == nil {
		return nil, api.ErrDataNotAvailable
	}

	overrun, err := hvac.GetHvacOverrunForId(*c.oneTimeDhwOverrunId)
	if err != nil {
		return nil, err
	}

	if overrun == nil || overrun.OverrunStatus == nil {
		return nil, api.ErrDataNotAvailable
	}

	return util.Ptr(ucapi.HvacOverrunStatusType(*overrun.OverrunStatus)), nil
}
