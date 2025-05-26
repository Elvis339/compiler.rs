package main

import "math/rand"

type graphNode struct {
	id        int
	neighbors []*graphNode
	visited   bool
}

// bfs visit all nodes in the graph starting from the root node
func bfs(root *graphNode) int {
	visited := make(map[int]bool)
	queue := []*graphNode{root}
	count := 0

	for len(queue) > 0 {
		current := queue[0] // create new slice on each iteration
		queue = queue[1:]

		if visited[current.id] {
			continue
		}

		visited[current.id] = true
		count++

		// pointer chasing trough the heap
		for _, neighbor := range current.neighbors {
			if !visited[neighbor.id] {
				queue = append(queue, neighbor)
			}
		}
	}

	return count
}

func createGraph(size int) *graphNode {
	nodes := make([]*graphNode, size)
	for i := 0; i < size; i++ {
		nodes[i] = &graphNode{id: i}
	}

	for i := 0; i < size; i++ {
		connections := rand.Intn(10) + 1
		for j := 0; j < connections; j++ {
			target := rand.Intn(size)
			nodes[i].neighbors = append(nodes[i].neighbors, nodes[target])
		}
	}

	return nodes[0]
}
