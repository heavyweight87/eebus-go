package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
)

// Actor: Monitoring Appliance
// UseCase: Monitoring of DHW System Function
type MaMDSFInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// Return the current operation mode of the HVAC system
	//
	// Parameters:
	//   - entity: the entity of the e.g. HVAC system
	//
	// Return values:
	//   - mode: the current operation mode
	//
	// Possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	OperationMode(entity spineapi.EntityRemoteInterface) (*HvacOperationModeType, error)

	// Scenario 2

	// Return the current overrun status of the HVAC system
	//
	// Parameters:
	//   - entity: the entity of the e.g. HVAC system
	//
	// Return values:
	//   - status: the current overrun status
	//
	// Possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	OverrunStatus(entity spineapi.EntityRemoteInterface) (*HvacOverrunStatusType, error)
}
