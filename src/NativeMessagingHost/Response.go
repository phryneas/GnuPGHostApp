package NativeMessagingHost

import (
	"io"
	"encoding/binary"
	"encoding/json"
	"bytes"
	"OpenPgpJsApi"
)

type ResponseData struct {
	Encrypt OpenPgpJsApi.OpenPgpJsEncryptResult `json:"encrypt"`
	Decrypt OpenPgpJsApi.OpenPgpJsDecryptResult `json:"decrypt"`
}

type Response struct {
	Status string `json:"status"`
	Message string `json:"message"`
	Data ResponseData `json:"data"`
}

func (r Response) String() string {
	buffer := new(bytes.Buffer)
	buffer.WriteString("Response: ")
	json.NewEncoder(buffer).Encode(r)
	return buffer.String()
}

func (r Response) Send(stdOut io.Writer) (err error) {
	return SendItem(stdOut, r)
}

func SendItem(stdOut io.Writer, item interface{} ) (err error){
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(item);
	binary.Write(stdOut, binary.LittleEndian, uint32(b.Len()))
	b.WriteTo(stdOut)
	return;
}