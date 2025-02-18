package api

import (
	"fmt"

	db "github.com/vrypan/farma/localdb"
	"github.com/vrypan/farma/utils"
	"google.golang.org/protobuf/proto"
)

func ShowLogs(userId uint64, limit int) string {
	prefix := "l:user:"
	if userId > 0 {
		prefix = fmt.Sprintf("l:user:%d:", userId)
	}

	data, _, err := db.GetPrefixP([]byte(prefix), []byte(prefix), limit)

	if err != nil {
		Error("Error fetching data", err)
	}
	logs := make([]*utils.UserLog, len(data))
	for i, log := range data {
		var logProto utils.UserLog
		proto.Unmarshal(log, &logProto)
		logs[i] = &logProto
	}
	response := Response{
		Status:  "SUCCESS",
		Message: "User logs",
		Data:    logs,
	}
	return response.String()
}
