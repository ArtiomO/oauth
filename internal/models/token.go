package models

type TokenIn struct {
	Code         string `form:"code"`
	GrantType    string `form:"grant_type"`
	RedirectUri  string `form:"redirect_uri"`
	CodeVerifier string `form:"code_verifier"`
}
