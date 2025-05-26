package main

import "fmt"

// enforcing go:noinline directive to have them preserved in the assembly output

type node struct {
	value int   // 8 bytes at offset 0
	left  *node // 8 bytes at offset 8
	right *node // 8 bytes at offset 16
}

//go:noinline
func (n *node) insert(val int) *node {
	if n == nil {
		// Creates new node on the Heap.
		// Depending on the use-case but for this specific use case of where we have search functionality,
		// this exhibits poor spatial locality
		return &node{value: val}
	}

	if val < n.value {
		n.left = n.left.insert(val)
	} else if val > n.value {
		n.right = n.right.insert(val)
	}

	return n
}

//go:noinline
func (n *node) search(val int) bool {
	if n == nil {
		return false
	}
	if val == n.value {
		return true
	}
	if val < n.value {
		return n.left.search(val)
	}
	return n.right.search(val)
}

type contiguousBST struct {
	data []int
	used []bool
	size int
}

//go:noinline
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

//go:noinline
func (cbst *contiguousBST) search(val int) bool {
	index := 0
	for index < len(cbst.data) && cbst.used[index] {
		if val == cbst.data[index] {
			return true
		}
		if val < cbst.data[index] {
			index = 2*index + 1
		} else {
			index = 2*index + 2
		}
	}
	return false
}

func main() {
	// Force the compiler to include functions

	// ptr BST
	var root *node
	root = root.insert(5)
	root = root.insert(3)
	root = root.insert(7)
	found1 := root.search(3)

	// Array BST usage
	bst := &contiguousBST{
		data: make([]int, 100),
		used: make([]bool, 100),
		size: 0,
	}
	bst.insert(5)
	bst.insert(3)
	bst.insert(7)
	found2 := bst.search(3)

	if found1 && found2 {
		fmt.Println("found")
	}
}
