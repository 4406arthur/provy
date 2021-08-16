package cmd

import (
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

type ReverseProxy struct {
	Protocol string
	Backend  string
	//TODO: is possible have a template to config custom handler logic ???
	Handler *http.ServeMux
}

func NewReveseProxy(protocol, backend string, director func(*http.Request)) *ReverseProxy {
	//setup handler
	mux := http.NewServeMux()
	rp := &httputil.ReverseProxy{
		Director: director,
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				return net.DialTimeout(protocol, backend, 5*time.Second)
			},
		},
	}

	mux.Handle("/", rp)

	return &ReverseProxy{
		Protocol: protocol,
		Backend:  backend,
		Handler:  mux,
	}
}
