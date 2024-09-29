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

func (m *HvacCommon) GetDescriptionsForFilter(
	filter model.MeasurementDescriptionDataType,
) ([]model.MeasurementDescriptionDataType, error) {
	function := model.FunctionTypeMeasurementDescriptionListData

	data, err := featureDataCopyOfType[model.MeasurementDescriptionListDataType](m.featureLocal, m.featureRemote, function)
	if err != nil || data == nil || data.MeasurementDescriptionData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.MeasurementDescriptionDataType](data.MeasurementDescriptionData, filter)
	return result, nil
}
