package internal

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type HvacCommon struct {
	featureLocal  spineapi.FeatureLocalInterface
	featureRemote spineapi.FeatureRemoteInterface
}

// NewLocalHvac creates a new HvacCommon helper for local entities
func NewLocalHvac(featureLocal spineapi.FeatureLocalInterface) *HvacCommon {
	return &HvacCommon{
		featureLocal: featureLocal,
	}
}

// NewRemoteHvac creates a new HvacCommon helper for remote entities
func NewRemoteHvac(featureRemote spineapi.FeatureRemoteInterface) *HvacCommon {
	return &HvacCommon{
		featureRemote: featureRemote,
	}
}

var _ api.HvacCommonInterface = (*HvacCommon)(nil)

// GetHvacOperationModeDescriptions returns the operation mode descriptions
func (h *HvacCommon) GetHvacOperationModeDescriptions() ([]model.HvacOperationModeDescriptionDataType, error) {
	function := model.FunctionTypeHvacOperationModeDescriptionListData
	operationModeDescriptions := make([]model.HvacOperationModeDescriptionDataType, 0)

	data, err := featureDataCopyOfType[model.HvacOperationModeDescriptionListDataType](h.featureLocal, h.featureRemote, function)
	if err == nil || data != nil {
		operationModeDescriptions = append(operationModeDescriptions, data.HvacOperationModeDescriptionData...)
	}

	return operationModeDescriptions, nil
}

// GetHvacSystemFunctionOperationModeRelations returns the operation mode relations (used to map operation modes to setpoints)
func (h *HvacCommon) GetHvacSystemFunctionOperationModeRelations() ([]model.HvacSystemFunctionSetpointRelationDataType, error) {
	function := model.FunctionTypeHvacSystemFunctionSetPointRelationListData
	relations := make([]model.HvacSystemFunctionSetpointRelationDataType, 0)

	data, err := featureDataCopyOfType[model.HvacSystemFunctionSetpointRelationListDataType](h.featureLocal, h.featureRemote, function)
	if err == nil || data != nil {
		relations = append(relations, data.HvacSystemFunctionSetpointRelationData...)
	}

	return relations, nil
}
