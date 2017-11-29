package main

import (
	"sync"
	"testing"
	"runtime"
)

var global = make(chan string, 10000)

func BenchmarkRWMutexMapGetConcurrent(b *testing.B) {
	m := map[string]string{
		"foo": "bar",
	}
	mu := sync.RWMutex{}
	wg := new(sync.WaitGroup)
	workers := runtime.GOMAXPROCS(0)
	each := b.N / workers
	wg.Add(workers)
	b.ResetTimer()
	for i := 0; i < workers; i++ {
		go func() {
			var tmp string
			for j := 0; j < each; j++ {
				mu.RLock()
				tmp, _ = m["foo"]
				mu.RUnlock()
			}
			global <- tmp
			wg.Done()
		}()
	}
	wg.Wait()
}

// go test -cpu=4 -run=XXX -bench=BenchmarkRegularParallel -benchtime=5s
func BenchmarkRWMutexMapGetParallel(b *testing.B) {
	m := map[string]string{
		"foo": "bar",
	}
	mu := sync.RWMutex{}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var tmp string
		for pb.Next() {
			mu.RLock()
			tmp, _ = m["foo"]
			mu.RUnlock()
		}
		global <- tmp
	})
}

func BenchmarkSyncMapGetConcurrent(b *testing.B) {
	m := new(sync.Map)
	m.Store("foo", "bar")
	wg := new(sync.WaitGroup)
	workers := runtime.GOMAXPROCS(0)
	each := b.N / workers
	wg.Add(workers)
	b.ResetTimer()
	for i := 0; i < workers; i++ {
		go func() {
			var tmp string
			for j := 0; j < each; j++ {
				a, _ := m.Load("foo")
				tmp = a.(string)
			}
			global <- tmp
			wg.Done()
		}()
	}
	wg.Wait()
}

// go test -cpu=4 -run=XXX -bench=BenchmarkRegularParallel -benchtime=5s
func BenchmarkSyncMapGetParallel(b *testing.B) {
	m := new(sync.Map)
	m.Store("foo", "bar")
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var tmp string
		for pb.Next() {
			a, _ := m.Load("foo")
			tmp = a.(string)
		}
		global <- tmp
	})
}
