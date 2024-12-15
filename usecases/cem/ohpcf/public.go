package ohpcf

import (
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/ship-go/logging"
	"github.com/enbility/ship-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Returns copy of the SmartEnergyManagementPsDataType data for the given entity.
func (e *OHPCF) getCopyOfSmartEnergyManagementData(
	entity spineapi.EntityRemoteInterface,
) (*model.SmartEnergyManagementPsDataType, error) {
	smartEnergyManagement, err := client.NewSmartEnergyManagementPs(e.LocalEntity, entity)
	if err != nil || smartEnergyManagement == nil {
		return nil, api.ErrDataNotAvailable
	}

	data, err := smartEnergyManagement.GetData()
	if err != nil {
		return nil, err
	}

	smartEnergyManagementPsData := model.SmartEnergyManagementPsDataType{}
	util.DeepCopy(data, &smartEnergyManagementPsData)

	return &smartEnergyManagementPsData, nil
}

// Scenario 1 - Monitor heat pump compressor's power consumption flexibility

// getPowerValues extracts the power and max power values from the given smart energy management data.
func getPowerValues(
	data *model.SmartEnergyManagementPsDataType,
) (power float64, maxPower float64, err error) {
	if data.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList == nil ||
		len(data.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList.Value) == 0 ||
		data.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList.Value[0].Value == nil ||
		data.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList.Value[0].ValueType == nil {
		return 0, 0, api.ErrDataNotAvailable
	}

	valueType := *data.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList.Value[0].ValueType
	value := (*data.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList.Value[0].Value).GetValue()

	switch valueType {
	case model.PowerTimeSlotValueTypeTypePower:
		power = value
	case model.PowerTimeSlotValueTypeTypePowerMax:
		maxPower = value
	default:
		return 0, 0, api.ErrDataNotAvailable
	}

	return power, maxPower, nil
}

// PowerConsumptionFlexibilitySettings returns the power consumption flexibility settings for the given entity.
//
// Parameters:
//   - entity: The entity of the heat pump compressor.
//
// Returns:
//   - *ucapi.PowerConsumptionFlexibilitySettings: The power consumption flexibility settings for the given entity.
//
// Possible errors:
//   - api.ErrNoCompatibleEntity: The entity is not compatible with the use case.
//   - api.ErrDataNotAvailable: The data is not (yet) available.
func (e *OHPCF) PowerConsumptionFlexibilitySettings(
	entity spineapi.EntityRemoteInterface,
) (*ucapi.PowerConsumptionFlexibilitySettings, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	smartEnergyManagement, err := client.NewSmartEnergyManagementPs(e.LocalEntity, entity)
	if err != nil || smartEnergyManagement == nil {
		return nil, api.ErrDataNotAvailable
	}

	smartEnergyManagementData, err := e.getCopyOfSmartEnergyManagementData(entity)
	if err != nil {
		return nil, err
	}

	sequence := &smartEnergyManagementData.Alternatives[0].PowerSequence[0]

	isAvailable := (sequence.State.State != nil && *sequence.State.State == model.PowerSequenceStateTypeInactive)
	operatingConstraintsInterrupt := sequence.OperatingConstraintsInterrupt
	isStoppable := (operatingConstraintsInterrupt.IsStoppable != nil && *operatingConstraintsInterrupt.IsStoppable)
	isPausable := (operatingConstraintsInterrupt.IsPausable != nil && *operatingConstraintsInterrupt.IsPausable)
	if !isStoppable && !isPausable {
		logging.Log().Error("Power sequence is neither stoppable nor pausable")
		return nil, api.ErrDataNotAvailable
	}

	earliestStartTime := time.Now()
	latestEndTime := time.Now()
	if sequence.ScheduleConstraints != nil {
		scheduleConstraints := sequence.ScheduleConstraints
		if scheduleConstraints.EarliestStartTime != nil {
			earliestStartTime, err = scheduleConstraints.EarliestStartTime.GetTime()
			if err != nil {
				logging.Log().Error("Earliest start time is missing")
				return nil, err
			}
		}

		if scheduleConstraints.LatestEndTime != nil {
			latestEndTime, err = scheduleConstraints.LatestEndTime.GetTime()
			if err != nil {
				logging.Log().Error("Latest end time is missing")
				return nil, err
			}
		}
	}

	// If Element "schedule" is absent, this denotes that this
	// sequence will not be started autonomously by the server.
	canStartAutonomously := sequence.Schedule != nil

	power, maxPower, err := getPowerValues(smartEnergyManagementData)
	if err != nil || (power == 0 && maxPower == 0) {
		logging.Log().Error("Both power and max power are missing")
		return nil, err
	}

	settings := &ucapi.PowerConsumptionFlexibilitySettings{
		Power:                power,
		MaxPower:             maxPower,
		IsAvailable:          isAvailable,
		IsStoppable:          isStoppable,
		IsPauseable:          isPausable,
		CanStartAutonomously: canStartAutonomously,
		EarliestStartTime:    earliestStartTime,
		LatestEndTime:        latestEndTime,
	}

	return settings, nil
}

func (e *OHPCF) State(entity spineapi.EntityRemoteInterface) (*ucapi.PowerSequenceStateType, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	smartEnergyManagement, err := client.NewSmartEnergyManagementPs(e.LocalEntity, entity)
	if err != nil || smartEnergyManagement == nil {
		return nil, api.ErrDataNotAvailable
	}

	data, err := smartEnergyManagement.GetData()
	if err != nil {
		return nil, err
	}

	if data == nil ||
		len(data.Alternatives) == 0 ||
		len(data.Alternatives[0].PowerSequence) == 0 ||
		data.Alternatives[0].PowerSequence[0].State == nil ||
		data.Alternatives[0].PowerSequence[0].State.State == nil {
		return nil, api.ErrDataNotAvailable
	}

	state := ucapi.PowerSequenceStateType(*data.Alternatives[0].PowerSequence[0].State.State)

	return &state, nil
}

// StartTime returns the start time of the power sequence for the given entity.
//
// Parameters:
//   - entity: The entity of the heat pump compressor.
//
// Returns:
//   - *time.Time: The start time of the power sequence for the given entity.
//
// Possible errors:
//   - api.ErrNoCompatibleEntity: The entity is not compatible with the use case.
//   - api.ErrDataNotAvailable: The data is not (yet) available.
//   - and others
func (e *OHPCF) StartTime(entity spineapi.EntityRemoteInterface) (*time.Time, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	smartEnergyManagement, err := client.NewSmartEnergyManagementPs(e.LocalEntity, entity)
	if err != nil || smartEnergyManagement == nil {
		return nil, api.ErrDataNotAvailable
	}

	smartEnergyManagementPsData, err := e.getCopyOfSmartEnergyManagementData(entity)
	if err != nil {
		return nil, err
	}

	sequence := &smartEnergyManagementPsData.Alternatives[0].PowerSequence[0]
	if sequence.Schedule == nil || sequence.Schedule.StartTime == nil {
		return nil, api.ErrDataNotAvailable
	}

	startTime, err := sequence.Schedule.StartTime.GetTime()
	if err != nil {
		return nil, err
	}

	return &startTime, nil
}

// Scenario 2 - Control heat pump compressor's power consumption flexibility

// WriteStartTime writes the start time for the optional power consumption.
//
// Parameters:
//   - entity: The entity of the heat pump compressor.
//   - startTime: The start time for the power sequence.
//
// Possible errors:
//   - api.ErrNoCompatibleEntity: The entity is not compatible with the use case.
//   - api.ErrDataNotAvailable: The data is not (yet) available.
//   - api.ErrNotSupported: The operation is not supported.
//   - and others
func (e *OHPCF) WriteStartTime(entity spineapi.EntityRemoteInterface, startTime time.Time) error {
	if !e.IsCompatibleEntityType(entity) {
		return api.ErrNoCompatibleEntity
	}

	smartEnergyManagement, err := client.NewSmartEnergyManagementPs(e.LocalEntity, entity)
	if err != nil || smartEnergyManagement == nil {
		return api.ErrDataNotAvailable
	}

	if startTime.Compare(time.Now()) < 0 {
		logging.Log().Error("Start time is in the past")
		return api.ErrNotSupported
	}

	smartEnergyManagementPsData, err := e.getCopyOfSmartEnergyManagementData(entity)
	if err != nil {
		return err
	}

	sequence := &smartEnergyManagementPsData.Alternatives[0].PowerSequence[0]
	sequence.Schedule = &model.PowerSequenceScheduleDataType{
		StartTime: model.NewAbsoluteOrRelativeTimeTypeFromDuration(time.Until(startTime)),
	}

	_, err = smartEnergyManagement.WriteData(smartEnergyManagementPsData)

	return err
}

// writeState updates the power sequence state of a given entity.
//
// Parameters:
//   - entity: The remote entity interface for which the power sequence state needs to be updated.
//   - state: The new power sequence state to be set.
//
// Returns:
//   - error: An error if any step in the process fails, otherwise nil.
func (e *OHPCF) writeState(entity spineapi.EntityRemoteInterface, state model.PowerSequenceStateType) error {
	if !e.IsCompatibleEntityType(entity) {
		return api.ErrNoCompatibleEntity
	}

	smartEnergyManagement, err := client.NewSmartEnergyManagementPs(e.LocalEntity, entity)
	if err != nil || smartEnergyManagement == nil {
		return api.ErrDataNotAvailable
	}

	smartEnergyManagementPsData, err := e.getCopyOfSmartEnergyManagementData(entity)
	if err != nil {
		return err
	}

	sequence := smartEnergyManagementPsData.Alternatives[0].PowerSequence[0]
	sequence.State.State = &state

	_, err = smartEnergyManagement.WriteData(smartEnergyManagementPsData)

	return err
}

// Resume resumes a paused power sequence for the given entity.
//
// Parameters:
//   - entity: The entity for which the power sequence should be resumed.
//
// Returns:
//   - error: An error if the operation is not supported or data is not available.
func (e *OHPCF) Resume(entity spineapi.EntityRemoteInterface) error {
	flexibilitySettings, _ := e.PowerConsumptionFlexibilitySettings(entity)
	currentState, _ := e.State(entity)

	if flexibilitySettings == nil || currentState == nil {
		return api.ErrDataNotAvailable
	}

	if *currentState != ucapi.PowerSequenceStateTypePaused {
		logging.Log().Error("Power sequence is not paused.")
		return api.ErrNotSupported
	}

	return e.writeState(entity, model.PowerSequenceStateTypeRunning)
}

// Stop stops a running power sequence for the given entity.
//
// Parameters:
//   - entity: The entity for which the power sequence should be stopped.
//
// Returns:
//   - error: An error if the operation is not supported or data is not available.
func (e *OHPCF) Stop(entity spineapi.EntityRemoteInterface) error {
	flexibilitySettings, _ := e.PowerConsumptionFlexibilitySettings(entity)
	currentState, _ := e.State(entity)

	if flexibilitySettings == nil || currentState == nil {
		return api.ErrDataNotAvailable
	}

	if !flexibilitySettings.IsStoppable {
		logging.Log().Error("Power sequence is not stoppable")
		return api.ErrNotSupported
	}

	if *currentState != ucapi.PowerSequenceStateTypeRunning {
		logging.Log().Error("Power sequence is not running. Nothing to stop")
		return api.ErrNotSupported
	}

	return e.writeState(entity, model.PowerSequenceStateTypeInvalid)
}

// Pause pauses a running power sequence for the given entity.
//
// Parameters:
//   - entity: The entity for which the power sequence should be paused.
//
// Returns:
//   - error: An error if the operation is not supported or data is not available.
func (e *OHPCF) Pause(entity spineapi.EntityRemoteInterface) error {
	flexibilitySettings, _ := e.PowerConsumptionFlexibilitySettings(entity)
	currentState, _ := e.State(entity)

	if flexibilitySettings == nil || currentState == nil {
		return api.ErrDataNotAvailable
	}

	if !flexibilitySettings.IsPauseable {
		logging.Log().Error("Power sequence is not pausable")
		return api.ErrNotSupported
	}

	if *currentState != ucapi.PowerSequenceStateTypeRunning {
		logging.Log().Error("Power sequence is not running. Nothing to pause")
		return api.ErrNotSupported
	}

	return e.writeState(entity, model.PowerSequenceStateTypePaused)
}
