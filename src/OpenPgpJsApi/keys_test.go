package OpenPgpJsApi

import (
	"testing"
)

func TestFindKeyRequest_Execute(t *testing.T) {
	// TODO

	r := FindKeyRequest{Email: "gnupghostapp_tests@example.com"}
	res, _ := r.Execute()

	t.Log(res)
}