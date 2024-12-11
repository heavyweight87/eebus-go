package api

import (
	"time"

	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
)

type CemOHPCFInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// Return the optional power consumption flexibility settings of the heat pump compressor
	//
	// parameters:
	//   - entity: the entity of the heat pump compressor
	//
	// returns:
	//   - the optional power consumption flexibility settings of the heat pump compressor
	//
	// possible errors:
	//   - ErrDataNotAvailable if no data is (yet) available
	//   - and others
	PowerConsumptionFlexibilitySettings(
		entity spineapi.EntityRemoteInterface,
	) (*PowerConsumptionFlexibilitySettings, error)

	// Return the power sequence state
	//
	// parameters:
	//   - entity: the entity of the heat pump compressor
	//
	// returns:
	//   - the power sequence state
	//
	// possible errors:
	//   - ErrDataNotAvailable if no data is (yet) available
	//   - and others
	State(entity spineapi.EntityRemoteInterface) (*PowerSequenceStateType, error)

	// Return the start time of the optional power consumption process
	//
	// parameters:
	//   - entity: the entity of the heat pump compressor
	//
	// Returns:
	//  - the start time of the optional power consumption process
	//
	// possible errors:
	//   - ErrDataNotAvailable if no data is (yet) available
	//   - and others
	StartTime(entity spineapi.EntityRemoteInterface) (*time.Time, error)

	// Scenario 2

	// Write the start time of the optional power consumption process
	//
	// parameters:
	//  - entity: the entity of the heat pump compressor
	//  - startTime: the start time of the optional power consumption process
	//
	// possible errors:
	//   - ErrDataNotAvailable if no data is (yet) available
	//   - and others
	WriteStartTime(entity spineapi.EntityRemoteInterface, startTime time.Time) error

	// Resumes the paused optional power consumption prcess (if currently paused).
	//
	// Parameters:
	//  - entity: the entity of the heat pump compressor
	//
	// possible errors:
	//   - ErrDataNotAvailable if no data is (yet) available
	//   - and others
	Resume(entity spineapi.EntityRemoteInterface) error

	// Stops the optional power consumption process (if stoppable).
	//
	// Parameters:
	//  - entity: the entity of the heat pump compressor
	//
	// possible errors:
	//   - ErrDataNotAvailable if no data is (yet) available
	//   - and others
	Stop(entity spineapi.EntityRemoteInterface) error

	// Pauses the optional power consumption process (if pausable).
	//
	// Parameters:
	//  - entity: the entity of the heat pump compressor
	//
	// possible errors:
	//   - ErrDataNotAvailable if no data is (yet) available
	//   - and others
	Pause(entity spineapi.EntityRemoteInterface) error
}
