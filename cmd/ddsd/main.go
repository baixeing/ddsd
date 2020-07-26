package main

import (
	"log"
	"time"

	"net/http"

	"fmt"

	"math/rand"

	"github.com/grandcat/zeroconf"
	"github.com/segmentio/ksuid"
)

const (
	Service = "_ddsd._tcp"
	Domain  = "local."
)

func NewDiscovery(port int) (*zeroconf.Server, error) {
	instance := ksuid.New().String()
	return zeroconf.Register(
		instance,
		Service,
		Domain,
		port,
		nil,
		nil,
	)
}

func main() {
	r := rand.NewSource(time.Now().UnixNano())
	port := int(r.Int63()%3000 + 2000)

	server, err := NewDiscovery(port)
	if err != nil {
		log.Fatalln(err)
	}
	defer server.Shutdown()

	router := NewRouter()
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
