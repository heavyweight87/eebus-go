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

func (m *HvacCommon) GetOperatingModeDescriptionsForFilter(
	filter model.HvacOperationModeDescriptionDataType,
) ([]model.HvacOperationModeDescriptionDataType, error) {
	function := model.FunctionTypeMeasurementDescriptionListData

	data, err := featureDataCopyOfType[model.HvacOperationModeDescriptionListDataType](m.featureLocal, m.featureRemote, function)
	if err != nil || data == nil || data.HvacOperationModeDescriptionData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.HvacOperationModeDescriptionDataType](data.HvacOperationModeDescriptionData, filter)
	return result, nil
}
