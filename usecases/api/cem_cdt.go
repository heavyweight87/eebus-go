package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type CemCDTInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// Return the current setpoints data
	//
	// parameters:
	//   - entity: the entity to get the setpoints data from
	//
	// return values:
	//   - setpoints: A map of the setpoints for supported modes
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	Setpoints(entity spineapi.EntityRemoteInterface) (map[model.HvacOperationModeTypeType]Setpoint, error)

	// Return the constraints for the setpoints
	//
	// parameters:
	//   - entity: the entity to get the setpoints constraints from
	//
	// return values:
	//   - setpointConstraints: A map of the constraints for supported modes
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	SetpointConstraints(entity spineapi.EntityRemoteInterface) (map[model.HvacOperationModeTypeType]SetpointConstraints, error)

	// Write a setpoint
	//
	// parameters:
	//   - entity: the entity to write the setpoint to
	//   - mode: the mode to write the setpoint for
	//   - degC: the setpoint value to write
	WriteSetpoint(entity spineapi.EntityRemoteInterface, mode model.HvacOperationModeTypeType, degC float64) error
}
