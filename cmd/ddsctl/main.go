package main

import (
	"flag"
	"log"
	"os"

	"github.com/baixeing/ddsd/storage"
)

var (
	put    string
	remove string
	list   bool
)

func init() {
	flag.StringVar(&put, "put", "", "file to save")
	flag.BoolVar(&list, "list", false, "list files on storage")
	flag.StringVar(&remove, "remove", "", "remove file with UID")
	flag.Parse()

	n := 0
	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() != f.DefValue {
			n++
		}
	})

	if n != 1 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	s, err := storage.NewStorage("/tmp/ddsd")
	if err != nil {
		log.Fatalln(err)
	}
	defer s.Close()

	switch {
	case put != "":
		if err = s.Put(put, 4096*1024); err != nil {
			log.Fatalln(err)
		}
	case remove != "":
		if err = s.DeleteFile(remove); err != nil {
			log.Fatalln(err)
		}
	case list:
		s.Info("dummy query/filter")
	}
}
