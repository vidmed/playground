package main

import (
	"sync"
	"testing"
)

// Used to prevent compiler optimizations to ensure no dead code elimination.
// These ensure our Load functions aren't eliminated because we capture the result.

// go test -cpu=4 -run=XXX -bench=BenchmarkRegularParallel -benchtime=5s
func BenchmarkRegularParallel(b *testing.B) {
	rm := NewRegularIntMap()
	values := populateMap(b.N, rm)

	// Holds our final results, to prevent compiler optimizations.
	globalResultChan = make(chan int, 64)
	//b.SetParallelism(1)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		currentResult := 0
		i := 0
		for pb.Next() {
			currentResult, _ = rm.Load(values[i])
			i++
		}
		globalResultChan <- currentResult
	})
}

// go test -cpu=4 -run=XXX -bench=BenchmarkRegularParallel -benchtime=5s
func BenchmarkSyncParallel(b *testing.B) {
	var sm sync.Map
	values := populateSyncMap(b.N, &sm)

	// Holds our final results, to prevent compiler optimizations.
	globalResultChan = make(chan int, 64)
	//b.SetParallelism(1)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		currentResult := 0
		i := 0
		for pb.Next() {
			r, _ := sm.Load(values[i])
			currentResult = r.(int)
			i++
		}
		globalResultChan <- currentResult
	})
}
