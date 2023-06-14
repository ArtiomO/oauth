package models

import (
	"encoding/json"
)

type LoginFormIn struct {
	ClientId            string `form:"client_id"`
	RedirectUri         string `form:"redirect_uri"`
	State               string `form:"state"`
	CodeChallenge       string `form:"code_challenge"`
	CodeChallengeMethod string `form:"code_challenge_method"`
}

func LoginInFromString(s string) LoginFormIn {
	var loginFormReg LoginFormIn
	err := json.Unmarshal([]byte(s), &loginFormReg)
	if err != nil {
		panic(err)
	}
	return loginFormReg
}

func (c* LoginFormIn) String() string {
	out, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(out)
}

type LoginIn struct {
	Email    string `form:"email"`
	Password string `form:"password"`
	ReqId    string `form:"reqid"`
}

