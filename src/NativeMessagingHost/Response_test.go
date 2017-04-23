package NativeMessagingHost

import (
	"testing"
	"encoding/binary"
	"bytes"
	"io"
	"encoding/json"
)

func TestResponse_Send(t *testing.T) {
	response := &Response{Status: "test", Data: "testData"}
	var buffer bytes.Buffer
	// send
	err := response.Send(&buffer)
	if err != nil {
		t.Errorf("got error sending: %s", err)
	}

	// read and validate
	// read length
	var n uint32
	err = binary.Read(&buffer, binary.LittleEndian, &n);
	if err != nil {
		t.Errorf("error reading 4 bytes from buffer: %s", err)
	}
	// read json
	reader := &io.LimitedReader{R: &buffer, N: int64(n)}
	var testResponse Response
	err = json.NewDecoder(reader).Decode(&testResponse)
	if err != nil {
		t.Errorf("could not read json from buffer: %s", err)
	}
	// validate json
	if testResponse.Status != response.Status || testResponse.Data != response.Data {
		t.Errorf("sent and received data is not equal. sent: %s, received: %s", response, testResponse)
	}
}