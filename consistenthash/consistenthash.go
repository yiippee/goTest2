package main

import (
	"hash/crc32"
	"strconv"
)

func main() {

}

type Hash func(data []byte) uint32

type Map struct {
	hash    Hash
	repli   int
	hashMap map[int]string
}

func New(repli int, fn Hash) *Map {
	m := &Map{
		hash:    fn,
		repli:   repli,
		hashMap: make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.repli; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.hashMap[hash] = key
		}
	}
}
