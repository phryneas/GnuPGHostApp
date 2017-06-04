package NativeMessagingHost

import (
	"io"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"OpenPgpJsApi"
	"bytes"
)

type RequestData struct {
	Encrypt OpenPgpJsApi.EncryptRequest `json:"encrypt"`
	Decrypt OpenPgpJsApi.DecryptRequest `json:"decrypt"`
	FindKeys OpenPgpJsApi.FindKeyRequest `json:"findKeys"`
	ExportPublicKeys OpenPgpJsApi.ExportPublicKeysRequest `json:"exportPublicKeys"`
}

type Request struct {
	Action string `json:"action"`
	Data   RequestData `json:"data"`
}

func (r Request) String() string {
	buffer := new(bytes.Buffer)
	buffer.WriteString("Request: ")
	json.NewEncoder(buffer).Encode(r)
	return buffer.String()
}


func (r RequestData) String() string {
	return fmt.Sprintf("Encrypt: [%s] , Decrypt: [%s]", r.Encrypt, r.Decrypt);
}

func PrepareDecoder(stdIn io.Reader) (decoder *json.Decoder, err error) {
	var n uint32
	err = binary.Read(stdIn, binary.LittleEndian, &n);
	if err != nil {
		return
	}
	reader := &io.LimitedReader{R: stdIn, N: int64(n)}
	decoder = json.NewDecoder(reader)
	return;
}

func ReadRequest(stdin io.Reader) (request Request, err error) {
	decoder, err := PrepareDecoder(stdin)
	if err != nil {
		return
	}
	err = decoder.Decode(&request)
	if err == io.EOF {
		err = nil
	}
	return
}