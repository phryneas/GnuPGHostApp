package main

import (
	"testing"
	"NativeMessagingHost"
	"OpenPgpJsApi"
	"bytes"
	"io"
	"github.com/go-errors/errors"
)

func TestLoopExecution(t *testing.T) {
	firstRequest := NativeMessagingHost.Request{Action: "encrypt" }
	secondRequest := OpenPgpJsApi.OpenPgpJsEncryptRequest{DataString: "test", Armor: true, PublicKeys: []string{"F9C2408278723D64985CA4A63F3B8061E714CD2C"}}
	writer := new(bytes.Buffer)
	reader := new(bytes.Buffer)
	NativeMessagingHost.SendResponse(firstRequest, writer)
	NativeMessagingHost.SendResponse(secondRequest, writer)
	NativeMessagingHost.SendResponse(secondRequest, writer)

	writer = bytes.NewBufferString(writer.String())

	err := LoopExecution(writer, reader)
	if err == io.EOF {
		t.Logf("got EOF: %s", err)
	} else if err != nil {
		t.Fatalf("failed with %s", err.(*errors.Error).ErrorStack())
	}

	var result OpenPgpJsApi.OpenPgpJsEncryptResult
	decoder, err := NativeMessagingHost.PrepareDecoder(reader)
	decoder.Decode(result)

	t.Log(result);
	t.Fail()

}
