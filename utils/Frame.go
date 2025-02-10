package utils

import (
	"fmt"

	db "github.com/vrypan/farma/localdb"
)

type Frame struct {
	Id       int
	Name     string
	Desc     string
	Domain   string
	Endpoint string
}

func NewFrame() *Frame {
	return &Frame{}
}

func (f *Frame) FromEndpoint(e string) error {
	db.AssertOpen()

	rows, err := db.Instance.Query("SELECT id, name, desc, domain FROM frames where endpoint=?", e)
	if err != nil {
		return err
	}

	if !rows.Next() {
		return fmt.Errorf("Frame not found")
	}

	if err = rows.Scan(&f.Id, &f.Name, &f.Desc, &f.Domain); err != nil {
		return fmt.Errorf("Frame not found")
	}
	f.Endpoint = e

	return nil
}

func (f Frame) Save() error {
	/// TBA
	return nil
}
