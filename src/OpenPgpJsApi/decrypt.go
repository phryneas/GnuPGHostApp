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
// PublicKeys []string //Key | Array.<Key> 	(optional) array of public keys or single key, to verify Signatures

type OpenPgpJsDecryptRequest struct {
	Message string  `json:"Message"`
	//Message 	the Message object with the encrypted Data
	// passed as Armored String
	Format string `json:"Format"`
	//String 	(optional) return Data Format either as 'utf8' or 'binary'
	// one of 'utf8' or 'binary'
	Signature string  `json:"Signature"`
	//Signature 	(optional) Detached Signature for verification
	// passed as Armored String
}

type OpenPgpJsDecryptSignature struct {
	Keyid string `json:"keyid"`
	Valid bool `json:"valid"`
}

type OpenPgpJsDecryptResult struct {
	Signatures []OpenPgpJsDecryptSignature `json:"signatures"`
	DataString string `json:"data_string"`
	DataBytes  []uint8 `json:"data_bytes"`
}

func (r OpenPgpJsDecryptRequest) Execute() (result OpenPgpJsDecryptResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			result = OpenPgpJsDecryptResult{}
			err = r.(error)
		}
	}()

	ctx, err := gpgme.New()
	handleErr(err)
	defer ctx.Release()

	signature, err := gpgme.NewDataReader(strings.NewReader(r.Signature))
	handleErr(err)
	defer signature.Close()

	message, err := gpgme.NewDataReader(strings.NewReader(r.Message))
	handleErr(err)
	defer message.Close()

	plain, err := gpgme.NewData()
	handleErr(err)
	defer plain.Close()

	var signatures []gpgme.Signature
	if r.Signature != "" {
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
		result.Signatures = append(result.Signatures, OpenPgpJsDecryptSignature{Keyid: signature.Fingerprint, Valid: validity})
	}

	plain.Seek(0, gpgme.SeekSet)
	buf := new(bytes.Buffer)
	io.Copy(buf, plain)
	result.DataString = buf.String()

	return
}
