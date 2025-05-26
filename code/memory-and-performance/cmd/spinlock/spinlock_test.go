package main

import (
	"sync"
	"testing"
)

// go test -bench=. -benchmem -cpu=1,2,4,8

func BenchmarkContentionLevels(b *testing.B) {
	tests := []struct {
		name       string
		goroutines int
		workCycles int
	}{
		{"LowContention", 2, 1000},
		{"MediumContention", 4, 100},
		{"HighContention", 8, 10},
		{"VeryHighContention", 16, 1},
	}

	for _, tt := range tests {
		b.Run(tt.name+"_SpinLock", func(b *testing.B) {
			benchmarkWithContention(b, newSpinLock(), tt.goroutines, tt.workCycles)
		})

		b.Run(tt.name+"_Mutex", func(b *testing.B) {
			benchmarkWithContention(b, &sync.Mutex{}, tt.goroutines, tt.workCycles)
		})
	}
}

type Locker interface {
	Lock()
	Unlock()
}

func benchmarkWithContention(b *testing.B, lock Locker, goroutines, workCycles int) {
	var wg sync.WaitGroup
	var sharedCounter int64
	iterations := b.N / goroutines

	b.ResetTimer()

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for j := 0; j < iterations; j++ {
				lock.Lock()

				temp := sharedCounter
				for k := 0; k < workCycles; k++ {
					temp += int64(k + workerID)
				}
				sharedCounter = temp

				lock.Unlock()

				localWork := 0
				for k := 0; k < workCycles*2; k++ {
					localWork += k
				}
				_ = localWork // Prevent optimization
			}
		}(i)
	}

	wg.Wait()
}
