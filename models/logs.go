package models

import (
	"fmt"

	db "github.com/vrypan/farma/localdb"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (l *UserLog) Key() string {
	return fmt.Sprintf("l:user:%s:%d:%d", l.FrameId, l.UserId, l.Ctime.Seconds)
}
func (l *UserLog) Json() []byte {
	json, err := protojson.Marshal(l)
	if err != nil {
		return nil
	}
	return json
}
func (l *UserLog) Type() string {
	return "UserLog"
}

func (l *UserLog) Save() error {
	l.Ctime = timestamppb.Now()
	key := l.Key()
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
	prefix := fmt.Sprintf("l:user:%s:%d:", l.FrameId, l.UserId)
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
