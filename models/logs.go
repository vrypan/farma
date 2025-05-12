package models

import (
	"fmt"

	db "github.com/vrypan/farma/localdb"
	"github.com/vrypan/farma/natsclient"
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
func (l *UserLog) Topic() string {
	eventTypes := map[EventType]string{
		EventType_NONE:                           "none",
		EventType_FRAME_ADDED:                    "user.subscription.on.miniapp",
		EventType_FRAME_REMOVED:                  "user.subscription.off.miniapp",
		EventType_NOTIFICATIONS_ENABLED:          "user.subscription.on",
		EventType_NOTIFICATIONS_DISABLED:         "user.subscription.off",
		EventType_NOTIFICATION_SENT:              "user.notification.sent",
		EventType_NOTIFICATION_FAILED_OTHER:      "user.notification.failed.other",
		EventType_NOTIFICATION_FAILED_INVALID:    "user.notification.failed.invalid",
		EventType_NOTIFICATION_FAILED_RATE_LIMIT: "user.notification.failed.rate_limit",
	}
	topic, ok := eventTypes[l.EvtType]
	if !ok {
		topic = "unknown"
	}
	return topic
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
	return natsclient.Publish(l)
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

func (l *NotificationLog) Key() string {
	return ""
}
func (l *NotificationLog) Json() []byte {
	json, err := protojson.Marshal(l)
	if err != nil {
		return nil
	}
	return json
}
func (l *NotificationLog) Topic() string {
	return "miniapp.notification.batch"
}
