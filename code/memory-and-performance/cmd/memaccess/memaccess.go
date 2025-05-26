package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

// node represents a traditional pointer-based binary search tree (BST)
// Each node is allocated separately on the heap, creating scattered memory layout
type node struct {
	value int   // The stored value
	left  *node // Pointer to left child (smaller values)
	right *node // Pointer to right child (larger values)
}

// insert recursively adds a value to the BST
func (n *node) insert(val int) *node {
	if n == nil {
		// Heap allocation: creates new node at unpredictable memory address
		// Depending on the use-case but for this specific use case of where we have search
		// this creates poor spatial locality
		return &node{value: val}
	}

	if val < n.value {
		n.left = n.left.insert(val)
	} else if val > n.value {
		n.right = n.right.insert(val)
	}

	return n
}

// search traverses the BST following pointer chains
// Each pointer dereference (n.left, n.right) likely causes cache miss
// due to scattered memory layout of nodes
func (n *node) search(val int) bool {
	if n == nil {
		return false
	}

	if val == n.value {
		return true
	}

	if val < n.value {
		// Follow left pointer - random memory jump, likely cache miss
		return n.left.search(val)
	}
	// Follow right pointer - random memory jump, likely cache miss
	return n.right.search(val)
}

// contiguousBST implements BST using array-based heap indexing
// All data stored in contiguous memory for better spatial locality
// Uses heap property: parent at i, left child at 2*i+1, right child at 2*i+2
type contiguousBST struct {
	data []int  // Contiguous array storing all values
	used []bool // Tracks which array positions are occupied
	size int    // Current number of elements
}

func newContiguousBST(size int) *contiguousBST {
	return &contiguousBST{
		data: make([]int, size),
		used: make([]bool, size),
		size: 0,
	}
}

// insert adds value using array indexing instead of pointer following
// Uses arithmetic (2*i+1, 2*i+2) instead of pointer dereferencing
// Accesses predictable memory locations - cache-friendly
func (cbst *contiguousBST) insert(val int) {
	if cbst.size == 0 {
		// First element goes at root (index 0)
		cbst.data[0] = val
		cbst.used[0] = true
		cbst.size++
		return
	}

	index := 0
	l := len(cbst.data)
	for {
		if index >= l {
			return
		}

		if !cbst.used[index] {
			cbst.used[index] = true
			cbst.data[index] = val
			cbst.size++
			return
		}

		// Navigate using heap indexing (arithmetic, no pointer dereferencing)
		if val < cbst.data[index] {
			index = 2*index + 1 // Left child index calculation
		} else if val > cbst.data[index] {
			index = 2*index + 2 // Right child index calculation
		} else {
			return // Duplicate value
		}
	}
}

// search traverses BST using array indexing instead of pointer chasing
// Each access is to predictable memory location within contiguous array
// Much more cache-friendly than random pointer dereferencing
func (cbst *contiguousBST) search(val int) bool {
	index := 0
	l := len(cbst.data)
	for index < l && cbst.used[index] {
		if val == cbst.data[index] {
			return true
		}
		// Navigate using arithmetic instead of pointer dereferencing
		if val < cbst.data[index] {
			index = 2*index + 1 // Left child - simple calculation
		} else {
			index = 2*index + 2 // Right child - simple calculation
		}
	}
	return false
}

const treeSize = 1_000_000

func setup(n int) ([]int, []int) {
	values := make([]int, n)
	for i := 0; i < n; i++ {
		values[i] = rand.Intn(100_000)
	}

	searches := make([]int, n)
	for i := 0; i < n; i++ {
		searches[i] = rand.Intn(100_000)
	}

	return values, searches
}

// go run . -v ptr
// go run . -v array

func main() {
	version := flag.String("v", "ptr", "BST version: ptr (scattered) or array (contiguous)")
	flag.Parse()

	values, searches := setup(treeSize)

	start := time.Now()
	switch *version {
	case "ptr":
		var root *node
		for _, val := range values {
			root = root.insert(val)
		}

		for _, search := range searches {
			root.search(search)
		}
		fmt.Printf("BST(%s): %s\n", *version, time.Since(start))
	default: // "array" case
		cbst := newContiguousBST(treeSize * 2)
		for _, val := range values {
			cbst.insert(val)
		}

		for _, search := range searches {
			cbst.search(search)
		}
		fmt.Printf("BST(%s): %s\n", *version, time.Since(start))
	}
}
