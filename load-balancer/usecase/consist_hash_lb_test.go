package usecase

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	ring := new(func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// Given the above hash function, this will give replicas with "hashes":
	// 2, 4, 6, 12, 14, 16, 22, 24, 26
	members := []Member{
		{
			Key:     "6",
			Replica: 40,
		},
		{
			Key:     "4",
			Replica: 20,
		},
		{
			Key:     "2",
			Replica: 20,
		},
	}
	ring.add(members...)

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		if ring.get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	// Adds 8, 18, 28
	ring.add(Member{
		Key:     "8",
		Replica: 20,
	})

	// 27 should now map to 8.
	testCases["27"] = "8"

	for k, v := range testCases {
		if ring.get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

}
