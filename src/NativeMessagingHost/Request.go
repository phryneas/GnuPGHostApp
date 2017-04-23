package NativeMessagingHost

import (
	"io"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

type Request struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}

func (r Request) String() string {
	return fmt.Sprintf("%s [%s]", r.Action, r.Data);
}

func ReadRequest(stdIn io.Reader) (request Request, err error) {
	var n uint32
	err = binary.Read(stdIn, binary.LittleEndian, &n);
	if err != nil {
		return
	}
	reader := &io.LimitedReader{R: stdIn, N: int64(n)}
	err = json.NewDecoder(reader).Decode(&request)
	return;
}