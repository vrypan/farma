package models

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"

	db "github.com/vrypan/farma/localdb"
	"google.golang.org/protobuf/proto"
)

func NewFrame() *Frame {
	return &Frame{}
}

func (f *Frame) Key(id uint64) string {
	return fmt.Sprintf("f:id:%d", id)
}

func (f *Frame) Save() error {
	db.AssertOpen()
	if f.Id == 0 {
		id, err := db.FrameIdSequence.GetNext()
		if err != nil || id == 0 {
			id, err = db.FrameIdSequence.GetNext()
			if err != nil {
				return err
			}
		}
		f.Id = id
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
	// f:endpoint:<endpoint> and
	// f:name:<name>
	// This makes it very fast to find a Frame by endpoint or name

	endpointKey := fmt.Sprintf("f:endpoint:%s", f.Webhook)
	if err = db.Set([]byte(endpointKey), []byte(key)); err != nil {
		return err
	}
	nameKey := fmt.Sprintf("f:name:%s", f.Name)
	if err = db.Set([]byte(nameKey), []byte(key)); err != nil {
		return err
	}
	if f.PublicKey != nil {
		publicKey := NewPublicKey(f.PublicKey, f.Id).Bytes()
		if err = db.Set(publicKey, []byte(key)); err != nil {
			return err
		}
	}
	return nil
}

func (f *Frame) NewPk() ([]byte, error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("Error generating Frame key pair: %v\n", err)
	}
	f.PublicKey = pubKey
	return privKey, nil
}

func (f *Frame) FromEndpoint(endpoint string) error {
	db.AssertOpen()

	refKey := fmt.Sprintf("f:endpoint:%s", endpoint)
	key, err := db.Get([]byte(refKey))
	if err != nil {
		return err
	}
	frameBytes, err := db.Get(key)
	if err != nil {
		return err
	}
	return proto.Unmarshal(frameBytes, f)
}

// Checks the database and updates f if the frame already exists
func (f *Frame) FromName(name string) error {
	db.AssertOpen()

	refKey := fmt.Sprintf("f:name:%s", name)
	key, err := db.Get([]byte(refKey))
	if err != nil {
		return err
	}
	frameBytes, err := db.Get(key)
	if err != nil {
		return err
	}
	return proto.Unmarshal(frameBytes, f)
}
func (f *Frame) FromId(id uint64) *Frame {
	db.AssertOpen()
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

func (f *Frame) Delete() error {
	db.AssertOpen()

	if f.Id == 0 {
		return fmt.Errorf("Frame.Id=0")
	}

	key := f.Key(f.Id)
	if err := db.Delete([]byte(key)); err != nil {
		return err
	}

	endpointKey := fmt.Sprintf("f:endpoint:%s", f.Webhook)
	if err := db.Delete([]byte(endpointKey)); err != nil {
		return err
	}
	nameKey := fmt.Sprintf("f:name:%s", f.Name)
	if err := db.Delete([]byte(nameKey)); err != nil {
		return err
	}
	return nil
}

func (f *Frame) Update() error {
	db.AssertOpen()

	if f.Id == 0 {
		return fmt.Errorf("Frame.Id=0")
	}

	key := f.Key(f.Id)
	data, err := proto.Marshal(f)
	if err != nil {
		return err
	}

	if err = db.Set([]byte(key), data); err != nil {
		return err
	}

	endpointKey := fmt.Sprintf("f:endpoint:%s", f.Webhook)
	if err = db.Set([]byte(endpointKey), []byte(key)); err != nil {
		return err
	}
	nameKey := fmt.Sprintf("f:name:%s", f.Name)
	if err = db.Set([]byte(nameKey), []byte(key)); err != nil {
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
