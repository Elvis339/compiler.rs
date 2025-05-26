package main

import (
	"math/rand"
	"sync"
)

type compactGraph struct {
	nodes []compactNode
	size  int
}

type compactNode struct {
	id        int
	neighbors []int
}

// Pre-allocated pools to avoid allocations in hot path
var (
	queuePool = sync.Pool{
		New: func() interface{} {
			return make([]int, 0, 1000)
		},
	}

	visitedPool = sync.Pool{
		New: func() interface{} {
			return make(map[int]bool, 100000)
		},
	}
)

// GC-friendly BFS using indices and object pools
func (g *compactGraph) bfs(startID int) int {
	visited := visitedPool.Get().(map[int]bool)
	queue := queuePool.Get().([]int)

	for k := range visited {
		delete(visited, k)
	}
	queue = queue[:0]

	defer func() {
		visitedPool.Put(visited)
		queuePool.Put(queue)
	}()

	queue = append(queue, startID)
	count := 0

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}

		visited[current] = true
		count++

		node := &g.nodes[current]
		for _, neighborID := range node.neighbors {
			if !visited[neighborID] {
				queue = append(queue, neighborID)
			}
		}
	}

	return count
}

func createCompactGraph(size int) *compactGraph {
	nodes := make([]compactNode, size)

	for i := 0; i < size; i++ {
		nodes[i].id = i
		expectedConnections := 5 // Average connections
		nodes[i].neighbors = make([]int, 0, expectedConnections)
	}

	// Build connections
	for i := 0; i < size; i++ {
		connections := rand.Intn(10) + 1
		for j := 0; j < connections; j++ {
			target := rand.Intn(size)
			nodes[i].neighbors = append(nodes[i].neighbors, target)
		}
	}

	return &compactGraph{
		nodes: nodes,
		size:  size,
	}
}
