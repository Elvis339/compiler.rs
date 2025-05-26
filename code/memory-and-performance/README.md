# Memory and Performance Analysis Tools

A collection of Go programs for exploring memory management and performance characteristics, created for a series of blog posts on memory optimization and garbage collection behavior.

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

### Binary Trees
Demonstrates memory allocation patterns and GC behavior with tree data structures.

**Standard Go runtime:**
```bash
make run EXEC=btree
```

**With Green Tea GC (experimental):**
```bash
make run EXEC=btreex
```

### Garbage Collection Analysis
Programs designed to stress-test garbage collection and analyze GC performance.

**Standard Go runtime:**
```bash
make run EXEC=garbage
```

**With Green Tea GC (experimental):**
```bash
make run EXEC=garbagex
```

### Memory Access Patterns
Tools for analyzing memory access performance and cache behavior.

**Run with tracing:**
```bash
make run EXEC=memaccess     # Standard runtime
make run EXEC=memaccessx    # Green Tea GC
```

**Run benchmarks:**
```bash
cd cmd/memaccess && go test -bench=. -benchmem -count=6
```

**Generate assembly analysis:**
```bash
cd cmd/memaccess/demo && ./assembly.sh
```

## Program Arguments

You can pass arguments to any program using the `ARGS` parameter:
```bash
make run EXEC=btree ARGS="25"
```

## Output and Traces

- **GC traces:** Saved to `traces/<executable>.gctrace`
- **Profiling data:** Generated as `*.pprof` files in the project root
- **Assembly output:** Generated in `cmd/memaccess/demo/`

## Requirements

- Go (standard runtime)
- `gotip` (for experimental Green Tea GC builds)
- macOS for code signing (optional, skipped on other platforms)

## Green Tea GC

Some programs include variants built with Go's experimental Green Tea garbage collector (`GOEXPERIMENT=greenteagc`). These variants have an 'x' suffix (e.g., `btreex`, `garbagex`) and can be used to compare performance characteristics between the standard and experimental GC implementations.