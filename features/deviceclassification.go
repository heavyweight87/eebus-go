package features

import (
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

type DeviceClassification struct {
	*FeatureImpl
}

func NewDeviceClassification(localRole, remoteRole model.RoleType, localEntity *spine.EntityLocalImpl, remoteEntity *spine.EntityRemoteImpl) (*DeviceClassification, error) {
	feature, err := NewFeatureImpl(model.FeatureTypeTypeDeviceClassification, localRole, remoteRole, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	dc := &DeviceClassification{
		FeatureImpl: feature,
	}

	return dc, nil
}

// request DeviceClassificationManufacturerData from a remote device entity
func (d *DeviceClassification) RequestManufacturerDetails() (*model.MsgCounterType, error) {
	return d.requestData(model.FunctionTypeDeviceClassificationManufacturerData, nil, nil)
}

// get the current manufacturer details for a remote device entity
func (d *DeviceClassification) GetManufacturerDetails() (*model.DeviceClassificationManufacturerDataType, error) {
	rData := d.featureRemote.DataCopy(model.FunctionTypeDeviceClassificationManufacturerData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.DeviceClassificationManufacturerDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	return data, nil
}
