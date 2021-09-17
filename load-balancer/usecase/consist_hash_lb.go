package usecase

import (
	"hash/crc32"
	"sort"
	"strconv"

	"github.com/sirupsen/logrus"
)

type CHLoadBalancer struct {
	Logger *logrus.Entry
	Ring   Ring
}

type Hash func(data []byte) uint32

// Ring constains all hashed keys
type Ring struct {
	Hash    Hash
	Keys    []int // Sorted
	VNodes  int
	HashMap map[int]string
}

type Member struct {
	Key     string `mapstructure:"key"`
	Replica int    `mapstructure:"replica"`
}

func NewLoadBalancerUsecase(logger *logrus.Entry, members []Member) *CHLoadBalancer {
	// logger.Infof("LB algo: %s \n LB rule %v , algo, rule")
	ring := Ring{
		Hash:    crc32.ChecksumIEEE,
		HashMap: make(map[int]string),
	}

	// Add adds some keys to the hash.
	for _, member := range members {
		for i := 0; i < member.Replica; i++ {
			hash := int(ring.Hash([]byte(strconv.Itoa(i) + member.Key)))
			ring.Keys = append(ring.Keys, hash)
			ring.HashMap[hash] = member.Key
		}
	}
	sort.Ints(ring.Keys)
	ring.VNodes = len(ring.Keys)

	return &CHLoadBalancer{
		Logger: logger,
		Ring:   ring,
	}
}

// Get gets the closest item in the hash to the provided key.
func (lb *CHLoadBalancer) Locate(key string) string {
	if len(lb.Ring.Keys) == 0 {
		return ""
	}

	hash := int(lb.Ring.Hash([]byte(key)))
	// Binary search for appropriate replica.
	idx := sort.Search(len(lb.Ring.Keys), func(i int) bool {
		return lb.Ring.Keys[i] >= hash
	})

	if idx == lb.Ring.VNodes {
		idx = 0
	}

	return lb.Ring.HashMap[lb.Ring.Keys[idx]]
}

// testing ussage
func new(fn Hash) *Ring {
	m := &Ring{
		//replicas: replicas,
		Hash:    fn,
		HashMap: make(map[int]string),
	}
	if m.Hash == nil {
		m.Hash = crc32.ChecksumIEEE
	}
	return m
}

// testing ussage
func (r *Ring) add(members ...Member) {
	for _, member := range members {
		for i := 0; i < member.Replica; i++ {
			hash := int(r.Hash([]byte(strconv.Itoa(i) + member.Key)))
			r.Keys = append(r.Keys, hash)
			r.HashMap[hash] = member.Key
		}
	}
	sort.Ints(r.Keys)
}

// testing ussage
func (r *Ring) get(key string) string {
	if len(r.Keys) == 0 {
		return ""
	}

	hash := int(r.Hash([]byte(key)))
	// Binary search for appropriate replica.
	idx := sort.Search(len(r.Keys), func(i int) bool {
		return r.Keys[i] >= hash
	})

	return r.HashMap[r.Keys[idx%len(r.Keys)]]
}
