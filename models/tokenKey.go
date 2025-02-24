package models

import (
	"fmt"
	"net/url"
	"strings"

	db "github.com/vrypan/farma/localdb"
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

func (k TokenKey) Set(subscriptionKey []byte) error {
	err := db.Set(k.Bytes(), subscriptionKey)
	return err
}

func (k TokenKey) Get() ([]byte, error) {
	return db.Get(k.Bytes())
}
func (k TokenKey) Delete() error {
	return db.Delete(k.Bytes())
}
