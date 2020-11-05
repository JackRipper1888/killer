package mapkit

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

//go test -v -bench=BenchmarkLock -benchmem mapkit_test.go mapkit.go
func BenchmarkLock(b *testing.B) {
	m := NewConcurrentMap(64)
	go func() {
		for i := 0; i < 1000; i++ {
			m.Set(strconv.Itoa(i), i)
		}
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			m.Set(strconv.Itoa(i), i)
		}
	}()

	b.StartTimer()
	for i := 0; i < 1000; i++ {
		k := rand.Int63n(1000)
		b.ResetTimer()
		b.Log(m.Get(strconv.Itoa(int(k))))
		b.StopTimer()
	}
}

//go test -v -bench=BenchmarkRangeMap -benchmem mapkit_test.go mapkit.go
func BenchmarkRangeMap(b *testing.B) {
	m := NewConcurrentMap(64)
	var n int
	for i := 0; i < b.N; i++ {
		n++
		m.Set(strconv.Itoa(i), i)
	}

	b.StartTimer()
	m.Range(func(k string, v interface{}) bool {
		return true
	})
	b.StopTimer()
}

//go:generata
//go test -v -bench=BenchmarkConcurrentMap -benchmem mapkit_test.go mapkit.go
func BenchmarkConcurrentMap(b *testing.B) {
	m := NewConcurrentMap(64)
	var n int
	for i := 0; i < b.N; i++ {
		n++
		m.Set(strconv.Itoa(i), i)
	}

	b.StartTimer()
	v, isExist := m.Get("4")
	if isExist {
		b.Logf("分片map%d个数据中查找%d", n, v.(int))
	}
	b.StopTimer()
}

//go test -v -bench=BenchmarkSyncMap -benchmem mapkit_test.go mapkit.go
func BenchmarkSyncMap(b *testing.B) {
	var m sync.Map
	var n int
	for i := 0; i < b.N; i++ {
		n++
		m.Store(strconv.Itoa(i), i)
	}

	b.StartTimer()
	v, isExist := m.Load("4")
	if isExist {
		b.Logf("sync map%d个数据中查找%d", n, v.(int))
	}
	//m.Range(func(k, v interface{})bool{
	//	return true
	//})
	b.StopTimer()
}

//go test -v -bench=BenchmarkMap -benchmem mapkit_test.go mapkit.go
func BenchmarkMap(b *testing.B) {
	var m = make(map[string]interface{})
	var n int
	for i := 0; i < b.N; i++ {
		n++
		m[strconv.Itoa(i)] = i
	}
	b.StartTimer()
	v, isExist := m["4"]
	if isExist {
		b.Logf("sync map%d个数据中查找%d", n, v.(int))
	}
	b.StopTimer()
}
