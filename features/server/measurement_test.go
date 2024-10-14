package server_test

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/server"
	"github.com/enbility/eebus-go/mocks"
	"github.com/enbility/eebus-go/service"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/cert"
	spineapi "github.com/enbility/spine-go/api"
	spinemocks "github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestMeasurementSuite(t *testing.T) {
	suite.Run(t, new(MeasurementSuite))
}

type MeasurementSuite struct {
	suite.Suite

	sut *server.Measurement

	service api.ServiceInterface

	localEntity spineapi.EntityLocalInterface

	remoteDevice     spineapi.DeviceRemoteInterface
	remoteEntity     spineapi.EntityRemoteInterface
	mockRemoteEntity *spinemocks.EntityRemoteInterface
}

func (s *MeasurementSuite) BeforeTest(suiteName, testName string) {
	cert, _ := cert.CreateCertificate("test", "test", "DE", "test")
	configuration, _ := api.NewConfiguration(
		"test", "test", "test", "test",
		[]shipapi.DeviceCategoryType{shipapi.DeviceCategoryTypeEnergyManagementSystem},
		model.DeviceTypeTypeEnergyManagementSystem,
		[]model.EntityTypeType{model.EntityTypeTypeCEM},
		9999, cert, time.Second*4)

	serviceHandler := mocks.NewServiceReaderInterface(s.T())
	serviceHandler.EXPECT().ServicePairingDetailUpdate(mock.Anything, mock.Anything).Return().Maybe()

	s.service = service.NewService(configuration, serviceHandler)
	_ = s.service.Setup()
	s.localEntity = s.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	mockRemoteDevice := spinemocks.NewDeviceRemoteInterface(s.T())
	s.mockRemoteEntity = spinemocks.NewEntityRemoteInterface(s.T())
	mockRemoteFeature := spinemocks.NewFeatureRemoteInterface(s.T())
	mockRemoteDevice.EXPECT().FeatureByEntityTypeAndRole(mock.Anything, mock.Anything, mock.Anything).Return(mockRemoteFeature).Maybe()
	mockRemoteDevice.EXPECT().Ski().Return(remoteSki).Maybe()
	s.mockRemoteEntity.EXPECT().Device().Return(mockRemoteDevice).Maybe()
	s.mockRemoteEntity.EXPECT().EntityType().Return(mock.Anything).Maybe()
	entityAddress := &model.EntityAddressType{}
	s.mockRemoteEntity.EXPECT().Address().Return(entityAddress).Maybe()
	mockRemoteFeature.EXPECT().DataCopy(mock.Anything).Return(mock.Anything).Maybe()

	var entities []spineapi.EntityRemoteInterface

	s.remoteDevice, entities = setupFeatures(s.service, s.T())
	s.remoteEntity = entities[1]

	var err error
	s.sut, err = server.NewMeasurement(nil)
	assert.NotNil(s.T(), err)

	s.sut, err = server.NewMeasurement(s.localEntity)
	assert.Nil(s.T(), err)
}

func (s *MeasurementSuite) Test_Description() {
	data, err := s.sut.GetDescriptionForId(model.MeasurementIdType(100))
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	desc := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
	}
	measId1 := s.sut.AddDescription(desc)
	assert.NotNil(s.T(), measId1)

	data, err = s.sut.GetDescriptionForId(*measId1)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	desc = model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrentA),
	}

	measId2 := s.sut.AddDescription(desc)
	assert.NotNil(s.T(), measId2)

	limitData, err := s.sut.GetDescriptionForId(*measId2)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), limitData)
}

func (s *MeasurementSuite) Test_GetDescriptionsForFilter() {
	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeStateOfCharge),
	}

	data, err := s.sut.GetDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))

	feature := s.localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeMeasurement, model.RoleTypeServer)

	desc := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: filter.MeasurementType,
				ScopeType:       filter.ScopeType,
			},
		},
	}
	feature.SetData(model.FunctionTypeMeasurementDescriptionListData, desc)

	data, err = s.sut.GetDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))

	filter = model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
	}

	data, err = s.sut.GetDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))
}

func (s *MeasurementSuite) Test_GetLimitData() {
	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeStateOfCharge),
	}

	data := model.MeasurementDataType{}
	err := s.sut.UpdateDataForFilter(data, nil, filter)
	assert.NotNil(s.T(), err)

	data = model.MeasurementDataType{
		MeasurementId: util.Ptr(model.MeasurementIdType(100)),
	}
	err = s.sut.UpdateDataForFilter(data, nil, filter)
	assert.NotNil(s.T(), err)

	mId := s.sut.AddDescription(filter)
	assert.NotNil(s.T(), mId)

	descData, err := s.sut.GetDescriptionForId(*mId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), descData)

	result, err := s.sut.GetDataForId(*mId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), result)

	data = model.MeasurementDataType{
		MeasurementId: mId,
		ValueType:     util.Ptr(model.MeasurementValueTypeTypeValue),
		Value:         model.NewScaledNumberType(16),
	}
	err = s.sut.UpdateDataForFilter(data, nil, filter)
	assert.Nil(s.T(), err)

	result, err = s.sut.GetDataForId(*mId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), result)

	result, err = s.sut.GetDataForId(model.MeasurementIdType(100))
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), result)

	data = model.MeasurementDataType{}
	deleteElements := &model.MeasurementDataElementsType{
		Value: util.Ptr(model.ElementTagType{}),
	}
	err = s.sut.UpdateDataForId(data, deleteElements, *mId)
	assert.Nil(s.T(), err)

	result, err = s.sut.GetDataForId(*mId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Nil(s.T(), result.Value)
}
