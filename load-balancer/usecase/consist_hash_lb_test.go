package usecase

import (
	"os"
	"strconv"
	"testing"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Entry

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(4)
	logger = logrus.WithFields(logrus.Fields{})
}

func TestLocate(t *testing.T) {
	members := []Member{
		{
			Key:     "node1",
			Replica: 40,
		},
		{
			Key:     "node2",
			Replica: 20,
		},
		{
			Key:     "node3",
			Replica: 20,
		},
	}
	lb := NewLoadBalancerUsecase(logger, members)

	testCases := map[string]string{
		"2":  "node2",
		"11": "node3",
		"23": "node1",
		"27": "node3",
	}

	for k, v := range testCases {
		if lb.Locate(k) != v {
			t.Errorf("Asking for %s, should have yielded %s but %s", k, lb.Locate(k), v)
		}
	}
}

func BenchmarkLocate(b *testing.B) {

	// Given the above hash function, this will give replicas with "hashes":
	// 2, 4, 6, 12, 14, 16, 22, 24, 26
	members := []Member{
		{
			Key:     "/v1",
			Replica: 800,
		},
		{
			Key:     "/v2",
			Replica: 200,
		},
		{
			Key:     "/v3",
			Replica: 200,
		},
	}

	lb := NewLoadBalancerUsecase(logger, members)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		lb.Locate(key)
	}
}
