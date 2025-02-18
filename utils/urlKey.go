package utils

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type UrlKey struct {
	FrameId  uint64
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
	return fmt.Sprintf("s:url:%d:%d:%d:%s:%s", k.FrameId, k.UserId, k.Status.Number(), url.QueryEscape(k.Endpoint), url.QueryEscape(k.Token))
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
	frameId, _ := strconv.ParseUint(parts[2], 10, 64)
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
