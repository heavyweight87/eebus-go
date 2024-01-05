package ship

import (
	"time"

	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/ship/model"
	"github.com/enbility/eebus-go/util"
)

// Handshake Hello covers the states smeHello...

// SME_HELLO_STATE_READY_INIT
func (c *ShipConnectionImpl) handshakeHello_Init() {
	if err := c.handshakeHelloSend(model.ConnectionHelloPhaseTypeReady, tHelloInit, false); err != nil {
		c.setAndHandleState(SmeHelloStateAbort)
		return
	}

	c.setState(SmeHelloStateReadyListen, nil)
}

// SME_HELLO_STATE_READY_LISTEN
func (c *ShipConnectionImpl) handshakeHello_ReadyListen(timeout bool, message []byte) {
	if timeout {
		c.handshakeHello_ReadyTimeout()
		return
	}

	var helloReturnMsg model.ConnectionHello
	if err := c.processShipJsonMessage(message, &helloReturnMsg); err != nil {
		c.setAndHandleState(SmeHelloStateAbort)
		return
	}

	hello := helloReturnMsg.ConnectionHello

	switch hello.Phase {
	case model.ConnectionHelloPhaseTypeReady:
		// HELLO_OK
		c.setState(SmeHelloStateOk, nil)

	case model.ConnectionHelloPhaseTypePending:
		// the phase is still pending an no prolongationRequest is set, ignore the message
		if hello.ProlongationRequest == nil {
			return
		}

		// if we got a prolongation request, accept it
		if *hello.ProlongationRequest {
			if c.serviceDataProvider.AllowWaitingForTrust(c.remoteShipID) {
				// re-init timer
				c.setHandshakeTimer(timeoutTimerTypeWaitForReady, tHelloInit)
			}

			if err := c.handshakeHelloSend(model.ConnectionHelloPhaseTypeReady, tHelloInit, false); err != nil {
				c.endHandshakeWithError(err)
			}

			return
		}

		// TODO: what to do if this is false?

	case model.ConnectionHelloPhaseTypeAborted:
		c.setAndHandleState(SmeHelloStateRemoteAbortDone)

		return

	default:
		// don't accept any other responses
		logging.Log().Errorf("Unexpected connection hello phase: %s", hello.Phase)
		c.setAndHandleState(SmeHelloStateAbort)
		return
	}

	c.handleState(false, nil)
}

func (c *ShipConnectionImpl) handshakeHello_ReadyTimeout() {
	c.setAndHandleState(SmeHelloStateAbort)
}

// SME_HELLO_ABORT
func (c *ShipConnectionImpl) handshakeHello_Abort() {
	c.stopHandshakeTimer()

	if err := c.handshakeHelloSend(model.ConnectionHelloPhaseTypeAborted, 0, false); err != nil {
		c.endHandshakeWithError(err)
		return
	}

	c.setAndHandleState(SmeHelloStateAbortDone)
}

// SME_HELLO_PENDING_INIT
func (c *ShipConnectionImpl) handshakeHello_PendingInit() {
	if err := c.handshakeHelloSend(model.ConnectionHelloPhaseTypePending, tHelloInit, false); err != nil {
		c.endHandshakeWithError(err)
		return
	}

	c.setState(SmeHelloStatePendingListen, nil)

	if !c.serviceDataProvider.AllowWaitingForTrust(c.remoteShipID) {
		c.setAndHandleState(SmeHelloStateAbort)
	}
}

// SME_HELLO_PENDING_LISTEN
func (c *ShipConnectionImpl) handshakeHello_PendingListen(timeout bool, message []byte) {
	if timeout {
		// The device needs to be in a state for the user to allow trusting the device
		// e.g. either the web UI or by other means
		if !c.serviceDataProvider.AllowWaitingForTrust(c.remoteShipID) {
			c.handshakeHello_PendingTimeout()
		} else {
			c.handshakeHello_PendingProlongationRequest()
		}

		return
	}

	var helloReturnMsg model.ConnectionHello
	if err := c.processShipJsonMessage(message, &helloReturnMsg); err != nil {
		c.setAndHandleState(SmeHelloStateAbort)
		return
	}

	hello := helloReturnMsg.ConnectionHello

	switch hello.Phase {
	case model.ConnectionHelloPhaseTypeReady:
		if hello.Waiting == nil {
			c.setAndHandleState(SmeHelloStateAbort)
			return
		}

		c.stopHandshakeTimer()

		newDuration := time.Duration(*hello.Waiting) * time.Millisecond
		duration := tHelloProlongThrInc
		if newDuration >= duration {
			// the duration has to be reduced
			duration = newDuration - duration

			// check if it is less than T_hello_prolong_min
			if newDuration >= tHelloProlongMin {
				c.setHandshakeTimer(timeoutTimerTypeSendProlongationRequest, duration)
				return
			}
		}

		if newDuration < tHelloProlongMin {
			// I interpret 13.4.4.1.3 Page 64 Line 1550-1553 as this resulting in a timeout state
			// TODO: verify this
			c.setAndHandleState(SmeHelloStateAbort)
		}

	case model.ConnectionHelloPhaseTypePending:
		if hello.Waiting != nil && hello.ProlongationRequest == nil {
			c.stopHandshakeTimer()

			newDuration := time.Duration(*hello.Waiting) * time.Millisecond
			c.lastReceivedWaitingValue = newDuration
			duration := tHelloProlongThrInc
			if newDuration >= duration {
				// the duration has to be reduced
				duration = newDuration - duration

				// check if it is less than T_hello_prolong_min
				if newDuration >= tHelloProlongMin {
					c.setHandshakeTimer(timeoutTimerTypeSendProlongationRequest, duration)
					return
				}
			}

			if newDuration < tHelloProlongMin {
				// I interpret 13.4.4.1.3 Page 64 Line 1557-1560 as this resulting in a timeout state
				// TODO: verify this
				c.setAndHandleState(SmeHelloStateAbort)
			}

			return
		}

		if hello.Waiting == nil && hello.ProlongationRequest != nil && *hello.ProlongationRequest {
			// if we got a prolongation request, accept it
			if err := c.handshakeHelloSend(model.ConnectionHelloPhaseTypePending, tHelloInit, false); err != nil {
				c.endHandshakeWithError(err)
			}

			return
		}

		c.setAndHandleState(SmeHelloStateAbort)

	case model.ConnectionHelloPhaseTypeAborted:
		c.setAndHandleState(SmeHelloStateRemoteAbortDone)
		return

	default:
		// don't accept any other responses
		logging.Log().Errorf("Unexpected connection hello phase: %s", hello.Phase)
		c.setAndHandleState(SmeHelloStateAbort)
		return
	}

	c.handleState(false, nil)
}

func (c *ShipConnectionImpl) handshakeHello_PendingProlongationRequest() {
	if err := c.handshakeHelloSend(model.ConnectionHelloPhaseTypePending, 0, true); err != nil {
		c.endHandshakeWithError(err)
		return
	}

	// TODO: we need to set the timer to the last received waiting value
	c.setHandshakeTimer(timeoutTimerTypeProlongRequestReply, tHelloInit)
}

func (c *ShipConnectionImpl) handshakeHello_PendingTimeout() {
	if c.getHandshakeTimerType() != timeoutTimerTypeSendProlongationRequest {
		c.setAndHandleState(SmeHelloStateAbort)
		return
	}

	if err := c.handshakeHelloSend(model.ConnectionHelloPhaseTypePending, 0, true); err != nil {
		c.endHandshakeWithError(err)
		return
	}

	if c.lastReceivedWaitingValue == 0 {
		newValue := float64(tHelloInit.Milliseconds()) * 1.1
		c.lastReceivedWaitingValue = time.Duration(newValue)
	}
	c.setHandshakeTimer(timeoutTimerTypeProlongRequestReply, c.lastReceivedWaitingValue)
}

func (c *ShipConnectionImpl) handshakeHelloSend(phase model.ConnectionHelloPhaseType, waitingDuration time.Duration, prolongation bool) error {
	helloMsg := model.ConnectionHello{
		ConnectionHello: model.ConnectionHelloType{
			Phase: phase,
		},
	}

	if waitingDuration > 0 {
		helloMsg.ConnectionHello.Waiting = util.Ptr(uint(waitingDuration.Milliseconds()))
	}
	if prolongation {
		helloMsg.ConnectionHello.ProlongationRequest = &prolongation
	}

	if err := c.sendShipModel(model.MsgTypeControl, helloMsg); err != nil {
		return err
	}
	return nil
}
