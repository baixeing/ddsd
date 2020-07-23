package storage

import (
	"os"
	"path/filepath"

	"github.com/segmentio/ksuid"

	"github.com/jinzhu/gorm"
)

type Chunk struct {
	Model
	UID     string `gorm:"unique"`
	Seq     uint64 `gorm:"column:seq"`
	FileUID string `gorm:"column:file"`
}

func NewChunk(fileUID string, seq uint64) *Chunk {
	return &Chunk{
		UID:     ksuid.New().String(),
		Seq:     seq,
		FileUID: fileUID,
	}
}

func (c *Chunk) Push(p []byte, path string, db *gorm.DB) error {
	to := filepath.Join(path, c.UID)

	// [TODO] get default permissions from config
	fd, err := os.OpenFile(to, os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		return err
	}

	if _, err = fd.Write(p); err != nil {
		return err
	}

	db.Create(c)
	return fd.Close()
}

func (c *Chunk) Delete(path string, db *gorm.DB) error {
	from := filepath.Join(path, c.UID)
	db.Delete(c)

	return os.Remove(from)
}
