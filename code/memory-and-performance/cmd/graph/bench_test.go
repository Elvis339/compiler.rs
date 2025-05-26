package main

import (
	"runtime"
	"testing"
)

func BenchmarkGraph(b *testing.B) {
	graph := createGraph(2_000_000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(bfs(graph))
	}
}

func BenchmarkCompact(b *testing.B) {
	graph := createCompactGraph(2_000_000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(graph.bfs(0))
	}
}
