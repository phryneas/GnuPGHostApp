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
					Encrypt: OpenPgpJsApi.EncryptRequest{DataString: "test", Armor: true, PublicKeys: []string{"1E43F132357B5AD55CECCCC3067D1766157F6495"}},
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
					Decrypt: OpenPgpJsApi.DecryptRequest{
						Message: `-----BEGIN PGP MESSAGE-----

hIwDRsAl1a6NFu4BA/4+anSlCSKchJbf8A+E05VEFkS0DLx823GmlpuEMk/dCv4U
uESYLDl1FZA/H/m7DAEyHVMra4NBxPCmqY3mLCHpn/C8GEuVvCYqXSDDEy72ujPG
o+1ISH3/FU1xonSPGy4D9xo2bUe1f/s3Rwp04/aFVJS5jree9ZqsRiqdjcfqk9LA
PAHHwnPLX4i9S1mR0uuDGTU1z8yup7UNjt2V4cJfudPDLQMB7XyviJOTlsG3K8JD
AJkqRttwzh5y30HK28+sO5OkgMJI6E8XCQ77QRmwNDPf+j7Qm4HYGY91VnYL3uTX
wPy+1Z6jN/5He/ObTgL1wRoe4Ng73IvfUyN9JBbfOYo/+xB2Iyk7vHzeZ1m/PEYG
MCsyPRXj6jA8/ql0hdjYMp//ScF+VsjF2840yX6mMgweW4Wss932sx34FNJ0FQoE
4Th8K+PmjJNyILW/X1an828KjTmPXmpnfE81gbaT0gDKAMERIoCpndnAbKArAsji
ZGQMh/4azTesg4Tgcg==
=/jod
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

				if len(result.Signatures)==0 || result.Signatures[0].Keyid != "1E43F132357B5AD55CECCCC3067D1766157F6495" {
					t.Errorf("wrong keyid, expected 1E43F132357B5AD55CECCCC3067D1766157F6495, got %s", result.Signatures)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.request.Action, func(t *testing.T) {
			writer := new(bytes.Buffer)
			reader := new(bytes.Buffer)

			NativeMessagingHost.SendItem(writer, test.request)
			t.Logf(writer.String())

			err := LoopExecution(writer, reader)

			if err == io.EOF {
				t.Logf("got EOF: %s", err)
			} else if err != nil {
				t.Fatalf("failed with %s", err)
			}

			t.Logf(reader.String())

			response := NativeMessagingHost.Response{}
			decoder, err := NativeMessagingHost.PrepareDecoder(reader)
			decoder.Decode(&response)

			test.validation(response, t)
		})
	}
}
