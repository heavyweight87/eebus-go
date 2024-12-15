package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
)

// Actor: Monitoring Appliance
// UseCase: Monitoring of Room Heating System Function
type MaMRHSFInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// Return the momentary operation mode
	//
	// Parameters:
	//   - entity: the entity of the device (e.g. HVAC)
	//
	// Possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	OperationMode(entity spineapi.EntityRemoteInterface) (*HvacOperationModeType, error)
}
