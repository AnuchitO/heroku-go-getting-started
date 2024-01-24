package middlewares

import (
	"encoding/json"

	"google.golang.org/api/idtoken"
)

// Extract from some fields of Google ID Token
// Warning: 'profile' scope may not included in every subsequent ID token received, omitempty to avoid overwriting blank strings to DB

type ParsedIdToken struct {
	// Basic Info
	Issuer   string `json:"iss"`
	Audience string `json:"aud"`
	Expires  int    `json:"exp"`
	IssuedAt int    `json:"iat"`
	Subject  string `json:"sub"`
	// Additional Info
	Email           string `json:"email"`
	EmailVerified   bool   `json:"email_verified"`
	FullName        string `json:"name"`
	GivenName       string `json:"given_name"`
	FamilyName      string `json:"family_name"`
	PictureURL      string `json:"picture"`
	Locale          string `json:"locale"`
	AuthorizedParty string `json:"azp"`
	AtHash          string `json:"at_hash"`
	Hd              string `json:"hd"`
}

func NewParsedIdToken(payload *idtoken.Payload) *ParsedIdToken {
	parsedIdToken := new(ParsedIdToken)
	bs, err := json.Marshal(payload.Claims)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(bs, parsedIdToken)
	if err != nil {
		return nil
	}
	return parsedIdToken
}
