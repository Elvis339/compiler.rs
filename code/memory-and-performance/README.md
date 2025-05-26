# Memory and Performance Analysis Tools

A collection of Go programs for exploring memory management and performance characteristics, created for a series of blog posts on memory optimization and garbage collection behavior. These tools demonstrate the impact of memory layout on algorithm performance and provide practical examples for identifying GC bottlenecks using real profiling tools.

## Quick Start

### Build All Targets
```bash
make all
```
This compiles all executables in the `cmd/` directory. On macOS, executables are automatically code-signed for use with Instruments profiling.

### Clean Up
```bash
make clean
```
Removes all compiled executables, GC traces, and profiling data.

## Available Programs

### Binary Trees (`btree`)
Demonstrates memory allocation patterns and GC behavior with tree data structures.

**Standard Go runtime:**
```bash
make run EXEC=btree
```

**With Green Tea GC (experimental):**
```bash
make run EXEC=btreex
```


### Graph Traversal (`graph`)
Implements breadth-first search on randomly connected graphs to demonstrate scattered memory access patterns and their impact on garbage collection performance.

**Standard Go runtime:**
```bash
make run EXEC=graph ARGS="-v compact -s 2000000 -p"
```

**With Green Tea GC (experimental):**
```bash
make run EXEC=graphx ARGS="-v compact -s 2000000 -p"
```

**Available flags:**
- `-v`: Algorithm version (`compact`, `ptr-chasing` (default))
- `-s`: Number of nodes in the graph (default: 1_000_000)
- `-p`: Enable CPU and memory profiling

### Memory Access Patterns (`memaccess`)
Tools for analyzing memory access performance, cache behavior, and the relationship between data structure layout and performance.

**Run with tracing:**
```bash
make run EXEC=memaccess ARGS="-v ptr -s 10000000 -p"    # Standard runtime
make run EXEC=memaccessx ARGS="-v ptr -s 10000000 -p"   # Green Tea GC
```

**Run benchmarks:**
```bash
cd cmd/memaccess && go test -bench=. -benchmem -count=6
```

**Generate assembly analysis:**
```bash
cd cmd/memaccess/demo && ./assembly.sh
```

**Available flags:**
- `-v`: Algorithm version (`ptr`, `array` (default))
- `-s`: Tree size (default: 5_000_000)
- `-p`: Enable CPU and memory profiling

## Profiling and Analysis

### CPU and Memory Profiling
Running programs with `make run EXEC=<binary>` they generate:
- `*.pprof` files for CPU and memory profiling
- Real-time visualization via statsviz (check console output for URL)

### Outputs and traces

- **GC traces:** Saved to `traces/<executable>.gctrace`
- **Profiling data:** Generated as `<executable>_cpu.pprof` and `<executable>_mem.pprof`
- **Assembly output:** Generated in `cmd/memaccess/demo/`

## Requirements

- **Go 1.23+** (standard runtime)
- **gotip** (for experimental Green Tea GC builds)

## Green Tea GC

Some programs include variants built with Go's experimental Green Tea garbage collector (`GOEXPERIMENT=greenteagc`). These variants have an 'x' suffix (e.g., `btreex`, `graphx`) and can be used to compare performance characteristics between the standard and experimental GC implementations.


## Related Blog Posts

- **Memory and Performance: Latency** - Introduction to memory hierarchy and spatial locality
- **Memory and Performance: Garbage Collection Pt. 1** - Understanding Go's GC and Green Tea
- **Memory and Performance: Garbage Collection Pt. 2** - Profiling and optimization techniques
