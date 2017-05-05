package OpenPgpJsApi

type Key struct {
	Revoked         bool `json:"revoked"`
	Expired         bool `json:"expired"`
	Disabled        bool `json:"disabled"`
	Secret          bool `json:"secret"`
	CanEncrypt      bool `json:"canEncrypt"`
	CanSign         bool `json:"canSign"`
	CanCertify      bool `json:"canCertify"`
	CanAuthenticate bool `json:"canAuthenticate"`
	OwnerTrust      Validity `json:"ownerTrust"`
	SubKeys         []SubKey `json:"subKeys"`
	UserIDs         []UserID  `json:"userIDs"`
}

type SubKey struct {
	Key
	Invalid     bool `json:"invalid"`
	KeyID       string `json:"keyId"`
	FingerPrint string `json:"fingerPrint"`
	Created     string `json:"created"`
	Expires     string `json:"expires"`
	CardNumber  string `json:"cardNumber"`
}

type UserID struct {
	Invalid  bool `json:"invalid"`
	Validity Validity `json:"validity"`
	UID      string `json:"UID"`
	Name     string `json:"name"`
	Comment  string `json:"comment"`
	Email    string `json:"email"`
}

type Validity int
