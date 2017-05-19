package OpenPgpJsApi

import (
	"testing"
)

const text = "test";

const keyFingerPrint = "1E43F132357B5AD55CECCCC3067D1766157F6495";

const encryptedAndSignedTest = `-----BEGIN PGP MESSAGE-----

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

const encrypted =`-----BEGIN PGP MESSAGE-----

hIwDRsAl1a6NFu4BA/9TY7j+3VNzUA2wl1Used0b0btx7/7LjQulhIb939pa+ee5
4/f2MGuwR6D37S+v972DjQHIMCE921KEcA84xV9LFD6qimZn9Y9m7B/v0V52Uqrg
eOJUmcBtyM6lMnqomZZpigni3btV3f8h8KGVzT70Kueq/w2Gxr+o3pvC0aKbVdI/
ARzE6XtHzjRgoeC5x+vMWymTl9WAoeSiGAJAq86kGCnbuKtJx/kbea2kgRXwOSAl
qk+rR5zdljYqR3XrnAvM
=vhs6
-----END PGP MESSAGE-----
`;

const detachedSignature = `-----BEGIN PGP SIGNATURE-----

iJwEAAEIAAYFAlkevLAACgkQBn0XZhV/ZJXrVQQAmGB/4zv2Ltq4M4dfcXPmMuqY
sS0SCqF+j1agLfy2Zh4XO66RhLJZ3o0obVyv7lJKdZ7RDzf69iJa8nZdRKamVMnF
2VifL9sOnFC0z4CipcsUkuIw7A+Azom32XQrPKf58CEzQRvnH5crkowNsSCvJlKd
wQ1iCtA2HokiN4rLgzA=
=58A4
-----END PGP SIGNATURE-----`;

func TestOpenPgpJsDecryptRequest_Execute_InlineSignature(t *testing.T) {
	request := DecryptRequest{Message: encryptedAndSignedTest}
	result, err := request.Execute()
	t.Logf("\nresult was %#v\nerror was %+v", result, err)
	if err != nil {
		t.Errorf("encountered error: %s", err)
	}
	if result.DataString != text {
		t.Errorf("decrypted text should be 'test', was '%s'", result.DataString)
	}
	if result.Signatures[0].Keyid != keyFingerPrint {
		t.Errorf("wrong keyid, expected %s, got %s", keyFingerPrint, result.Signatures[0].Keyid)
	}
	if !result.Signatures[0].Valid {
		t.Error("invalid signature. should not happen")
	}
}

func TestOpenPgpJsDecryptRequest_Execute_DetachedSignature(t *testing.T) {
	request := DecryptRequest{Message: encrypted, Signature: detachedSignature}
	result, err := request.Execute()
	t.Logf("\nresult was %#v\nerror was %+v", result, err)
	if err != nil {
		t.Errorf("encountered error: %s", err)
	}
	if result.DataString != text {
		t.Errorf("decrypted text should be 'test', was '%s'", result.DataString)
	}
	if result.Signatures[0].Keyid != keyFingerPrint {
		t.Errorf("wrong keyid, expected %s, got %s", keyFingerPrint, result.Signatures[0].Keyid)
	}
	if !result.Signatures[0].Valid {
		t.Error("invalid signature. should not happen")
	}
}