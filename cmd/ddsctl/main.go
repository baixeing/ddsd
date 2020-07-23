package main

import (
	"flag"
	"log"
	"os"

	"github.com/baixeing/ddsd/storage"
)

var (
	detail string
	push   string
	pull   string
	remove string
	list   bool
)

func init() {
	flag.StringVar(&detail, "detail", "", "file UID details")
	flag.StringVar(&push, "push", "", "file to push to DDSD")
	flag.StringVar(&pull, "pull", "", "file UID to pull from DDSD")
	flag.BoolVar(&list, "ls", false, "list files on storage")
	flag.StringVar(&remove, "rm", "", "remove file with UID")
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
	case detail != "":
		s.Detail(detail)
	case push != "":
		if err = s.Push(push, 4096*1024); err != nil {
			log.Fatalln(err)
		}
	case pull != "":
		log.Println("TODO: ", pull)
	case remove != "":
		if err = s.DeleteFile(remove); err != nil {
			log.Fatalln(err)
		}
	case list:
		s.Info("dummy query/filter")
	}
}
