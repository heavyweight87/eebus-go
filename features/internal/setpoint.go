package internal

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/ship-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type SetPointCommon struct {
	featureLocal  spineapi.FeatureLocalInterface
	featureRemote spineapi.FeatureRemoteInterface
}

func NewLocalSetPoint(featureLocal spineapi.FeatureLocalInterface) *SetPointCommon {
	return &SetPointCommon{
		featureLocal: featureLocal,
	}
}

func NewRemoteSetPoint(featureRemote spineapi.FeatureRemoteInterface) *SetPointCommon {
	return &SetPointCommon{
		featureRemote: featureRemote,
	}
}

var _ api.SetPointCommonInterface = (*SetPointCommon)(nil)

func (s *SetPointCommon) GetSetpointForId(id model.SetpointIdType) (*model.SetpointDataType, error) {
	filter := model.SetpointDataType{
		SetpointId: &id,
	}

	result, err := s.GetSetpointDataForFilter(filter)

	if err != nil || len(result) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return util.Ptr(result[0]), nil
}

func (s *SetPointCommon) GetSetpointConstraintsForId(id model.SetpointIdType) (*model.SetpointConstraintsDataType, error) {
	filter := model.SetpointConstraintsDataType{
		SetpointId: &id,
	}

	result, err := s.GetSetpointConstraintsForFilter(filter)

	if err != nil || len(result) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return util.Ptr(result[0]), nil
}

func (s *SetPointCommon) GetSetpointDataForFilter(
	filter model.SetpointDataType,
) ([]model.SetpointDataType, error) {
	function := model.FunctionTypeSetpointListData

	data, err := featureDataCopyOfType[model.SetpointListDataType](s.featureLocal, s.featureRemote, function)
	if err != nil || data == nil || data.SetpointData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.SetpointDataType](data.SetpointData, filter)

	return result, nil
}

func (s *SetPointCommon) GetSetpoints() ([]model.SetpointDataType, error) {
	function := model.FunctionTypeSetpointListData

	data, err := featureDataCopyOfType[model.SetpointListDataType](s.featureLocal, s.featureRemote, function)
	if err != nil || data == nil || data.SetpointData == nil {
		return nil, api.ErrDataNotAvailable
	}

	return data.SetpointData, nil
}

func (s *SetPointCommon) GetSetpointConstraints() ([]model.SetpointConstraintsDataType, error) {
	function := model.FunctionTypeSetpointConstraintsListData

	data, err := featureDataCopyOfType[model.SetpointConstraintsListDataType](s.featureLocal, s.featureRemote, function)
	if err != nil || data == nil || data.SetpointConstraintsData == nil {
		return nil, api.ErrDataNotAvailable
	}

	return data.SetpointConstraintsData, nil
}

func (s *SetPointCommon) GetSetpointConstraintsForFilter(
	filter model.SetpointConstraintsDataType,
) ([]model.SetpointConstraintsDataType, error) {
	function := model.FunctionTypeSetpointConstraintsListData

	data, err := featureDataCopyOfType[model.SetpointConstraintsListDataType](s.featureLocal, s.featureRemote, function)
	if err != nil || data == nil || data.SetpointConstraintsData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.SetpointConstraintsDataType](data.SetpointConstraintsData, filter)

	return result, nil
}
