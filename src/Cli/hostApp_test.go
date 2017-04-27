package main

import (
	"testing"
	"NativeMessagingHost"
	"OpenPgpJsApi"
	"bytes"
	"io"
	"strings"
)

type testCase struct {
	request    NativeMessagingHost.Request
	validation func(result NativeMessagingHost.Response, t *testing.T)
}

func TestLoopExecution(t *testing.T) {
	tests := []testCase{
		{
			NativeMessagingHost.Request{
				Action: "encrypt",
				Data: NativeMessagingHost.RequestData{
					Encrypt: OpenPgpJsApi.OpenPgpJsEncryptRequest{DataString: "test", Armor: true, PublicKeys: []string{"F9C2408278723D64985CA4A63F3B8061E714CD2C"}},
				},
			},
			func(response NativeMessagingHost.Response, t *testing.T) {
				result := response.Data.Encrypt
				if !strings.Contains(result.Data, "BEGIN PGP MESSAGE") {
					t.Errorf("armored encrypted data should contain PGP header, was '%s'", result.Data)
				}
			},
		},
		{
			NativeMessagingHost.Request{
				Action: "decrypt",
				Data: NativeMessagingHost.RequestData{
					Decrypt: OpenPgpJsApi.OpenPgpJsDecryptRequest{
						Message: `-----BEGIN PGP MESSAGE-----
Version: GnuPG v2

hQEMA0lF8LJU85UOAQgA5BjlfmBncN+9G8SFFhHvpZLUzj+z9f6AEtUCE8oI8oSb
TkQKX00cbYNYbUvJ14bOMZA/Z6rBRqa5ko3RbZ11DYHvOqwcgUv4k+lic1nU7Uwt
X+gFO6RgiJVze1bmtWAg/6OsAIOV5sf5pYEkAV2u+VzaW8adCHyvIzdINa9j8SMm
JZh/MAWRC5A+GzMnyR7qpmtVswAOR/1aiqm0nE+WnvIPXs0JqNCoLT0jPFsBHKv5
CEDBHlgxJOFUNvXEG7klii/Yuoe8uW6RDLArksh4OTkELOT4GwRQC+ZelycghdY8
2tiEmVv4EqrGv4/0AnLz8XM2iOA7AZmo8ZdSm58OY9LAsAFD++l9M3RkWBd6Na3C
CGwes62uYCzzBINX6ftBoblZc1QbBeCgbatjqiEeMxnrp52lSr9e79lBUr4y4QAU
E+pV6oKprDdne68aKjFmfbToBpneUz5xR4NK9+NnFZ2MtGEUp6Psk74x9AyvUm59
IPIyZg6uZuGl8nUfRoGql9qOdXR3WL9xgQ4+NPkRsMB07y+0RC2INRft+vH2CVOi
aUQLMWs3FkHAwxZ2DmHvwZwCgh6fStPLmCVHL2pttqrFnPjigYzQgZF1ski5kAP2
+itB8BiuprIjC1BlpuyMg16O86VfQX2NfgL5KMnpyxnWFo/SIxvOrEM6MK9xblln
riRfMsR0aeygr35U4nIYU4innYeNCLg5l53FdVlWdTDhyd2ixJ1gQKBumoDzFcph
M+voe5nNDHzqGui80YASmll9wmTF6+MVEB5078VLfqDnjWCyt+knMFmoO9EVwqCD
u8sctMeF0yFMwDfshzokmSE/
=zqnU
-----END PGP MESSAGE-----`,
						Format: "utf8",
					},
				},
			},
			func(response NativeMessagingHost.Response, t *testing.T) {
				result := response.Data.Decrypt
				if result.DataString !="test" {
					t.Errorf("decrypted text should be 'test', was '%s'", result.DataString)
				}

				if len(result.Signatures)==0 || result.Signatures[0].Keyid != "F9C2408278723D64985CA4A63F3B8061E714CD2C" {
					t.Errorf("wrong keyid, expected F9C2408278723D64985CA4A63F3B8061E714CD2C, got %s", result.Signatures)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.request.Action, func(t *testing.T) {
			writer := new(bytes.Buffer)
			reader := new(bytes.Buffer)

			NativeMessagingHost.SendItem(writer, test.request)
//			fmt.Println(writer.String())

			err := LoopExecution(writer, reader)

			if err == io.EOF {
				t.Logf("got EOF: %s", err)
			} else if err != nil {
				t.Fatalf("failed with %s", err)
			}

			response := NativeMessagingHost.Response{}
			decoder, err := NativeMessagingHost.PrepareDecoder(reader)
			decoder.Decode(&response)

			test.validation(response, t)
		})
	}
}
