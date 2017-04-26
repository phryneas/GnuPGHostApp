package NativeMessagingHost

import (
	"io"
	"encoding/binary"
	"encoding/json"
	"bytes"
	"fmt"
)

type Response struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

func (r Response) String() string {
	return fmt.Sprintf("Response{Status: %s, Data: %s}", r.Status, r.Data)
}

func SendResponse(r interface{}, stdOut io.Writer) (err error) {
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(r);
	binary.Write(stdOut, binary.LittleEndian, uint32(b.Len()));
	b.WriteTo(stdOut)
	return;
}