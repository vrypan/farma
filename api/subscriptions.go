package api

import (
	"fmt"

	db "github.com/vrypan/farma/localdb"
	"github.com/vrypan/farma/utils"
	"google.golang.org/protobuf/proto"
)

func ShowSubscriptions(frameId uint64, limit int) string {
	prefix := "s:id:"
	if frameId > 0 {
		prefix = fmt.Sprintf("s:id:%d:", frameId)
	}

	data, _, err := db.GetPrefixP([]byte(prefix), []byte(prefix), limit)

	if err != nil {
		Error("Error fetching data", err)
	}
	list := make([]*utils.Subscription, len(data))
	for i, item := range data {
		var pb utils.Subscription
		proto.Unmarshal(item, &pb)
		list[i] = &pb
	}
	response := Response{
		Status:  "SUCCESS",
		Message: "Subscriptions",
		Data:    list,
	}
	return response.String()
}
