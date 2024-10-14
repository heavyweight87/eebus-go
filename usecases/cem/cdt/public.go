package cdt

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	usecasesapi "github.com/enbility/eebus-go/usecases/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

func (e *CDT) Setpoints(
	entity spineapi.EntityRemoteInterface,
) (map[model.HvacOperationModeTypeType]usecasesapi.Setpoint, error) {
	if len(e.operationModeToSetpoint) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	// setpoints := make(map[model.HvacOperationModeTypeType]usecasesapi.Setpoint)
	// for mode, setpointId := range e.operationModeToSetpoint {
	// 	if setPoint, err := client.NewSetPoint(e.LocalEntity, entity); err == nil {
	// 		setpoint, err := setPoint.GetSetpointForId(setpointId)
	// 		if err != nil {
	// 			return nil, err
	// 		}

	// 	}
	// }

	return nil, nil
}

func (e *CDT) SetpointConstraints(
	entity spineapi.EntityRemoteInterface,
) (map[model.HvacOperationModeTypeType]usecasesapi.SetpointConstraints, error) {
	if len(e.operationModeToSetpoint) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	setpointConstraints := make(map[model.HvacOperationModeTypeType]usecasesapi.SetpointConstraints)
	for mode, setpointId := range e.operationModeToSetpoint {
		if setPoint, err := client.NewSetPoint(e.LocalEntity, entity); err == nil {
			constraints, err := setPoint.GetSetpointConstraintsForId(setpointId)
			if err != nil {
				return nil, err
			}

			setpointConstraints[mode] = usecasesapi.SetpointConstraints{
				Min:     constraints.SetpointRangeMin.GetValue(),
				Max:     constraints.SetpointRangeMax.GetValue(),
				SetSize: constraints.SetpointStepSize.GetValue(),
			}
		}
	}

	return setpointConstraints, nil
}

func (e *CDT) WriteSetpoint(
	entity spineapi.EntityRemoteInterface,
	mode model.HvacOperationModeTypeType,
	degC float64,
) error {
	if mode == model.HvacOperationModeTypeTypeAuto {
		return nil
	}

	// if setpointId, ok := e.operationModeToSetpoint[mode]; ok {
	// 	setPoint, err := client.NewSetPoint(e.LocalEntity, entity)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	setpoint, err := setPoint.GetSetpointForId(setpointId)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
