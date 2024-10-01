package client

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type Hvac struct {
	*Feature

	*internal.HvacCommon
}

// Get a new LoadControl features helper
//
// - The feature on the local entity has to be of role client
// - The feature on the remote entity has to be of role server
func NewHvac(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface) (*Hvac, error) {
	feature, err := NewFeature(model.FeatureTypeTypeLoadControl, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	hvac := &Hvac{
		Feature:    feature,
		HvacCommon: internal.NewRemoteHvac(feature.featureRemote),
	}

	return hvac, nil
}

var _ api.LoadControlClientInterface = (*LoadControl)(nil)

// return current values for Hvac System Function Descriptions
func (h *Hvac) GetSystemFunctionDescriptions(
	selector *model.HvacSystemFunctionDescriptionListDataSelectorsType,
	elements *model.HvacSystemFunctionDescriptionDataElementsType,
) (*model.MsgCounterType, error) {
	return h.requestData(model.FunctionTypeHvacSystemFunctionDescriptionListData, selector, elements)
}
