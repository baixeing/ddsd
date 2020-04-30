package storage

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/segmentio/ksuid"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/dustin/go-humanize"
)

type File struct {
	Model
	UID         string `gorm:"unique"`
	Name        string
	Path        string
	FileSize    uint64
	ChunkSize   uint64
	ChunksTotal uint64
	Perm        os.FileMode
	Chunks      []*Chunk
	Checksum    string
	ContentType string
}

func NewFile(name string, cs uint64) (*File, error) {
	// Chunks
	info, err := os.Lstat(name)
	if err != nil {
		return nil, err
	}

	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("not a regular file: %s", info.Name())
	}

	size := uint64(info.Size())
	n := size / cs
	if size%cs != 0 {
		n++
	}
	chunks := make([]*Chunk, n)

	// Checksum, ContentType
	fd, err := os.OpenFile(name, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	p := make([]byte, 512)
	if _, err = fd.Read(p); err != nil {
		return nil, err
	}
	contentType := http.DetectContentType(p)

	h := md5.New()
	if _, err := io.Copy(h, fd); err != nil {
		return nil, err
	}
	checksum := fmt.Sprintf("%x", h.Sum(nil))

	return &File{
		UID:         ksuid.New().String(),
		Name:        filepath.Base(name),
		Path:        filepath.Dir(name),
		Perm:        info.Mode().Perm(),
		FileSize:    size,
		ChunkSize:   cs,
		ChunksTotal: n,
		Chunks:      chunks,
		Checksum:    checksum,
		ContentType: contentType,
	}, nil
}

func (f *File) Info(db *gorm.DB) string {
	var cn int

	db.Model(&Chunk{}).Where("file = ?", f.UID).Count(&cn)

	return fmt.Sprintf("%s\t%s\t%s\t%s\t%d\t%d\t%s\t%s",
		f.UID, f.Name, f.Path, humanize.IBytes(f.FileSize), cn, f.ChunksTotal, humanize.IBytes(f.ChunkSize), f.ContentType)
}

func (f *File) Put(path string, db *gorm.DB) error {
	savePath := filepath.Join(path, f.UID)

	// [TODO] Set dir permission from config
	if err := os.MkdirAll(savePath, 0700); err != nil {
		return err
	}

	fd, err := os.OpenFile(filepath.Join(f.Path, f.Name), os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer fd.Close()

	buff := make([]byte, f.ChunkSize)
	i := 0
	for n, err := fd.Read(buff); err == nil; i++ {
		f.Chunks[i] = NewChunk(f.UID, uint64(i))

		if err := f.Chunks[i].Put(buff[:n], savePath, db); err != nil {
			return err
		}

		n, err = fd.Read(buff)
		if err != nil && err != io.EOF {
			return err
		}

	}
	db.Create(f)

	return nil
}

func (f *File) Delete(path string, db *gorm.DB) error {
	from := filepath.Join(path, f.UID)
	for _, c := range f.Chunks {
		db.Delete(c)
	}
	db.Delete(f)
	return os.RemoveAll(from)
}
