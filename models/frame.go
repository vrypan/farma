package models

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	db "github.com/vrypan/farma/localdb"
	"google.golang.org/protobuf/proto"
)

var (
	FRAME_NOT_FOUND  = errors.New("Not Found")
	FRAME_NOT_STORED = errors.New("Not Stored")
)

func NewFrame() *Frame {
	return &Frame{}
}

func (f *Frame) Key(id string) string {
	return fmt.Sprintf("f:id:%s", id)
}

func (f *Frame) Save() error {
	db.AssertOpen()
	if f.Id == "" {
		newId := uuid.New()
		encodedId := base64.RawURLEncoding.EncodeToString(newId[:])
		f.Id = encodedId
	}

	key := f.Key(f.Id)
	data, err := proto.Marshal(f)
	if err != nil {
		return fmt.Errorf("Error marshaling Frame protobuf: %v\n", err)
	}

	if err = db.Set([]byte(key), data); err != nil {
		return fmt.Errorf("Error saving Frame: %v\n", err)
	}

	// In addition to f:id:<id>, we also save
	// f:pk:<public_key>
	// This makes it very fast to find a Frame by public_key

	if f.PublicKey != nil {
		if err = f.PublicKey.Save(); err != nil {
			return err
		}
	}
	return nil
}

func (f *Frame) FromId(id string) *Frame {
	key := f.Key(id)
	data, err := db.Get([]byte(key))
	if err != nil {
		return nil
	}
	if proto.Unmarshal(data, f) != nil {
		return nil
	}
	return f
}

func (f *Frame) FromEndpoint(endpoint string) *Frame {
	parts := strings.Split(endpoint, "/")
	if len(parts) != 3 {
		return nil
	}
	frameId := f.Key(parts[2])
	return f.FromId(frameId)
}

func (f *Frame) Delete() error {
	db.AssertOpen()

	if f.Id == "" {
		return fmt.Errorf("Frame.Id is empty")
	}

	key := f.Key(f.Id)
	if err := db.Delete([]byte(key)); err != nil {
		return err
	}

	return nil
}

func (f *Frame) Update() error {
	db.AssertOpen()

	if f.Id == "" {
		return fmt.Errorf("Frame.Id is empty")
	}

	key := f.Key(f.Id)
	data, err := proto.Marshal(f)
	if err != nil {
		return err
	}

	if err = db.Set([]byte(key), data); err != nil {
		return err
	}
	return nil
}

func AllFrames() []*Frame {
	db.AssertOpen()
	res, _, _ := db.GetPrefixP([]byte("f:id:"), []byte("f:id:"), 1000)
	frames := make([]*Frame, len(res))
	for i, frameBytes := range res {
		frame := NewFrame()
		if proto.Unmarshal(frameBytes, frame) != nil {
			return nil
		}
		frames[i] = frame
	}
	return frames
}
