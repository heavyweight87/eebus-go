package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type CaCDSFInterface interface {
	api.UseCaseInterface

	// return the operation mode of the DHW system
	//
	// parameters:
	//   - entity: the entity of the e.g. EVSE
	//
	// return values:
	//   - The operation mode of the remote DHW entity
	OperationMode(entity spineapi.EntityRemoteInterface) (model.HvacOperationModeTypeType, error)
}
