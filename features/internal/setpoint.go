package internal

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
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
