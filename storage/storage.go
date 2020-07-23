package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/jinzhu/gorm"
)

type Storage struct {
	Path string
	db   *gorm.DB
}

func NewStorage(path string) (*Storage, error) {
	if err := os.MkdirAll(path, 0700); err != nil {
		return nil, err
	}

	db, err := gorm.Open("sqlite3", filepath.Join(path, "storage.db"))
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&File{}, &Chunk{})

	return &Storage{
		path, db,
	}, nil
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) pushFile(name string, cs uint64) error {
	f, err := NewFile(name, cs)
	if err != nil {
		return err
	}

	if err = f.Push(s.Path, s.db); err != nil {
		return err
	}

	return nil
}

func (s *Storage) pushDir(name string, cs uint64) error {
	// [TODO] LOG instead of return error
	err := filepath.Walk(name,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Mode().IsRegular() {
				return s.pushFile(path, cs)
			}
			return nil
		},
	)
	return err
}

func (s *Storage) Push(name string, cs uint64) error {
	info, err := os.Lstat(name)
	if err != nil {
		return err
	}

	switch m := info.Mode(); {
	case m.IsRegular():
		return s.pushFile(name, cs)
	case m.IsDir():
		return s.pushDir(name, cs)
	}

	return nil
}

func (s *Storage) Info(_ string) {
	files := make([]File, 0)
	s.db.Find(&files)

	w := tabwriter.NewWriter(os.Stdout, 1, 2, 3, ' ', 0)
	_, _ = fmt.Fprintf(w, "ID\tFILENAME\tPATH\tSIZE\tCONTENT-TYPE\n")

	for _, f := range files {
		_, _ = fmt.Fprintln(w, f.Info(s.db))
	}
	w.Flush()
}

func (s *Storage) Detail(uuid string) {
	f := File{}
	s.db.Find(&f).Where("file = ?", uuid) // [FIXME] SQLite query
	//	s.db.Where("file = ?", uuid).Find(&f)

	w := tabwriter.NewWriter(os.Stdout, 1, 2, 3, ' ', 0)

	_, _ = fmt.Fprintf(w, "ID\tFILENAME\tPATH\tPERMISSION\tSIZE\tCHUNKSIZE\tCHUNKS\tMD5\tCONTENT-TYPE\n")
	_, _ = fmt.Fprintln(w, f.Detail())

	w.Flush()
}

func (s *Storage) DeleteFile(fileUUID string) error {
	f := File{UID: fileUUID}

	s.db.Where(&f).First(&f)
	s.db.Where("file = ?", fileUUID).Find(&f.Chunks)

	return f.Delete(s.Path, s.db)
}
