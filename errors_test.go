package noise

import (
	"errors"
	"fmt"
	"testing"
)

const MOCK_ADDRESS = "127.0.0.1:2379"
const STATEMENT = "Expected error: %v, got: %v"

func TestDialingError(t *testing.T) {
	customError := "Fail dial"
	err := errors.New(customError)
	output := errDialingNode(err)
	expected := fmt.Sprintf("net: error during dialing -> %v", err)

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}

}

func TestBindingError(t *testing.T) {
	customError := "Fail binding"
	err := errors.New(customError)
	output := errBindingConnection(err)
	expected := fmt.Sprintf("net: connection closed or cannot be established -> %v", err)

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}
}

func TestExceededMaxPeersError(t *testing.T) {

	output := errExceededMaxPeers(10)
	expected := fmt.Sprintf("overflow: it is not possible to accept more than %d connections -> max peers exceeded", 10)

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}
}

func TestExceededMaxPayloadSize(t *testing.T) {
	customError := "Fail setting up"
	err := errors.New(customError)
	output := errSettingUpConnection(err)
	expected := fmt.Sprintf("ops: error trying to configure connection -> %v", err)

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}
}

func TestSendingMessageToInvalidPeer(t *testing.T) {
	err := fmt.Errorf("peer disconnected: %s", MOCK_ADDRESS)
	output := errSendingMessage(err)
	expected := fmt.Sprintf("ops: error sending message -> peer disconnected: %s", MOCK_ADDRESS)

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}
}

func TestErrVerifyingSignature(t *testing.T) {
	err := errors.New("invalid signature: ABC")
	output := errVerifyingSignature(err)
	expected := "sec: error verifying signature -> invalid signature: ABC"

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}
}

func TestErrDuringHandshake(t *testing.T) {
	err := errors.New("fail")
	output := errDuringHandshake(err)
	expected := "ops: error during handshake -> fail"

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}
}
