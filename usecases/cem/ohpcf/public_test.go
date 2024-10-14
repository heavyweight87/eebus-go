package ohpcf

import (
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *CemOHPCFSuite) Test_NodeScheduleInformation() {
	data, err := s.sut.SmartEnergyManagementData(nil)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, data.NodeScheduleInformation.AlternativesCount)
	assert.Equal(s.T(), false, data.NodeScheduleInformation.NodeRemoteControllable)
	assert.Equal(s.T(), false, data.NodeScheduleInformation.SupportsReselection)
	assert.Equal(s.T(), false, data.NodeScheduleInformation.SupportsSingleSlotSchedulingOnly)
	assert.Equal(s.T(), 0, data.NodeScheduleInformation.TotalSequencesCountMax)

	data, err = s.sut.SmartEnergyManagementData(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, data.NodeScheduleInformation.AlternativesCount)
	assert.Equal(s.T(), false, data.NodeScheduleInformation.NodeRemoteControllable)
	assert.Equal(s.T(), false, data.NodeScheduleInformation.SupportsReselection)
	assert.Equal(s.T(), false, data.NodeScheduleInformation.SupportsSingleSlotSchedulingOnly)
	assert.Equal(s.T(), 0, data.NodeScheduleInformation.TotalSequencesCountMax)

	limitData := &model.SmartEnergyManagementPsDataType{
		NodeScheduleInformation: &model.PowerSequenceNodeScheduleInformationDataType{

			NodeRemoteControllable:           util.Ptr(true),
			SupportsSingleSlotSchedulingOnly: util.Ptr(true),
			AlternativesCount:                util.Ptr(uint(1)),
			TotalSequencesCountMax:           util.Ptr(uint(3)),
			SupportsReselection:              util.Ptr(false),
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeSmartEnergyManagementPs, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeLoadControlLimitListData, limitData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.SmartEnergyManagementData(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, data.NodeScheduleInformation.AlternativesCount)
	assert.Equal(s.T(), true, data.NodeScheduleInformation.NodeRemoteControllable)
	assert.Equal(s.T(), false, data.NodeScheduleInformation.SupportsReselection)
	assert.Equal(s.T(), true, data.NodeScheduleInformation.SupportsSingleSlotSchedulingOnly)
	assert.Equal(s.T(), 3, data.NodeScheduleInformation.TotalSequencesCountMax)
}
