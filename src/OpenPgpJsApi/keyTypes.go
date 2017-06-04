package OpenPgpJsApi

import (
	"github.com/proglottis/gpgme"
	"time"
)

type Key struct {
	Revoked         bool `json:"revoked"`
	Expired         bool `json:"expired"`
	Disabled        bool `json:"disabled"`
	Secret          bool `json:"secret"`
	CanEncrypt      bool `json:"canEncrypt"`
	CanSign         bool `json:"canSign"`
	CanCertify      bool `json:"canCertify"`
	CanAuthenticate bool `json:"canAuthenticate"`
	OwnerTrust      gpgme.Validity `json:"ownerTrust"`
	SubKeys         []SubKey `json:"subKeys"`
	UserIDs         []UserID  `json:"userIDs"`
}

func (key Key) getFingerPrint() string {
	if len(key.SubKeys) == 0 {
		return ""
	}
	return key.SubKeys[0].FingerPrint
}

type SubKey struct {
	Revoked     bool `json:"revoked"`
	Expired     bool `json:"expired"`
	Disabled    bool `json:"disabled"`
	Invalid     bool `json:"invalid"`
	Secret      bool `json:"secret"`
	KeyID       string `json:"keyID"`
	FingerPrint string `json:"fingerPrint"`
	Created     time.Time `json:"created"`
	Expires     time.Time `json:"expires"`
	CardNumber  string `json:"cardNumber"`
}

type UserID struct {
	Revoked  bool `json:"revoked"`
	Invalid  bool `json:"invalid"`
	Validity gpgme.Validity `json:"validity"`
	UID      string `json:"UID"`
	Name     string `json:"name"`
	Comment  string `json:"comment"`
	Email    string `json:"email"`
}

func newKey(key *gpgme.Key) (ret Key) {
	ret = Key{
		Revoked:         key.Revoked(),
		Expired:         key.Expired(),
		Disabled:        key.Disabled(),
		Secret:          key.Secret(),
		CanEncrypt:      key.CanEncrypt(),
		CanSign:         key.CanSign(),
		CanCertify:      key.CanCertify(),
		CanAuthenticate: key.CanAuthenticate(),
		OwnerTrust:      key.OwnerTrust(),
	}

	for subKey := key.SubKeys(); subKey != nil; subKey = subKey.Next() {
		ret.SubKeys = append(ret.SubKeys, SubKey{
			Revoked:     subKey.Revoked(),
			Expired:     subKey.Expired(),
			Disabled:    subKey.Disabled(),
			Invalid:     subKey.Invalid(),
			Secret:      subKey.Secret(),
			KeyID:       subKey.KeyID(),
			FingerPrint: subKey.Fingerprint(),
			Created:     subKey.Created(),
			Expires:     subKey.Expires(),
			CardNumber:  subKey.CardNumber(),
		})
	}

	for userID := key.UserIDs(); userID != nil; userID = userID.Next() {
		ret.UserIDs = append(ret.UserIDs, UserID{
			Invalid:  userID.Invalid(),
			Revoked:  userID.Revoked(),
			Validity: userID.Validity(),
			UID:      userID.UID(),
			Name:     userID.Name(),
			Comment:  userID.Comment(),
			Email:    userID.Email(),
		})
	}

	return
}
