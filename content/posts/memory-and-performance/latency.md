+++
date = '2025-05-24T10:33:25+04:00'
draft = false
title = 'Memory and Performance: Latency'
tags = ['go', 'rust', 'software-performance', 'profiling', 'memory', 'arm', 'cpu']
categories = ['software', 'software-performance']
series = ['memory-and-performance']
description = 'Learn the importance of memory latency and how to reduce its impact. Identify how being memory aware (or unaware) affects performance.'
+++

Memory latency is one of the factors impacting application performance. Join me on a quest to understand memory latency, how to measure it, and knowing when it can be improved.

For this type of posts people usually use C/C++ but I'm going to stick to Go and occasionally bring in Rust for comparison since I'm mostly using these two languages and I want to show you the tools available to you.

We're going to mix in both top-down and bottom-up approaches in this series of blog posts. The goal of this series is to learn how to identify memory latency and act uppon those signals.

# CPU Primer

CPUs can't store data - they need to communicate with memory. There are different types of memory involved, each with different latencies:

| Memory Type | Latency | Size |
|-------------|---------|------|
| L1 Cache | 1-3 cycles | 32-64KB |
| L2 Cache | 10-25 cycles | 256KB-1MB |
| L3 Cache | 40-75 cycles | 8-32MB |
| Main Memory (RAM) | 200-400 cycles | GBs |
| SSD | 50,000+ cycles | TBs |

CPUs work in cycles. A CPU cycle, also known as a clock cycle or machine cycle, is the basic unit of time in a CPU. It represents one complete operation of the CPU's internal clock.

## The Instruction Pipeline

Without pipeline optimization, CPU operations happen sequentially in this order:

1. **Fetching** - retrieving the instruction from memory
2. **Decoding** - interpreting the instruction and determining required actions
3. **Executing** - performing the operation specified by the instruction
4. **Memory Access** (if needed) - fetching data from memory
5. **Write Back** (if needed) - writing the result back to memory

If each stage takes one clock cycle, a single instruction would take 5 cycles to complete. In a simple sequential model, processing 1000 instructions would take 5000 cycles.

## How Pipelining Changes Everything

CPU pipelining transforms instruction processing from a sequential to simultaneous. Instead of waiting for each instruction to complete entirely before starting next, the CPU divides instruction processing into discrete stages and processes them simultaneously.

Here's how it works: while one instruction is being executed (stage 3), the next instruction can be decoded (stage 2), and the one after that can be fetched (stage 1). This creates a pipeline where multiple instructions are "in flight" at different stages.

```
Cycle 1: Fetch A
Cycle 2: Fetch B, Decode A  
Cycle 3: Fetch C, Decode B, Execute A
Cycle 4: Fetch D, Decode C, Execute B, Memory Access A
Cycle 5: Fetch E, Decode D, Execute C, Memory Access B, Write Back A
```

After the initial 5-cycle startup delay, the CPU completes one instruction per cycle instead of one every 5 cycles - a 5x improvement in throughput. Processing 1000 instructions now takes approximately 1004 cycles instead of 5000.

## Pipeline Hazards and Memory Interaction

The pipeline's efficiency depends heavily on memory access patterns. Here's some of the hazards that can disrupt the smooth flow:

**Data Hazards** occur if one instruction needs data that isn’t ready yet

**Control Hazards** happen with branching (if statements, loops). The CPU doesn't know which instruction to fetch next until the branch condition is evaluated. Modern CPUs use branch prediction to guess the outcome and speculatively execute instructions, but mispredictions cost a lof of cycles (pipeline flushes).

CPU tries to prevent these hazards:
1. Keep frequently used data close to the CPU in caches (L1, L2, L3)
2. Guess what data it might need next and fetch it early (Prefetching)
3. Rearrange instructions to avoid stalls (out-of-order execution)

**Example of pipeline stall**
```
Cycle 1: [Fetch ADD] 
Cycle 2: [Decode ADD] [Fetch LDR]
Cycle 3: [Execute ADD] [Decode LDR] [Fetch ADD]
Cycle 4: [Store ADD] [Execute LDR] [Decode ADD] ← Memory request starts
Cycle 5:             [STALL] [STALL] ← Waiting for memory
Cycle 6:             [STALL] [STALL] ← Still waiting...
...
Cycle 304:           [Store LDR] [Execute ADD] ← Finally got data!
```

# Memory Latency
Let's explore ways to identify memory latencies in Go.

We have two implementations of binary search tree. One implementation is recursive memory-unaware and second is memory-aware by keeping data in contagious block and having predictable access patterns.

Full implementation can be found here: https://github.com/Elvis339/compiler.rs/blob/main/code/memory-and-performance/cmd/memaccess/memaccess.go

**Pointer Base BST**
```go
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
```

**Contiguous BST**
```go
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
```

## Analysis
This benchmark compares two Binary Search Tree implementations a traditional pointer-based approach versus a contiguous array-based approach. The contiguous implementation demonstrates **53% better performance** despite using significantly more memory.

| Implementation | Avg Time/Op | Memory Used | Allocations | Performance Gain |
|---------------|-------------|-------------|-------------|------------------|
| PointerBST    | 279.54 ms   | ~600 KB     | ~25,000     | Baseline         |
| ContiguousBST | 128.91 ms   | ~9 MB       | 2           | **53% faster**   |

The reason for ~2.17x improvement partially comes from not having recursive calls, the burden of function setup is heavy. Additionally as previously described we're experiencing better spatial locality. Let's dive deeper into the specific performance factors:

Modern CPUs fetch memory in **cache lines** (typically 64 bytes). When you request one byte, the CPU automatically fetches the entire cache line:

```go
// Pointer BST - each node scattered in memory
node1 at 0x1000  // Cache line 1
node2 at 0x5000  // Cache line 2 (different, far away)
node3 at 0x9000  // Cache line 3 (different, far away)

// Array BST - multiple values in same cache line
data[0] at 0x2000  // Cache line 1
data[1] at 0x2008  // Cache line 1 (same cache line!)
data[2] at 0x2016  // Cache line 1 (same cache line!)
```

CPUs have automatic hardware prefetchers that detect sequential access patterns:

```go
// Predictable sequential array access pattern
data[0] → data[1] → data[2] → data[3]

// Unpredictable pointer chaising trought the heap access pattern
nodeA → nodeX → nodeM → nodeB
```

# Tools

Unfortunately, Go's pprof tool doesn't provide clear insights into memory access patterns for this type of performance analysis. For deeper memory latency investigation, we need to examine the actual CPU instructions generated.

While `perf` on Linux provides excellent memory profiling capabilities, we'll use assembly analysis to understand the performance differences:

[Assembly script](https://github.com/Elvis339/compiler.rs/blob/main/code/memory-and-performance/cmd/memaccess/demo/assembly.sh).

## Assembly Reveals the Truth

The ARM assembly output shows exactly why the array BST is faster:

**Pointer BST** - Random memory access:
```assembly
ldr x0, [x0, #8]     # Load n.left from random heap address
bl  search_function  # Function call overhead + stack management
```

**Array BST** - Predictable arithmetic:
```assembly
lsl x2, x2, #1       # index = index * 2 (bit shift)
add x2, x2, #1       # index = 2*index + 1 (simple addition)
b   loop_start       # Jump back to loop (no function calls)
```

The key insights that assembly proves is that the memory layour affects the fundamental CPU instructions generated. Better spatial locality translates directly to more efficient machine code.

# Conclusion

Modern CPUs are incredibly fast, but their performance depends on keeping the instruction pipeline fed with data. Modern hardware and compilers have techniques to avoid pipeline stalls, but we as software engineers must understand abstractions below and help with the optimization.


- CPUs access data through multiple cache levels (L1, L2, L3) before reaching main memory. Understanding this hierarchy helps explain why some code runs faster than others.
- Accessing nearby memory locations allows hardware prefetchers to predict patterns and keep data flowing to the CPU. Sequential array access leverages cache lines efficiently, while scattered pointer chasing defeats these optimizations entirely. This is why in specific cases arrays outpreform maps.

**Being memory-aware means:**
- Choosing data structures that promote spatial locality (arrays vs linked lists)
- Understanding how your data layout affects cache behavior  
- Recognizing that algorithm complexity isn't everything - you must include memory access into the mix

**Tools for analysis:** While profilers like Go's pprof show where time is spent, they don't reveal memory access patterns. Assembly analysis helps understand the actual CPU instructions generated. On macOS, Instruments provides memory profiling capabilities, though it's more valuable for higher-level analysis than micro-optimizations.

Our BST comparison proved that identical algorithms with different memory layouts produce fundamentally different performance characteristics.

## What's next?

In the following blog post series, we'll analyze Go's experimental [Green Tea](https://github.com/golang/go/issues/73581) garbage collector, which applies these same spatial locality principles to automatic memory management.