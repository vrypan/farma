package utils

import (
	"fmt"
	"net/url"
	"strings"
)

type TokenKey struct {
	Token string
}

func NewTokenKey(token string) TokenKey {
	return TokenKey{
		Token: token,
	}
}

func (k TokenKey) String() string {
	return fmt.Sprintf("s:token:%s", url.QueryEscape(k.Token))
}
func (k TokenKey) Bytes() []byte {
	return []byte(k.String())
}
func (k TokenKey) DecodeBytes(b []byte) TokenKey {
	return k.DecodeString(string(b))
}
func (k TokenKey) DecodeString(s string) TokenKey {
	parts := strings.Split(s, ":")
	token, _ := url.QueryUnescape(parts[2])
	return TokenKey{
		Token: token,
	}
}
