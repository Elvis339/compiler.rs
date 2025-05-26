package main

import (
	"fmt"
	"sync/atomic"
)

type spinLock struct {
	flag int32
}

func newSpinLock() *spinLock {
	return &spinLock{flag: 0}
}

func (s *spinLock) Lock() {
	for !atomic.CompareAndSwapInt32(&s.flag, 0, 1) {
		// Locked try again.
		// We keep spinning and wasting resources, effectively starving the CPU.
		// Ideally, we want to try X number of times and if unsuccessful hint to the CPU, we are ready to yield - allowing progress on other tasks.

		// Uncomment this for better performance
		//runtime.Gosched()
	}
}

func (s *spinLock) Unlock() {
	atomic.StoreInt32(&s.flag, 0)
}

func main() {
	fmt.Println("spinlock")
}
