package _map

import (
	"log"
	"strconv"
	"testing"
)

func BenchmarkMutexMap(t *testing.B) {
	var mMap = &MutexMap{
		Map: make(map[string]interface{}),
	}

	for i := 0; i < 1000; i++ {
		key := strconv.Itoa(i)
		go func() {
			mMap.writeMap(key, i)
		}()

		go func() {
			value, ok := mMap.readMap(key)
			if ok {
				log.Print(value)
			}
		}()
	}
}

func BenchmarkRWMutexMap(b *testing.B) {
	var mMap = &RWMutexMap{
		Map: make(map[string]interface{}),
	}

	for i := 1; i < 1000; i++ {
		key := strconv.Itoa(i)
		go func() {
			mMap.Set(key, i)
		}()

		go func() {
			value, ok := mMap.Get(key)
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
