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

func NewLocalHvac(featureLocal spineapi.FeatureLocalInterface) *HvacCommon {
	return &HvacCommon{
		featureLocal: featureLocal,
	}
}

func NewRemoteHvac(featureRemote spineapi.FeatureRemoteInterface) *HvacCommon {
	return &HvacCommon{
		featureRemote: featureRemote,
	}
}

var _ api.HvacCommonInterface = (*HvacCommon)(nil)

func (h *HvacCommon) GetHvacOperationModeDescriptions() ([]model.HvacOperationModeDescriptionDataType, error) {
	function := model.FunctionTypeHvacOperationModeDescriptionListData
	operationModeDescriptions := make([]model.HvacOperationModeDescriptionDataType, 0)

	data, err := featureDataCopyOfType[model.HvacOperationModeDescriptionListDataType](h.featureLocal, h.featureRemote, function)
	if err != nil || data == nil || data.HvacOperationModeDescriptionData == nil {
		return make([]model.HvacOperationModeDescriptionDataType, 0), api.ErrDataNotAvailable
	}

	operationModeDescriptions = append(operationModeDescriptions, data.HvacOperationModeDescriptionData...)

	return operationModeDescriptions, nil
}

func (h *HvacCommon) GetHvacSystemFunctionOperationModeRelations() ([]model.HvacSystemFunctionSetpointRelationDataType, error) {
	function := model.FunctionTypeHvacSystemFunctionSetPointRelationListData
	relations := make([]model.HvacSystemFunctionSetpointRelationDataType, 0)

	data, err := featureDataCopyOfType[model.HvacSystemFunctionSetpointRelationListDataType](h.featureLocal, h.featureRemote, function)
	if err != nil || data == nil || data.HvacSystemFunctionSetpointRelationData == nil {
		return make([]model.HvacSystemFunctionSetpointRelationDataType, 0), api.ErrDataNotAvailable
	}

	relations = append(relations, data.HvacSystemFunctionSetpointRelationData...)

	return relations, nil
}
