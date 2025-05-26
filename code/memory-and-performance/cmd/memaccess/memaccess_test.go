package main

import (
	"testing"
)

// go test -bench=BenchmarkPointerBST -memprofile=mem.prof
// go test -bench=BenchmarkContiguousBST -memprofile=mem.prof

// go tool pprof mem.prof
// top10

// go test -bench=. -benchmem -count=6
func BenchmarkPointerBST(b *testing.B) {
	values, searches := setup(treeSize)
	b.ResetTimer()

	var root *node
	for i := 0; i < b.N; i++ {
		for _, val := range values {
			root = root.insert(val)
		}

		s := false
		for _, search := range searches {
			s = root.search(search)
		}
		_ = s
	}
}

func BenchmarkContiguousBST(b *testing.B) {
	values, searches := setup(treeSize)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cbst := newContiguousBST(treeSize)
		for _, val := range values {
			cbst.insert(val)
		}

		s := false
		for _, search := range searches {
			s = cbst.search(search)
		}
		_ = s
	}
}
