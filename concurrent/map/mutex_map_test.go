package _map

import (
	"sync"
	"testing"
)

type MutexMap struct {
	sync.Mutex
	Map map[int]int
}

func (m *MutexMap) readMap(index int) (int, bool) {
	m.Lock()
	defer m.Unlock()
	value, ok := m.Map[index]
	return value, ok
}

func (m *MutexMap) writeMap(index int, value int) {
	m.Lock()
	defer m.Unlock()
	m.Map[index] = value
}

func TestVerifyMutex(t *testing.T) {
	var mMap = &MutexMap{
		Map: make(map[int]int),
	}

	for i := 0; i < 1000; i++ {
		go func() {
			mMap.writeMap(i, i)
		}()

		go mMap.readMap(i)
	}
}


