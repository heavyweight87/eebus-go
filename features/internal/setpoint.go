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

// NewLocalSetPoint creates a new SetPointCommon helper for local entities
func NewLocalSetPoint(featureLocal spineapi.FeatureLocalInterface) *SetPointCommon {
	return &SetPointCommon{
		featureLocal: featureLocal,
	}
}

// NewRemoteSetPoint creates a new SetPointCommon helper for remote entities
func NewRemoteSetPoint(featureRemote spineapi.FeatureRemoteInterface) *SetPointCommon {
	return &SetPointCommon{
		featureRemote: featureRemote,
	}
}

var _ api.SetPointCommonInterface = (*SetPointCommon)(nil)

// GetSetpointForId returns the setpoint data for a given setpoint ID
func (s *SetPointCommon) GetSetpointForId(
	id model.SetpointIdType,
) (*model.SetpointDataType, error) {
	filter := model.SetpointDataType{
		SetpointId: &id,
	}

	result, err := s.GetSetpointDataForFilter(filter)
	if err != nil || len(result) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return util.Ptr(result[0]), nil
}

// GetSetpointConstraintsForId returns the setpoint constraints for a given setpoint ID
func (s *SetPointCommon) GetSetpointConstraintsForId(
	id model.SetpointIdType,
) (*model.SetpointConstraintsDataType, error) {
	filter := model.SetpointConstraintsDataType{
		SetpointId: &id,
	}

	result, err := s.GetSetpointConstraintsForFilter(filter)

	if err != nil || len(result) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return util.Ptr(result[0]), nil
}

// GetSetpointDataForFilter returns the setpoint data for a given filter
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

// GetSetpointConstraintsForFilter returns the setpoint constraints for a given filter
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
