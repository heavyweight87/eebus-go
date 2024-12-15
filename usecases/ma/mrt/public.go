package mrt

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	"github.com/enbility/ship-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Scenario 1

// return the momentary room temperature
//
// Parameters:
//   - entity: the entity of the room
//
// Returns:
//   - the room temperature
//
// possible errors:
//   - ErrDataNotAvailable if room temperature is not available
//   - and others
func (e *MRT) RoomTemperature(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	measurement, err := client.NewMeasurement(e.LocalEntity, entity)
	if err != nil {
		return 0, api.ErrMetadataNotAvailable
	}

	filter := model.MeasurementDescriptionDataType{
		ScopeType: util.Ptr(model.ScopeTypeTypeRoomAirTemperature),
	}

	measurements, err := measurement.GetDataForFilter(filter)
	if err == nil && len(measurements) == 1 {
		return measurements[0].Value.GetValue(), nil
	}

	return 0.0, api.ErrDataNotAvailable
}
