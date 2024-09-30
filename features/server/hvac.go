package server

import (
	"github.com/enbility/eebus-go/features/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

type Hvac struct {
	*Feature
	*internal.HvacCommon
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

func (h *Hvac) AddOperatingModeDescription(
	description model.HvacOperationModeDescriptionDataType,
) *model.HvacOperationModeIdType {
	if description.OperationModeId != nil {
		return nil
	}

	data, err := h.GetOperatingModeDescriptionsForFilter(model.HvacOperationModeDescriptionDataType{})
	if err != nil {
		data = []model.HvacOperationModeDescriptionDataType{}
	}

	maxId := model.HvacOperationModeIdType(0)

	for _, item := range data {
		if item.OperationModeId != nil && *item.OperationModeId >= maxId {
			maxId = *item.OperationModeId + 1
		}
	}

	operationModeId := util.Ptr(maxId)
	description.OperationModeId = operationModeId

	partial := model.NewFilterTypePartial()
	datalist := &model.HvacOperationModeDescriptionListDataType{
		HvacOperationModeDescriptionData: []model.HvacOperationModeDescriptionDataType{description},
	}

	if err := h.featureLocal.UpdateData(model.FunctionTypeLoadControlLimitDescriptionListData, datalist, partial, nil); err != nil {
		return nil
	}

	return operationModeId
}
