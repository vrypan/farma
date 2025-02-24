package models

import (
	"fmt"

	db "github.com/vrypan/farma/localdb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (l *UserLog) Save() error {
	now := timestamppb.Now()
	key := fmt.Sprintf("l:user:%d:%d:%d", l.UserId, l.FrameId, now.Seconds)
	l.Ctime = now
	data, err := proto.Marshal(l)
	if err != nil {
		return err
	}
	if err = db.Set([]byte(key), data); err != nil {
		return err
	}
	return nil
}

func (l *UserLog) Load(limit int) ([]*UserLog, error) {
	prefix := fmt.Sprintf("l:user:%d:", l.UserId)
	data, _, err := db.GetPrefixP([]byte(prefix), []byte(prefix), limit)

	if err != nil {
		return nil, err
	}
	logs := make([]*UserLog, len(data))
	for i, log := range data {
		proto.Unmarshal(log, logs[i])
	}
	return logs, nil
}
