package NativeMessagingHost

import (
	"testing"
	"strings"
	"encoding/binary"
)

func TestReadRequest(t *testing.T) {
	// read illegal value
	testReader := strings.NewReader("illegal value")
	request, err := ReadRequest(testReader)
	if err == nil {
		t.Errorf("illegal value should cause an error, got %s %s", request, err)
	}

	// read legal value
	testString := "{\"action\":\"test\",\"data\":\"testData\"}";
	stringLengthBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(stringLengthBytes, uint32(len(testString)))
	testReader = strings.NewReader(string(stringLengthBytes[:]) + testString)
	request, err = ReadRequest(testReader)
	if err != nil {
		t.Errorf("legal string should not cause an error, got %s %s", request, err)
	}
	// validate data
	if request.Action != "test" || request.Data != "testData" {
		t.Errorf("illegal values read from json")
	}
}
