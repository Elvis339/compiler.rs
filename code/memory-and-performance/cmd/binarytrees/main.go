// The Computer Language Benchmarks Game
// http://benchmarksgame.alioth.debian.org/
//
// Go implementation of binary-trees, based on the reference implementation
// gcc #3, on Go #8 (which is based on Rust #4)
//
// Comments aim to be analogous as those in the reference implementation and are
// intentionally verbose, to help programmers unexperienced in GO to understand
// the implementation.
//
// The following alternative implementations were considered before submitting
// this code. All of them had worse readability and didn't yield better results
// on my local machine:
//
// 0. general:
// 0.1 using uint32, instead of int;
//
// 1. func Count:
// 1.1 using a local stack, instead of using a recursive implementation; the
//     performance degraded, even using a pre-allocated slice as stack and
//     manually handling its size;
// 1.2 assigning Left and Right to nil after counting nodes; the idea to remove
//     references to instances no longer needed was to make GC easier, but this
//     did not work as intended;
// 1.3 using a walker and channel, sending 1 on each node; although this looked
//     idiomatic to GO, the performance suffered a lot;
// 2. func NewTree:
// 2.1 allocating all tree nodes on a tree slice upfront and making references
//     to those instances, instead of allocating two sub-trees on each call;
//     this did not improve performance;
//
// Contributed by Gerardo Lima
// Reviewed by Diogo Simoes
// Based on previous work from Adam Shaver, Isaac Gouy, Marcel Ibes Jeremy,
//  Zerfas, Jon Harrop, Alex Mizrahi, Bruno Coutinho, ...
//

// https://benchmarksgame-team.pages.debian.net/benchmarksgame/program/binarytrees-go-2.html

package main

import (
	"flag"
	"fmt"
	"github.com/arl/statsviz"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type Tree struct {
	Left  *Tree
	Right *Tree
}

// Count the nodes in the given complete binary tree.
func (t *Tree) Count() int {
	// Only test the Left node (this binary tree is expected to be complete).
	if t.Left == nil {
		return 1
	}
	return 1 + t.Right.Count() + t.Left.Count()
}

// Create a complete binary tree of `depth` and return it as a pointer.
func NewTree(depth int) *Tree {
	if depth > 0 {
		return &Tree{Left: NewTree(depth - 1), Right: NewTree(depth - 1)}
	} else {
		return &Tree{}
	}
}

func Run(maxDepth int) {

	var wg sync.WaitGroup

	// Set minDepth to 4 and maxDepth to the maximum of maxDepth and minDepth +2.
	const minDepth = 4
	if maxDepth < minDepth+2 {
		maxDepth = minDepth + 2
	}

	// Create an indexed string buffer for outputing the result in order.
	outCurr := 0
	outSize := 3 + (maxDepth-minDepth)/2
	outBuff := make([]string, outSize)

	// Create binary tree of depth maxDepth+1, compute its Count and set the
	// first position of the outputBuffer with its statistics.
	wg.Add(1)
	go func() {
		tree := NewTree(maxDepth + 1)
		msg := fmt.Sprintf("stretch tree of depth %d\t check: %d",
			maxDepth+1,
			tree.Count())

		outBuff[0] = msg
		wg.Done()
	}()

	// Create a long-lived binary tree of depth maxDepth. Its statistics will be
	// handled later.
	var longLivedTree *Tree
	wg.Add(1)
	go func() {
		longLivedTree = NewTree(maxDepth)
		wg.Done()
	}()

	// Create a lot of binary trees, of depths ranging from minDepth to maxDepth,
	// compute and tally up all their Count and record the statistics.
	for depth := minDepth; depth <= maxDepth; depth += 2 {
		iterations := 1 << (maxDepth - depth + minDepth)
		outCurr++

		wg.Add(1)
		go func(depth, iterations, index int) {
			acc := 0
			for i := 0; i < iterations; i++ {
				// Create a binary tree of depth and accumulate total counter with its
				// node count.
				a := NewTree(depth)
				acc += a.Count()
			}
			msg := fmt.Sprintf("%d\t trees of depth %d\t check: %d",
				iterations,
				depth,
				acc)

			outBuff[index] = msg
			wg.Done()
		}(depth, iterations, outCurr)
	}

	wg.Wait()

	// Compute the checksum of the long-lived binary tree that we created
	// earlier and store its statistics.
	msg := fmt.Sprintf("long lived tree of depth %d\t check: %d",
		maxDepth,
		longLivedTree.Count())
	outBuff[outSize-1] = msg

	// Print the statistics for all of the various tree depths.
	for _, m := range outBuff {
		fmt.Println(m)
	}
}

func getExecutableName() string {
	executable, err := os.Executable()
	if err != nil {
		return "unknown"
	}
	return filepath.Base(executable)
}

func startStatsvizServer() string {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("Failed to find available port:", err)
	}
	p, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		log.Fatal("Failed to find available port")
	}
	port := p.Port
	_ = listener.Close()

	err = statsviz.Register(http.DefaultServeMux)
	if err != nil {
		log.Fatal("Failed to register statsviz:", err)
	}

	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	}()

	return fmt.Sprintf("http://localhost:%d/debug/statsviz/", port)
}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	execName := getExecutableName()

	// Start statsviz HTTP server
	statsvizURL := startStatsvizServer()
	fmt.Printf("statsviz server: %s\n", statsvizURL)

	// CPU profiling
	cpuFile, err := os.Create(filepath.Join("traces", fmt.Sprintf("%s_cpu.pprof", execName)))
	if err != nil {
		log.Fatal(err)
	}
	defer cpuFile.Close()

	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()

	// Parse command line arguments
	n := 0
	flag.Parse()
	if flag.NArg() > 0 {
		n, _ = strconv.Atoi(flag.Arg(0))
	}

	fmt.Println("n:", n)
	start := time.Now()
	Run(n)
	elapsed := time.Since(start)
	fmt.Println("elapsed:", elapsed)

	// Memory profiling
	runtime.GC()
	memFile, err := os.Create(filepath.Join("traces", fmt.Sprintf("%s_mem.pprof", execName)))
	if err != nil {
		log.Fatal(err)
	}
	defer memFile.Close()

	if err := pprof.WriteHeapProfile(memFile); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Generated: %s_cpu.pprof, %s_mem.pprof\n", execName, execName)
	fmt.Println("Press Ctrl+C to stop the server")
	<-stop
}
