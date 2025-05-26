package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"time"
)

func getExecutableName() string {
	executable, err := os.Executable()
	if err != nil {
		return "unknown"
	}
	return filepath.Base(executable)
}

func startProfiling(enable bool, execName string) func() {
	if !enable {
		return func() {}
	}

	cpuFile, err := os.Create(filepath.Join("traces", fmt.Sprintf("%s_cpu.pprof", execName)))
	if err != nil {
		log.Fatal("Failed to create CPU profile file:", err)
	}

	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		cpuFile.Close()
		log.Fatal("Failed to start CPU profiling:", err)
	}

	return func() {
		pprof.StopCPUProfile()
		cpuFile.Close()

		runtime.GC()
		memFile, err := os.Create(filepath.Join("traces", fmt.Sprintf("%s_mem.pprof", execName)))
		if err != nil {
			log.Printf("Failed to create memory profile file: %v", err)
			return
		}
		defer memFile.Close()

		if err := pprof.WriteHeapProfile(memFile); err != nil {
			log.Printf("Failed to write heap profile: %v", err)
			return
		}
	}
}

// make run EXEC=graph ARGS="-s 100000 -p"
// make run EXEC=graph ARGS="-v compact -s 100000 -p"
func main() {
	version := flag.String("v", "", "Add compact flag if you want to run optimized version")
	size := flag.Int("s", 1_000_000, "Number of nodes in the graph")
	enableProfiling := flag.Bool("p", false, "Enable CPU and memory profiling")
	flag.Parse()

	execName := getExecutableName()

	v := *version
	if len(v) == 0 {
		v = "ptr-chasing"
	}

	fmt.Printf("Configuration:\n")
	fmt.Printf("  Implementation: %s\n", v)
	fmt.Printf("  Graph size: %d nodes\n", *size)
	fmt.Printf("  Profiling: %t\n", *enableProfiling)
	fmt.Printf("\n")

	stopProfiling := startProfiling(*enableProfiling, fmt.Sprintf("%s_%s", execName, v))
	defer stopProfiling()

	start := time.Now()
	switch *version {
	case "compact":
		graph := createCompactGraph(*size)
		runtime.KeepAlive(graph.bfs(0))
	default:
		graph := createGraph(*size)
		runtime.KeepAlive(bfs(graph))
	}

	duration := time.Since(start)
	fmt.Println("Execution time", duration)
}
