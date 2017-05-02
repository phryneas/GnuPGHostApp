package OpenPgpJsApi

import (
	"testing"
)

const encryptedTest = `-----BEGIN PGP MESSAGE-----

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