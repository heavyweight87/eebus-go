package mrt

import (
	"github.com/enbility/eebus-go/features/client"
	"github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	"github.com/enbility/ship-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// HandleEvent handles events for the MRT use case
func (e *MRT) HandleEvent(payload spineapi.EventPayload) {
	if !e.IsCompatibleEntityType(payload.Entity) {
		return
	}

	if internal.IsEntityConnected(payload) {
		e.deviceConnected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.MeasurementDescriptionListDataType:
		e.deviceMeasurementDescriptionDataUpdate(payload)

	case *model.MeasurementConstraintsListDataType:
		e.measurementConstraintsDataUpdate(payload)

	case *model.MeasurementListDataType:
		e.deviceMeasurementDataUpdate(payload)
	}
}

func (e *MRT) measurementConstraintsDataUpdate(payload spineapi.EventPayload) {
	if measurement, err := client.NewMeasurement(e.LocalEntity, payload.Entity); err == nil {
		// Scenario 1
		filter := model.MeasurementDescriptionDataType{
			ScopeType: util.Ptr(model.ScopeTypeTypeRoomAirTemperature),
		}
		if measurement.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateRoomTemperature)
		}
	}
}

func (e *MRT) deviceConnected(entity spineapi.EntityRemoteInterface) {
	if measurement, err := client.NewMeasurement(e.LocalEntity, entity); err == nil {
		if !measurement.HasSubscription() {
			if _, err := measurement.Subscribe(); err != nil {
				logging.Log().Error(err)
			}
		}

		selector := &model.MeasurementDescriptionListDataSelectorsType{
			ScopeType: util.Ptr(model.ScopeTypeTypeRoomAirTemperature),
		}
		if _, err := measurement.RequestDescriptions(selector, nil); err != nil {
			logging.Log().Error(err)
		}
	}
}

func (e *MRT) deviceMeasurementDescriptionDataUpdate(payload spineapi.EventPayload) {
	measurement, err := client.NewMeasurement(e.LocalEntity, payload.Entity)
	if err != nil {
		logging.Log().Errorf("Error creating measurement client: %v", err)
		return
	}

	filter := model.MeasurementDescriptionDataType{
		ScopeType: util.Ptr(model.ScopeTypeTypeRoomAirTemperature),
	}
	measurements, err := measurement.GetDataForFilter(filter)
	if err != nil || len(measurements) == 0 {
		logging.Log().Errorf("Error getting measurement data for filter: %v", err)
		return
	}

	if measurements[0].MeasurementId == nil {
		logging.Log().Error("Measurement ID is nil")
		return
	}

	// Preform partial reads for measurement constrains and measurements data with
	// measurementId from the measurement description of the room air temperature.
	measurementId := *measurements[0].MeasurementId

	measurementsSelector := &model.MeasurementListDataSelectorsType{
		MeasurementId: &measurementId,
	}
	if _, err := measurement.RequestData(measurementsSelector, nil); err != nil {
		logging.Log().Error("Error getting measurement list values:", err)
	}

	constraintsSelector := &model.MeasurementConstraintsListDataSelectorsType{
		MeasurementId: &measurementId,
	}
	if _, err := measurement.RequestConstraints(constraintsSelector, nil); err != nil {
		logging.Log().Error(err)
	}
}

func (e *MRT) deviceMeasurementDataUpdate(payload spineapi.EventPayload) {
	if measurement, err := client.NewMeasurement(e.LocalEntity, payload.Entity); err == nil {
		// Scenario 1
		filter := model.MeasurementDescriptionDataType{
			ScopeType: util.Ptr(model.ScopeTypeTypeRoomAirTemperature),
		}
		if measurement.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateRoomTemperature)
		}
	}
}
