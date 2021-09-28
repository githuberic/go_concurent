package _map

import (
	"log"
	"testing"
)

func BenchmarkMutexMap(t *testing.B) {
	var mMap = &MutexMap{
		Map: make(map[int]int),
	}

	for i := 0; i < 1000; i++ {
		go func() {
			mMap.writeMap(i, i)
		}()

		go func() {
			value, ok := mMap.readMap(i)
			if ok {
				log.Print(value)
			}
		}()
	}
}

func BenchmarkRWMutexMap(b *testing.B) {
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

func BenchmarkChMap(t *testing.B) {
	mMap := NewChannelMap()

	for i := 0; i < 1000; i++ {
		go func() {
			mMap.add(i, i)
		}()

		go func() {
			value := mMap.find(i)
			if value > 0 {
				log.Printf("Value=%d", value)
			}
		}()
	}
}

func BenchmarkChPoolMap(t *testing.B) {
	mMap := NewChannelPoolMap()

	for i := 0; i < 1000; i++ {
		go func() {
			mMap.add(i, i)
		}()

		go func() {
			value := mMap.find(i)
			if value != nil && *value > 0 {
				log.Print(*value)
			}
		}()
	}
}
