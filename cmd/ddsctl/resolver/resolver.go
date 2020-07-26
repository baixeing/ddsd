package resolver

import (
	"context"
	"fmt"
	"time"

	"net/url"

	"net"

	"strconv"

	"net/http"

	"github.com/grandcat/zeroconf"
)

func Find(timeout time.Duration) ([]url.URL, error) {
	r, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, err
	}

	nodes := make(chan *zeroconf.ServiceEntry)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = r.Browse(ctx, "_ddsd._tcp", "local.", nodes)
	if err != nil {
		return nil, err
	}

	endpoints := make([]url.URL, 0)

	node := <-nodes
	if node == nil || (len(node.AddrIPv6) == 0 && len(node.AddrIPv4) == 0) {
		return nil, fmt.Errorf("no ddsd server found")
	}

	port := strconv.Itoa(node.Port)
	u := func(addr net.IP) url.URL {
		return url.URL{
			Scheme: "http",
			Host:   net.JoinHostPort(addr.String(), port),
		}
	}

	for _, addr := range node.AddrIPv6 {
		endpoints = append(endpoints, u(addr))
	}

	for _, addr := range node.AddrIPv4 {
		endpoints = append(endpoints, u(addr))
	}

	return endpoints, nil
}

func Endpoint(timeout time.Duration) (*url.URL, error) {
	endpoints, err := Find(timeout)
	if err != nil {
		return nil, err
	}

	var r *http.Response
	for _, endpoint := range endpoints {
		if r, err = http.Get(endpoint.String() + "/status"); err == nil && r.StatusCode == http.StatusOK {
			return &endpoint, nil

		}
	}

	if err == nil {
		err = fmt.Errorf("no ddsd alive")
	}

	return nil, err
}
