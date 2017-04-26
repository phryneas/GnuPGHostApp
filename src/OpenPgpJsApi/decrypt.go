package OpenPgpJsApi

import (
	"github.com/proglottis/gpgme"
	"io"
	"strings"
	"bytes"
)

// @see https://openpgpjs.org/openpgpjs/doc/module-openpgp.html#.decrypt

// unused options from the openpgp.js API:
// privateKey  //Key 	(optional) private key with decrypted secret key Data or session key
// sessionKey //Object 	(optional) session key in the form: { Data:Uint8Array, algorithm:String }
// password   string //String 	(optional) single password to decrypt the Message
// PublicKeys []string //Key | Array.<Key> 	(optional) array of public keys or single key, to verify signatures

type OpenPgpJsDecryptRequest struct {
	message string  `json:"Message"`
	//Message 	the Message object with the encrypted Data
	// passed as Armored String
	format string `json:"format"`
	//String 	(optional) return Data format either as 'utf8' or 'binary'
	// one of 'utf8' or 'binary'
	signature string  `json:"Signature"`
	//Signature 	(optional) Detached Signature for verification
	// passed as Armored String
}

type OpenPgpJsDecryptSignature struct {
	keyid string
	valid bool
}

type OpenPgpJsDecryptResult struct {
	signatures []OpenPgpJsDecryptSignature
	dataString string
	dataBytes  []uint8
}

func (r OpenPgpJsDecryptRequest) Execute() (result RequestResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			result = nil
			err = r.(error)
		}
	}()

	res := OpenPgpJsDecryptResult{}

	ctx, err := gpgme.New()
	signature, err := gpgme.NewDataReader(strings.NewReader(r.signature))
	handleErr(err)
	defer signature.Close()

	message, err := gpgme.NewDataReader(strings.NewReader(r.message))
	handleErr(err)
	defer message.Close()

	plain, err := gpgme.NewData()
	handleErr(err)
	defer plain.Close()

	var signatures []gpgme.Signature
	if (r.signature != ""){
		_, signatures, err = ctx.Verify(signature, message, nil)
		handleErr(err)

		err = ctx.Decrypt(message, plain)
		handleErr(err)
	} else {
		_, signatures, err = ctx.Verify(message, nil, plain)
		handleErr(err)
	}

	for _, signature := range signatures {
		validity := (signature.Summary & gpgme.SigSumGreen) != 0
		res.signatures = append(res.signatures, OpenPgpJsDecryptSignature{keyid: signature.Fingerprint, valid: validity})
	}


	plain.Seek(0, gpgme.SeekSet)

	buf := new(bytes.Buffer)
	io.Copy(buf, plain)
	res.dataString = buf.String()
	return res, nil
}
