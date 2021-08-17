package main

import (
	"fmt"
	"math/rand"
	"os"
	"provy/load-balancer/usecase"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(4)
	logger := logrus.WithFields(logrus.Fields{
		"service": "reverse-proxy",
	})
	members := []usecase.Member{
		{
			Key:     "/v1",
			Replica: 4000,
		},
		{
			Key:     "/v2",
			Replica: 2000,
		},
		{
			Key:     "/v3",
			Replica: 2000,
		},
		{
			Key:     "/v4",
			Replica: 1000,
		},
	}

	lb := usecase.NewLoadBalancerUsecase(logger, members)

	keyCount := 1000000
	distribution := make(map[string]int)
	input := make([]byte, 4)
	for i := 0; i < keyCount; i++ {
		rand.Read(input)
		member := lb.Locate(string(input))
		distribution[member]++
	}
	for member, count := range distribution {
		fmt.Printf("member: %s, key count: %d\n", member, count)
	}
}
