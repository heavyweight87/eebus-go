package ohpcf

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "cem-ohpcf-UseCaseSupportUpdate"

	// Optional heat pump commpressor's consumption data was updated
	//
	// Use `SmartEnergyManagementData` to get the current data
	//
	// Use Case OHPCF, Scenario 1
	DataUpdateSmartEnergyManagementData api.EventType = "cem-ohpcf-DataUpdateSmartEnergyManagementData"
)
