package OpenPgpJsApi

import (
	"testing"
	"fmt"
)

func TestOpenPgpJsEncryptRequest_Execute(t *testing.T) {
	request := OpenPgpJsEncryptRequest{DataString:"test", PublicKeys:[]string{"B76C94DD"}, Armor:true}
	result, err := request.Execute()
	if err != nil {
		t.Errorf("encrypt request errored: %s", err)
	}

	decrypt := OpenPgpJsDecryptRequest{Message:result.Data, Format:"utf8"}
	decrypted, err := decrypt.Execute()
	if err != nil {
		t.Errorf("decrypting encrypted data errored: %s \n request was %s", err, decrypt)
	}
	fmt.Println(decrypted)
}