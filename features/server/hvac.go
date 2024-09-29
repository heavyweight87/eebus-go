package server

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type Hvac struct {
	*Feature
}

func NewHvac(localEntity spineapi.EntityLocalInterface) (*Hvac, error) {
	feature, err := NewFeature(model.FeatureTypeTypeHvac, localEntity)
	if err != nil {
		return nil, err
	}

	lc := &Hvac{
		Feature: feature,
	}

	return lc, nil
}

func (h *Hvac) AddDescription(
	description model.LoadControlLimitDescriptionDataType,
) *model.LoadControlLimitIdType {
	return nil
}
