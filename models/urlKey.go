package models

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	db "github.com/vrypan/farma/localdb"
)

type UrlKey struct {
	FrameId  string
	UserId   uint64
	Status   SubscriptionStatus
	Endpoint string
	Token    string
}

func (k UrlKey) FromSubscription(sub *Subscription) UrlKey {
	return UrlKey{
		FrameId:  sub.FrameId,
		UserId:   sub.UserId,
		Status:   sub.Status,
		Endpoint: sub.Url,
		Token:    sub.Token,
	}
}

func (k UrlKey) String() string {
	return fmt.Sprintf("s:url:%s:%d:%d:%s:%s", k.FrameId, k.UserId, k.Status.Number(), url.QueryEscape(k.Endpoint), url.QueryEscape(k.Token))
}
func (k UrlKey) Bytes() []byte {
	return []byte(k.String())
}
func (k UrlKey) DecodeBytes(b []byte) UrlKey {
	return k.DecodeString(string(b))
}
func (k UrlKey) DecodeString(s string) UrlKey {
	parts := strings.Split(s, ":")
	//if len(parts) == 7 {
	frameId := parts[2]
	userId, _ := strconv.ParseUint(parts[3], 10, 64)
	statusNum, _ := strconv.Atoi(parts[4])
	status := SubscriptionStatus(statusNum)
	endpoint, _ := url.QueryUnescape(parts[5])
	token, _ := url.QueryUnescape(parts[6])
	return UrlKey{
		FrameId:  frameId,
		UserId:   userId,
		Status:   status,
		Endpoint: endpoint,
		Token:    token,
	}
	//}
	//return UrlKey{}
}

func (k UrlKey) Set(subscriptionKey []byte) error {
	prefix := fmt.Appendf([]byte(""), "s:url:%s:%d:", k.FrameId, k.UserId)
	existingKeys, _, err := db.GetKeysWithPrefix(prefix, prefix, 10)
	if err != nil {
		return err
	}
	for _, existingKey := range existingKeys {
		db.Delete(existingKey)
	}
	err = db.Set(k.Bytes(), subscriptionKey)
	return err
}

func (k UrlKey) Get() ([]byte, error) {
	return db.Get(k.Bytes())
}

func (k UrlKey) Delete() error {
	return db.Delete(k.Bytes())
}
