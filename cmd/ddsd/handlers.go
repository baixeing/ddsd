package main

import (
	"log"
	"net/http"

	"encoding/json"

	"github.com/baixeing/ddsd/storage"
)

func List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	s, err := storage.NewStorage("/tmp/ddsd")
	if err != nil {
		log.Fatalln(err)
	}
	defer s.Close()

	files := s.List("")

	if err := json.NewEncoder(w).Encode(files); err != nil {
		log.Println(err)
	}
}

func Status(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func Put(w http.ResponseWriter, r *http.Request) {
	log.Println(*r)
}
