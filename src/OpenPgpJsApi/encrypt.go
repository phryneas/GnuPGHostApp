package OpenPgpJsApi

import (
	"github.com/proglottis/gpgme"
	"io"
	"strings"
	"bytes"
)

// @see https://openpgpjs.org/openpgpjs/doc/module-openpgp.html#.encrypt

// unused options from the openpgp.js API:
//passwords   //String | Array.<String> 	(optional) array of passwords or a single password to encrypt the Message
//filename    //String 	(optional) a filename for the literal Data packet
type OpenPgpJsEncryptRequest struct {
	DataString  string `json:"DataString"`      //String | Uint8Array 	text/Data to be encrypted as JavaScript binary string or Uint8Array
	DataBytes   []byte  `json:"DataBytes"`      //String | Uint8Array 	text/Data to be encrypted as JavaScript binary string or Uint8Array
	PublicKeys  []string `json:"public_keys"`   //Key | Array.<Key> 	(optional) array of keys or single key, used to encrypt the Message
	PrivateKeys []string  `json:"private_keys"` //Key | Array.<Key> 	(optional) private keys for signing. If omitted Message will not be signed
	Armor       bool `json:"Armor"`             //Boolean 	(optional) if the return values should be ascii armored or the Message/Signature objects
	Detached    bool `json:"Detached"`          //Boolean 	(optional) if the Signature should be Detached (if true, Signature will be added to returned object)
	Signature   interface{} `json:"Signature"`  //Signature 	(optional) a Detached Signature to add to the encrypted Message
}

type OpenPgpJsEncryptResult struct {
	Data      string // ASCII armored Message if 'Armor' is true
	Message   []byte // full Message object if 'Armor' is false
	Signature []byte //Detached Signature if 'Detached' is true
}

func (r OpenPgpJsEncryptRequest) Execute() (result RequestResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			result = nil
			err = r.(error)
		}
	}()

	res := OpenPgpJsEncryptResult{}

	recipients := make([]*gpgme.Key, 0)
	for _, pubKey := range r.PublicKeys {
		foundKeys, err := gpgme.FindKeys(pubKey, false)
		handleErr(err)
		for _, foundKey := range foundKeys {
			recipients = append(recipients, foundKey)
		}
	}

	var plain *gpgme.Data
	cipher, err := gpgme.NewData()
	if r.DataString != "" {
		plain, err = gpgme.NewDataReader(strings.NewReader(r.DataString))
	} else {
		plain, err = gpgme.NewDataReader(bytes.NewReader(r.DataBytes))
	}
	handleErr(err)

	ctx, err := gpgme.New()
	handleErr(err)
	ctx.SetArmor(r.Armor)
	err = ctx.Encrypt(recipients, 0, plain, cipher)
	handleErr(err)

	buf := new(bytes.Buffer)
	io.Copy(buf, cipher)
	if r.Armor {
		res.Data = buf.String()
	} else {
		res.Message = buf.Bytes()
	}

	return res, nil
}
