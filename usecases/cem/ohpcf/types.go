package ohpcf

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "cem-ohpcf-UseCaseSupportUpdate"

	// Returns the current data for the optional power consumption
	//
	// Use `PowerConsumptionFlexibility` to get the current data
	//
	// Use Case OHPCF, Scenario 1
	DataUpdatePowerConsumptionFlexibilitySettings api.EventType = "cem-ohpcf-DataUpdatePowerConsumptionFlexibilitySettings"

	// Returns the current data for the optional power consumption state
	//
	// Use `State` to get the current data
	//
	// Use Case OHPCF, Scenario 1
	DataUpdateState api.EventType = "cem-ohpcf-DataUpdateState"
)
