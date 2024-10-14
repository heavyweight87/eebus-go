package cdt

import (
	"log"

	"github.com/enbility/eebus-go/features/client"
	"github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	"github.com/enbility/ship-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

func (e *CDT) HandleEvent(payload spineapi.EventPayload) {
	if !e.IsCompatibleEntityType(payload.Entity) {
		return
	}

	if internal.IsEntityConnected(payload) {
		e.connected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.SetpointDescriptionListDataType:
		log.Println("Got SetpointDescriptionListDataType")

	case *model.SetpointConstraintsListDataType:
		log.Println("Got SetpointConstraintsListDataType")

	case *model.SetpointListDataType:
		log.Println("Got SetpointListDataType")

	case *model.HvacOperationModeDescriptionListDataType:
		log.Println("Got HvacOperationModeDescriptionListDataType")

	case *model.HvacSystemFunctionSetpointRelationListDataType:
		log.Println("Got HvacSystemFunctionSetpointRelationListDataType")
	}
}

func (e *CDT) connected(entity spineapi.EntityRemoteInterface) {
	if setPoint, err := client.NewSetPoint(e.LocalEntity, entity); err == nil {
		if !setPoint.HasSubscription() {
			if _, err := setPoint.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}
		}

		selector := &model.SetpointDescriptionListDataSelectorsType{
			ScopeType: util.Ptr(model.ScopeTypeTypeDhwTemperature),
		}
		if _, err := setPoint.RequestSetPointDescriptions(selector, nil); err != nil {
			logging.Log().Debug(err)
		}

		if _, err := setPoint.RequestSetPointConstraints(nil, nil); err != nil {
			logging.Log().Debug(err)
		}

		if _, err := setPoint.RequestSetPoints(nil, nil); err != nil {
			logging.Log().Debug(err)
		}
	}

	if hvac, err := client.NewHvac(e.LocalEntity, entity); err == nil {
		if !hvac.HasSubscription() {
			if _, err := hvac.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}
		}

		if _, err := hvac.RequestHvacSystemFunctionSetPointRelations(nil, nil); err != nil {
			logging.Log().Debug(err)
		}

		if _, err := hvac.RequestHvacOperationModeDescriptions(nil, nil); err != nil {
			logging.Log().Debug(err)
		}
	}
}
