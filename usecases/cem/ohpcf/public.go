package ohpcf

import (
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

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

func getPowerValues(data *model.SmartEnergyManagementPsDataType) (float64, float64, error) {
	power := float64(0)
	maxPower := float64(0)

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

	data, err := smartEnergyManagement.GetData()
	if err != nil {
		return nil, err
	}

	if data == nil ||
		len(data.Alternatives) == 0 ||
		len(data.Alternatives[0].PowerSequence) == 0 ||
		data.Alternatives[0].PowerSequence[0].Description.SequenceId == nil ||
		data.Alternatives[0].PowerSequence[0].State == nil ||
		data.Alternatives[0].PowerSequence[0].State.State == nil {
		return nil, api.ErrDataNotAvailable
	}

	isAvailable := false
	if data.Alternatives[0].PowerSequence[0].State.State != nil {
		isAvailable = (*data.Alternatives[0].PowerSequence[0].State.State == model.PowerSequenceStateTypeInactive)
	}

	operatingConstraintsInterrupt := data.Alternatives[0].PowerSequence[0].OperatingConstraintsInterrupt
	isStoppable := (operatingConstraintsInterrupt.IsStoppable != nil && *operatingConstraintsInterrupt.IsStoppable)
	isPausable := (operatingConstraintsInterrupt.IsPausable != nil && *operatingConstraintsInterrupt.IsPausable)
	if !isStoppable && !isPausable {
		logging.Log().Error("Power sequence is neither stoppable nor pausable")
		return nil, api.ErrDataNotAvailable
	}

	earliestStartTime := time.Now()
	latestEndTime := time.Now()
	if data.Alternatives[0].PowerSequence[0].ScheduleConstraints != nil {
		scheduleConstraints := data.Alternatives[0].PowerSequence[0].ScheduleConstraints
		if scheduleConstraints.EarliestStartTime != nil {
			earliestStartTime, err = scheduleConstraints.EarliestStartTime.GetTime()
			if err != nil {
				logging.Log().Error("Earliest start time is missing")
				return nil, err
			}
		}

		if scheduleConstraints.LatestEndTime != nil {
			latestEndTime, err = data.Alternatives[0].PowerSequence[0].ScheduleConstraints.LatestEndTime.GetTime()
			if err != nil {
				logging.Log().Error("Latest end time is missing")
				return nil, err
			}
		}
	}

	// If Element "schedule" is absent, this denotes that this
	// sequence will not be started autonomously by the server.
	canStartAutonomously := data.Alternatives[0].PowerSequence[0].Schedule != nil

	power, maxPower, err := getPowerValues(data)
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

func (e *OHPCF) Start(entity spineapi.EntityRemoteInterface) error {
	flexibilitySettings, _ := e.PowerConsumptionFlexibilitySettings(entity)
	currentState, _ := e.State(entity)

	if flexibilitySettings == nil || currentState == nil {
		return api.ErrDataNotAvailable
	}

	switch *currentState {
	case ucapi.PowerSequenceStateTypeInactive:
		// Schedule (activate) optional power consumption process to begin within 5 seconds
		return e.WriteStartTime(entity, time.Now().Add(time.Second*5))

	case ucapi.PowerSequenceStateTypeRunning:
		// Already running
		return nil

	case ucapi.PowerSequenceStateTypePaused:
		if flexibilitySettings.IsPauseable {
			// Resume the process
			return e.writeState(entity, model.PowerSequenceStateTypeRunning)
		}
	}

	return api.ErrNotSupported
}

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
