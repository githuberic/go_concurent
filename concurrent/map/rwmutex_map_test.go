package _map

import (
	"log"
	"sync"
	"testing"
)

type RWMutexMap struct {
	sync.RWMutex
	Map map[int]int
}

func (m *RWMutexMap) Get(index int) (int, bool) {
	m.RLock()
	defer m.RUnlock()
	value, ok := m.Map[index]
	return value, ok
}

func (m *RWMutexMap) Set(index int, value int) {
	m.Lock()
	defer m.Unlock()
	m.Map[index] = value
}

func TestVerifyRWMutex(t *testing.T) {
	var mMap = &RWMutexMap{
		Map: make(map[int]int),
	}

	for i := 1; i < 1000; i++ {
		go func() {
			mMap.Set(i, i)
		}()

		go func() {
			value, ok := mMap.Get(i)
			if ok {
				log.Print(value)
			}
		}()
	}
}
