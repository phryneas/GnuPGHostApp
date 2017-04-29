package OpenPgpJsApi

import (
	"testing"
)

const encryptedTest = `-----BEGIN PGP MESSAGE-----
Version: GnuPG v2

hQEMA0Hq9rAfhmAEAQgAokUPFgHuN5ODm4jfzujXBfVVmqhzATht8Boo5tmIcQxa
zIxznHmwIRDsURe9GzzxWh/NftXZ8Xc6leKd8K6dLD57I6lj5NGDbgRD56+mHg5w
cTkREcajytlbT7eTVGZdY5MZ66bfs67qosGUQA7ltM6qMVRN5ucnEHzGZGgiZ31J
JfnClKquWWFEXU2+Xe8lX5HWVe7BXk/K4OjAndkAzoa7aTHdUF8FQU5l3IZk5Vg1
+McK/PsBm5sH5DXaa3bBBbnqAammBZGxXG3P6bKGCiMjx38v4Ks5Nbq/v1VM4uz7
eyyQJ5ri+V69teO64IwcBQ564ohCku1oLGueh+pZe9LAKwGTKQbhyDIb2qLTA2n8
+xE3kMuNzzy7WKz+4RRTdqdN8uOJS7AAj6bfA+vnmqEFPOf5v8ATI5vv0YLR9d+1
4YPgJ+4rOoHXpBlErAxWH2iZRYwYvaU2tIZfLu6d2otJHdIUN4JgyFl59r+92iJM
8VKFoymetgfmR75WRaZZNVkp/39ivZ2bMgJphhsWAP033/DSzQSzBDYTr1viVsWg
j2G2g6zJ8g/FAYeCikj6VsOimGIfG3MabIIOTfyv0atCMWfcqbZvWHzF+q6ZaGYe
yFUFvIGXdTUh+8TfsCLbVu/2UgjD5M0stBi4oAc=
=ig18
-----END PGP MESSAGE-----`

func TestOpenPgpJsDecryptRequest_Execute(t *testing.T) {
	request := OpenPgpJsDecryptRequest{Message: encryptedTest}
	result, err := request.Execute()
	if err != nil {
		t.Errorf("encountered error: %s", err)
	}
	if result.DataString !="test" {
		t.Errorf("decrypted text should be 'test', was '%s'", result.DataString)
	}
	if result.Signatures[0].Keyid != "1E43F132357B5AD55CECCCC3067D1766157F6495" {
		t.Errorf("wrong keyid, expected 1E43F132357B5AD55CECCCC3067D1766157F6495, got %s", result.Signatures[0].Keyid)
	}
	if !result.Signatures[0].Valid {
		t.Errorf("invalid signature. should not happen")
	}
}