package cdt

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "cem-cdt-UseCaseSupportUpdate"

	// Setpoints data updated
	//
	// Use `Setpoints` to get the current data
	//
	// Use Case CDT, Scenario 1
	DataUpdateSetpoints api.EventType = "eg-lpc-DataUpdateSetpoints"

	// Setpoint constraints data updated
	//
	// Use `SetpointConstraints` to get the current data
	//
	// Use Case CDT, Scenario 1
	DataUpdateSetpointConstraints api.EventType = "eg-lpc-DataUpdateSetpointConstraints"
)
