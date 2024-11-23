package mdsf

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "ma-mdsf-UseCaseSupportUpdate"

	// Data update of the operation mode
	//
	// Use `OperationMode` to get the current data
	//
	// Use Case MDSF, Scenario 1
	DataUpdateOperationMode api.EventType = "ma-mdsf-DataUpdateOperationMode"

	// Data update of the overrun status
	//
	// Use `OverrunStatus` to get the current data
	//
	// Use Case MDSF, Scenario 2
	DataUpdateOverrunStatus api.EventType = "ma-mdsf-DataUpdateOverrunStatus"
)
