package models

import (
	"fmt"

	"github.com/google/uuid"
	db "github.com/vrypan/farma/localdb"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func NewNotification(
	frameId string,
	appId uint64,
	id string,
	title string,
	message string,
	link string,
	endpoint string,
	tokens map[string]uint64,
) *Notification {
	if id == "" {
		id = uuid.New().String()
	}
	return &Notification{
		FrameId:  frameId,
		AppId:    appId,
		Id:       id,
		Endpoint: endpoint,
		Title:    title,
		Message:  message,
		Link:     link,
		Tokens:   tokens,
	}
}

func (n *Notification) Key() string {
	return n.Prefix() + fmt.Sprintf("%03d", *n.Version)
}
func (n *Notification) Json() []byte {
	json, err := protojson.Marshal(n)
	if err != nil {
		return nil
	}
	return json
}
func (n *Notification) Type() string {
	return "Notification"
}

func (n *Notification) Prefix() string {
	return "n:id:" + n.FrameId + ":" + n.Id + ":"
}
func (n *Notification) PrefixBytes() []byte {
	return []byte(n.Prefix())
}

func (n *Notification) Save() (int, error) {
	n.Ctime = timestamppb.Now()
	nextVersion := uint64(0)
	var err error
	prefix := n.PrefixBytes()
	for {
		keys, next, err := db.GetKeysWithPrefix(prefix, prefix, 100)
		if err != nil {
			return 0, fmt.Errorf("Error getting keys: %v", err)
		}
		nextVersion += uint64(len(keys))
		if next == nil {
			break
		}
	}
	nextKey := []byte(n.Prefix() + fmt.Sprintf("%03d", nextVersion))
	n.Version = &nextVersion
	notificationBytes, err := proto.Marshal(n)
	if err != nil {
		return 0, fmt.Errorf("Error marshaling notification: %v", err)
	}
	err = db.Set(nextKey, notificationBytes)
	if err != nil {
		return 0, fmt.Errorf("Error saving notification: %v", err)
	}
	return int(nextVersion), nil
}

func (n *Notification) Update() (int, error) {
	key := []byte(n.Prefix() + fmt.Sprintf("%03d", n.GetVersion()))
	notificationBytes, err := proto.Marshal(n)
	if err != nil {
		return 0, fmt.Errorf("Error marshaling notification: %v", err)
	}
	err = db.Set(key, notificationBytes)
	if err != nil {
		return 0, fmt.Errorf("Error saving notification: %v", err)
	}
	return int(*n.Version), nil
}

func (n *Notification) Load(frameId string, notificationId string) ([]*Notification, error) {
	var notifications []*Notification
	n.FrameId = frameId
	n.Id = notificationId
	prefix := n.PrefixBytes()
	next := prefix
	for {
		keys, next, err := db.GetKeysWithPrefix(prefix, next, 100)
		if err != nil {
			return nil, fmt.Errorf("Error getting keys: %v", err)
		}
		if len(keys) > 0 {
			for _, key := range keys {
				value, err := db.Get(key)
				if err != nil {
					return nil, fmt.Errorf("Error getting key %s: %v", key, err)
				}
				notification := Notification{}
				err = proto.Unmarshal(value, &notification)
				if err != nil {
					return nil, fmt.Errorf("Error unmarshaling value for key %s: %v", key, err)
				}
				notifications = append(notifications, &notification)
			}
		}
		if next == nil {
			break
		}
	}
	return notifications, nil
}
