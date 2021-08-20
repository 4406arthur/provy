package domain

import "net/http"

type ReverseProxy interface {
	Handler() http.HandlerFunc
}
type LoadBalancerUsecase interface {
	Locate(key string) string
}
