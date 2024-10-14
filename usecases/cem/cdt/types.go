package cdt

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "cem-cdt-UseCaseSupportUpdate"

	// setpoint data updated
	//
	// Use `Setpoint` to get the current data
	//
	// Use Case CDT, Scenario 1
	DataUpdatesetpoint api.EventType = "eg-lpc-DataUpdatesetpoint"

	// setpoint constraints data updated
	//
	// Use `SetpointConstraints` to get the current data
	//
	// Use Case CDT, Scenario 1
	DataUpdatesetpointConstraints api.EventType = "eg-lpc-DataUpdatesetpointConstraints"
)
