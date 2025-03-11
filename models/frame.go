package models

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"

	db "github.com/vrypan/farma/localdb"
	"google.golang.org/protobuf/proto"
)

var (
	FRAME_NOT_FOUND  = errors.New("Not Found")
	FRAME_NOT_STORED = errors.New("Not Stored")
)

func NewFrame() *Frame {
	f := Frame{}
	f.PublicKey = &PubKey{}
	return &f
}

func (f *Frame) Key(id string) string {
	return fmt.Sprintf("f:id:%s", id)
}

// Returns a 5-character unique id
func (f *Frame) newFrameId() string {
	var dbkey string
	var str string
	tries := 0
	for {
		// Generate random bytes.
		r := make([]byte, 3)
		_, err := rand.Read(r)
		if err != nil {
			panic(err)
		}
		// Encode bytes to base36 (digits, lower case letters)
		num := new(big.Int).SetBytes(r)
		str = num.Text(36)
		dbkey = f.Key(str)
		// Check if the id is already used
		if _, err = db.Get([]byte(dbkey)); err == db.ERR_NOT_FOUND {
			break
		}
		tries++
		if tries > 10 {
			// This is very very unlikely, if it happens
			// something is teribly wrong.
			panic("Too many tries to get a unique random id?!!")
		}
	}
	// Pad with zeros to ensure a 5-character id
	for len(str) < 5 {
		str = "0" + str
	}
	return str
}

func (f *Frame) Save() error {
	db.AssertOpen()
	if f.Id == "" {
		newId := f.newFrameId()
		f.Id = newId
		f.PublicKey.FrameId = f.Id
	}

	key := f.Key(f.Id)
	data, err := proto.Marshal(f)
	if err != nil {
		return fmt.Errorf("Error marshaling Frame protobuf: %v\n", err)
	}

	if err = db.Set([]byte(key), data); err != nil {
		return fmt.Errorf("Error saving Frame: %v\n", err)
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

func (f *Frame) FromEndpoint(endpoint string) error {
	parts := strings.Split(endpoint, "/")
	if len(parts) != 3 {
		return errors.New("Expected /f/<frameId> path.")
	}
	frameId := parts[2]
	if f.FromId(frameId) == nil {
		return errors.New("Error fetching frame from db.")
	}
	return nil
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
