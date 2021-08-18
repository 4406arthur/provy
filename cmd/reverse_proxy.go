package cmd

import (
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type ReverseProxy struct {
	Protocol   string
	Backend    string
	Instrument prometheus.Summary
	//TODO: is possible have a template to config custom handler logic ???
	Proxy *httputil.ReverseProxy
}

func NewReveseProxy(protocol, backend string, instrument prometheus.Summary, director func(*http.Request)) *ReverseProxy {

	reverProxy := &httputil.ReverseProxy{
		Director: director,
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				return net.DialTimeout(protocol, backend, 5*time.Second)
			},
		},
	}

	return &ReverseProxy{
		Protocol:   protocol,
		Backend:    backend,
		Instrument: instrument,
		Proxy:      reverProxy,
	}
}

func (p *ReverseProxy) Handler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		p.Proxy.ServeHTTP(w, r)
		elapsed := float64(time.Since(now)) / float64(time.Microsecond)
		p.Instrument.Observe(elapsed)
	})
}
