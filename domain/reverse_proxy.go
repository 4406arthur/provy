package domain

type LoadBalancerUsecase interface {
	Locate(key string) string
}
