package internal

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
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
