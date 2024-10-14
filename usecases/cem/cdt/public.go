package cdt

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	usecasesapi "github.com/enbility/eebus-go/usecases/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Setpoints retrieves the setpoints for various HVAC operation modes from a remote entity.
//
// Possible errors:
//   - ErrDataNotAvailable: If the mapping of operation modes to setpoints or the setpoints themselves are not available.
//   - Other errors: Any other errors encountered during the process.
func (e *CDT) Setpoints(
	entity spineapi.EntityRemoteInterface,
) (map[model.HvacOperationModeTypeType]usecasesapi.Setpoint, error) {
	if len(e.operationModeToSetpoint) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	setPoint, err := client.NewSetPoint(e.LocalEntity, entity)
	if err != nil {
		return nil, err
	}

	setpoints := make(map[model.HvacOperationModeTypeType]usecasesapi.Setpoint)
	for mode, setpointId := range e.operationModeToSetpoint {
		if setpoint, err := setPoint.GetSetpointForId(setpointId); err == nil {
			var value float64 = 0
			var minValue float64 = 0
			var maxValue float64 = 0
			var timePeriod model.TimePeriodType = model.TimePeriodType{}

			if setpoint.Value != nil {
				value = setpoint.Value.GetValue()
			}

			if setpoint.ValueMax != nil {
				maxValue = setpoint.ValueMax.GetValue()
			}

			if setpoint.ValueMin != nil {
				minValue = setpoint.ValueMin.GetValue()
			}

			if setpoint.TimePeriod != nil {
				timePeriod = *setpoint.TimePeriod
			}

			isActive := (setpoint.IsSetpointActive != nil && *setpoint.IsSetpointActive)
			isChangeable := (setpoint.IsSetpointChangeable != nil && *setpoint.IsSetpointChangeable)

			setpoints[mode] = usecasesapi.Setpoint{
				Value:        value,
				MinValue:     minValue,
				MaxValue:     maxValue,
				IsActive:     isActive,
				IsChangeable: isChangeable,
				TimePeriod:   timePeriod,
			}
		}
	}

	return setpoints, nil
}

// SetpointConstraints retrieves the setpoint constraints for various HVAC operation modes from a remote entity.
//
// Possible errors:
//   - ErrDataNotAvailable: If the mapping of operation modes to setpoints or the setpoint constraints are not available.
//   - Other errors: Any other errors encountered during the process.
func (e *CDT) SetpointConstraints(
	entity spineapi.EntityRemoteInterface,
) (map[model.HvacOperationModeTypeType]usecasesapi.SetpointConstraints, error) {
	if len(e.operationModeToSetpoint) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	setPoint, err := client.NewSetPoint(e.LocalEntity, entity)
	if err != nil {
		return nil, err
	}

	setpointConstraints := make(map[model.HvacOperationModeTypeType]usecasesapi.SetpointConstraints)
	for mode, setpointId := range e.operationModeToSetpoint {
		if constraints, err := setPoint.GetSetpointConstraintsForId(setpointId); err == nil {
			var minValue float64 = 0
			var maxValue float64 = 0
			var setSize float64 = 0

			if constraints.SetpointRangeMin != nil {
				minValue = constraints.SetpointRangeMin.GetValue()
			}

			if constraints.SetpointRangeMax != nil {
				maxValue = constraints.SetpointRangeMax.GetValue()
			}

			if constraints.SetpointStepSize != nil {
				setSize = constraints.SetpointStepSize.GetValue()
			}

			setpointConstraints[mode] = usecasesapi.SetpointConstraints{
				MinValue: minValue,
				MaxValue: maxValue,
				StepSize: setSize,
			}
		}
	}

	return setpointConstraints, nil
}

// WriteSetpoint sets the temperature setpoint for a specified HVAC operation mode on a remote entity.
//
// Possible errors:
//   - ErrDataNotAvailable: If the mapping of operation modes to setpoints or the setpoint constraints are not available.
//   - Other errors: Any other errors encountered during the process.
func (e *CDT) WriteSetpoint(
	entity spineapi.EntityRemoteInterface,
	mode model.HvacOperationModeTypeType,
	degC float64,
) error {
	if mode == model.HvacOperationModeTypeTypeAuto {
		return nil
	}

	setpointId, found := e.operationModeToSetpoint[mode]
	if !found {
		return api.ErrDataNotAvailable
	}

	setPoint, err := client.NewSetPoint(e.LocalEntity, entity)
	if err != nil {
		return err
	}

	setpointToWrite := []model.SetpointDataType{
		{
			SetpointId: &setpointId,
			Value:      model.NewScaledNumberType(degC),
		},
	}

	if _, err = setPoint.WriteSetPointListData(setpointToWrite); err != nil {
		return err
	}

	return nil
}
