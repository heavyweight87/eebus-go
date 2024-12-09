package api

import (
	"time"

	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
)

type CemOHPCFInterface interface {
	api.UseCaseInterface

	// Scenario 1

	PowerConsumptionFlexibilitySettings(
		entity spineapi.EntityRemoteInterface,
	) (*PowerConsumptionFlexibilitySettings, error)

	State(entity spineapi.EntityRemoteInterface) (*PowerSequenceStateType, error)

	StartTime(entity spineapi.EntityRemoteInterface) (*time.Time, error)

	// Scenario 2

	WriteStartTime(entity spineapi.EntityRemoteInterface, startTime time.Time) error

	Start(entity spineapi.EntityRemoteInterface) error

	Stop(entity spineapi.EntityRemoteInterface) error

	Pause(entity spineapi.EntityRemoteInterface) error
}
