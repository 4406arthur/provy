package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32
type Hash func(data []byte) uint32

// Map constains all hashed keys
type Map struct {
	hash Hash
	//replicas int
	keys    []int // Sorted
	hashMap map[int]string
}

type Member struct {
	key     string
	replica int
}

// New creates a Map instance
func New(fn Hash) *Map {
	m := &Map{
		//replicas: replicas,
		hash:    fn,
		hashMap: make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add adds some keys to the hash.
func (m *Map) Add(members ...Member) {
	for _, member := range members {
		for i := 0; i < member.replica; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + member.key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = member.key
		}
	}
	sort.Ints(m.keys)
}

// Get gets the closest item in the hash to the provided key.
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// Binary search for appropriate replica.
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}
