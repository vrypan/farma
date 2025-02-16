package localdb

import "github.com/dgraph-io/badger/v4"

type Sequence struct {
	name string
	seq  *badger.Sequence
}

func NewSequence(name string, size uint64) *Sequence {
	AssertOpen()
	c := &Sequence{}
	seq, err := db.GetSequence([]byte(name), size)
	if err != nil {
		return nil
	}
	c.seq = seq
	return c
}

func (s Sequence) Release() {
	s.seq.Release()
}

func (s Sequence) GetNext() (uint64, error) {
	return s.seq.Next()
}
