+++
date = '2025-05-26T21:13:14+04:00'
draft = false
title = 'Memory and Performance: Garbage Collection Pt. 2'
tags = ['go', 'rust', 'software-performance', 'profiling', 'memory', 'arm', 'cpu']
categories = ['software', 'software-performance']
series = ['memory-and-performance']
description = "Analyzing and tinkering with memory through lens of Go's experimental Green Tea garbage-collector. Experimenting with tools available to gain performance oriented signals."
+++

For the sake of the article length, I won't walk through running every command here. 

Instead, I've documented everything to make it seamless for you to reproduce these results yourself. I encourage you to run the profiling tools and discover the patterns firsthand.

**All code, scripts, and detailed instructions are available [here](https://github.com/Elvis339/compiler.rs/blob/main/code/memory-and-performance/README.md#graph)**

# Recap

- In the first blog post, we explored the effects of thoughtful memory layout. The array-based BST was 53% faster than the pointer-based version despite using more memory, because sequential access leverages CPU cache lines and having predictable access patterns optimize hardware prefetchers.
- In the previous blog post, we learned what garbage collection is and discovered that Go's current GC suffers from the exact same problem: it spends 85% of its time chasing scattered pointers through memory, with over 35% of CPU cycles stalled on memory accesses - confirming our findings from the first post about spatial locality.
- Now, we're going to explore how to identify when GC is bottlenecking your algorithms using real profiling tools like pprof and Instruments, and demonstrate how Green Tea's span-based scanning approach could solve these memory access problems.

# Reading GC Traces

We have a [graph](https://github.com/Elvis339/compiler.rs/blob/fb3031882458b4f035584024a560f9997076a909/code/memory-and-performance/cmd/graph/graph.go) implementation that creates 2M random nodes and performs BFS by visiting all nodes. 
Running `make run EXEC=graph` sets `GODEBUG=gctrace=1`, which hints to any running Go application to start garbage collection tracing.

**Current Go GC trace output:**
```
gc 1 @0.001s 20%:                    # gc 1 - First garbage collection cycle, occurred 0.001 seconds after program start, 20% of total program time spent in GC so far
   Clock: 0.023+9.3+0.17 ms          # 0.023ms - STW (stop-the-world) sweep termination, 9.3ms - concurrent marking phase, 0.17ms - STW mark termination
   CPU: 0.19+0.22/17/0.26+1.4 ms     # 0.19ms - STW sweep CPU time, 0.22ms - mutator assist, 17ms - background workers, 0.26ms - dedicated workers, 1.4ms - STW mark CPU time
   Heap: 15->20->19 MB, 15 MB goal   # 15MB - heap size at GC start, 20MB - peak heap during GC, 19MB - live heap after GC, 15MB - target for next GC
    
gc 2 @0.024s 14%:
    Clock: 0.008+9.3+0.016 ms
    CPU: 0.069+2.1/18/1.6+0.13 ms
    Heap: 35->41->41 MB, 39 MB goal
    
gc 3 @0.055s 10%:
    Clock: 0.012+8.2+0.027 ms
    CPU: 0.10+0.049/16/34+0.22 ms
    Heap: 74->83->82 MB, 82 MB goal
    
gc 4 @0.179s 6%:
    Clock: 0.014+26+0.014 ms
    CPU: 0.11+0.037/53/130+0.11 ms
    Heap: 140->144->127 MB, 165 MB goal
    
gc 5 @0.415s 7%:
    Clock: 0.014+107+0.011 ms
    CPU: 0.11+1.0/213/497+0.094 ms
    Heap: 218->228->184 MB, 256 MB goal
```

**Green Tea trace output**


Running `make run EXEC=graphx` also sets `GODEBUG=gctrace=1` but runs the binary with [Green Tea](https://github.com/golang/go/issues/73581) garbage collector.

**Trace output:**
```
gc 1 @0.001s 6%:                     # gc 1 - First garbage collection cycle, occurred 0.001 seconds after program start, 6% of total program time spent in GC so far
   Clock: 0.015+0.29+0.074 ms        # 0.015ms - STW sweep termination, 0.29ms - concurrent marking phase, 0.074ms - STW mark termination
   CPU: 0.12+0/0.29/0+0.59 ms        # 0.12ms - STW sweep CPU time, 0ms - mutator assist, 0.29ms - background workers, 0ms - dedicated workers, 0.59ms - STW mark CPU time
   Heap: 15->15->15 MB, 15 MB goal   # 15MB - heap size at GC start, 15MB - peak heap during GC, 15MB - live heap after GC, 15MB - target for next GC

gc 2 @0.015s 5%:
    Clock: 0.009+3.7+0.043 ms
    CPU: 0.075+0.048/7.2/9.3+0.35 ms
    Heap: 30->33->32 MB, 31 MB goal
    
gc 3 @0.035s 7%:
    Clock: 0.015+7.9+0.15 ms
    CPU: 0.12+0.66/14/21+1.2 ms
    Heap: 58->65->64 MB, 65 MB goal
    
gc 4 @0.078s 5%:
    Clock: 0.019+7.5+0.031 ms
    CPU: 0.15+0.045/14/35+0.25 ms
    Heap: 112->113->110 MB, 128 MB goal

gc 5 @0.352s 3%:
    Clock: 0.020+40+0.029 ms
    CPU: 0.16+1.3/80/181+0.23 ms
    Heap: 194->197->156 MB, 220 MB goal
```

The same implementation produces different performance characteristics simply by changing the garbage collection algorithm.

- Current Go's GC starts high with 20% of total program time spent in GC before settling down to 4-7% range
- Green Tea starts lower and stays consistently low 2-7% range

## Peaks

| Metric | Current GC Peak | Green Tea Peak | Difference |
|--------|-----------------|----------------|------------|
| **Worst Marking Time** | 107ms (GC 5) | 40ms (GC 5) | **63% lower** |
| **Highest GC %** | 20% (GC 1) | 7% (GC 3) | **65% lower** |
| **Peak CPU Overhead** | 497ms (GC 5) | 181ms (GC 5) | **64% reduction** |
| **Max Heap Growth** | 228MB (15→228) | 197MB (15→197) | **14% more efficient** |

## What If You Can't Change the GC Algorithm?

We can still achieve significant improvements by designing your data structures and allocation patterns to be GC-friendly.

We have a [compact graph](https://github.com/Elvis339/compiler.rs/blob/fb3031882458b4f035584024a560f9997076a909/code/memory-and-performance/cmd/graph/compact.go) implementation which is using [sync.Pool](https://pkg.go.dev/sync#Pool) for reusing temporary allocations (queue and visited map). 
This improves GC performance because it eliminates repeated allocations in the hot path - instead of creating new slices and maps on every BFS call, we reuse pre-allocated objects from the pool, dramatically reducing allocation pressure and giving the garbage collector fewer objects to track and clean up.

**Compact GC Trace**
```
gc 1 @0.000s 23%: 0.007+8.8+0.014 ms clock, 0.061+0.076/17/16+0.11 ms cpu, 62->68->68 MB, 62 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 2 @0.030s 8%: 0.025+3.6+0.010 ms clock, 0.20+0.11/6.5/13+0.084 ms cpu, 124->130->130 MB, 137 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 3 @0.216s 2%: 0.032+10+0.009 ms clock, 0.26+0/20/0.10+0.077 ms cpu, 246->253->194 MB, 261 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 4 @0.601s 1%: 0.034+9.5+0.025 ms clock, 0.27+0/14/5.1+0.20 ms cpu, 375->375->245 MB, 390 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 5 @1.820s 0%: 0.064+0.36+0.006 ms clock, 0.51+0/0.25/0.33+0.053 ms cpu, 371->371->64 MB, 490 MB goal, 0 MB stacks, 0 MB globals, 8 P (forced)
```

**Pointer Chasing Implementation**
```
gc 1 @0.000s 18%: 0.027+2.8+0.011 ms clock, 0.22+0.12/5.1/4.7+0.094 ms cpu, 16->19->18 MB, 17 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 2 @0.012s 9%: 0.033+3.8+0.011 ms clock, 0.26+0.047/7.1/13+0.092 ms cpu, 33->36->36 MB, 37 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 3 @0.031s 8%: 0.036+6.0+0.026 ms clock, 0.29+0.25/11/27+0.21 ms cpu, 63->70->69 MB, 73 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 4 @0.092s 6%: 0.049+15+0.036 ms clock, 0.39+0.081/31/76+0.28 ms cpu, 119->121->115 MB, 140 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 5 @0.330s 7%: 0.031+89+0.024 ms clock, 0.25+1.1/177/428+0.19 ms cpu, 196->205->166 MB, 230 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 6 @0.802s 7%: 0.063+155+0.025 ms clock, 0.50+0/311/798+0.20 ms cpu, 283->318->241 MB, 332 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 7 @1.426s 6%: 0.037+156+0.019 ms clock, 0.29+0/312/780+0.15 ms cpu, 436->436->256 MB, 483 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 8 @2.780s 3%: 0.041+0.21+0.004 ms clock, 0.33+0/0.19/0.040+0.037 ms cpu, 353->353->0 MB, 513 MB goal, 0 MB stacks, 0 MB globals, 8 P (forced)
```

Looking at the GC traces, the compact version performs 3 fewer GC cycles (5 vs 8), but the traces seem similar.

Benchmarks should be your first resort for performance indicators. These [benchmarks](https://github.com/Elvis339/compiler.rs/blob/fb3031882458b4f035584024a560f9997076a909/code/memory-and-performance/cmd/graph/bench_test.go) show about **477,828,000 ns/op difference** on my machine in favor of the compact implementation, which translates to roughly **478ms faster per operation**. 

### Pprof Analysis
If we run `go tool pprof traces/graph_compact_cpu.pprof` in one terminal and `go tool pprof traces/graph_ptr-chasing_cpu.pprof` in another, we're observing different CPU profiles.

**Left Terminal (Compact Implementation):**
```
(pprof) top10
Showing nodes accounting for 1340ms, 90.07% of 1510ms total
Showing top 10 nodes out of 67
flat  flat%   sum%        cum   cum%
370ms 25.83% 25.83%      480ms 31.79%  runtime.mapaccess1_fast64
290ms 19.21% 45.03%      290ms 19.21%  runtime.(remap).dverflow (inline)
240ms 15.89% 60.93%     1140ms 75.50%  main.(*compactGraph).bfs
150ms  9.93% 70.86%      150ms  9.93%  runtime.madvise
80ms  5.30% 76.16%      120ms  7.95%  runtime.scanobject
50ms  3.31% 79.47%      370ms 25.83%  runtime.mapaccess_fast64
50ms  3.31% 82.78%       50ms  3.31%  runtime.memclrNoHeapPointers
40ms  2.65% 85.43%       40ms  2.65%  runtime.heapBitsSetType
60ms  2.65% 88.08%       60ms  2.65%  runtime.pthread-cond-signal
30ms  1.99% 90.07%       30ms  1.99%  main.createCompactGraph
```

**Right Terminal (Pointer-Chasing Implementation):**
```
(pprof) top10
Showing nodes accounting for 3720ms, 79.49% of 4680ms total
Dropped 57 nodes (cum <= 23.40ms)
Showing top 10 nodes out of 83
flat  flat%   sum%        cum   cum%
940ms 20.09% 20.09%      940ms 20.09%  runtime.madvise
570ms 12.18% 32.26%      760ms 14.96%  runtime.mapaccess1_fast64
530ms 11.32% 43.59%     1350ms 32.85%  main.bfs
530ms 11.32% 54.91%     1980ms 40.60%  runtime.scanobject
390ms  8.33% 63.25%      390ms 17.50%  runtime.greyobject
250ms  5.34% 68.59%      250ms  5.34%  runtime.(▸span).heapBitsSmallForAddr
160ms  3.42% 72.01%      160ms  3.42%  runtime.(▸htstack).pop
130ms  2.78% 74.79%      130ms  2.78%  runtime.dverflow (inline)
110ms  2.35% 77.14%      600ms  8.55%  runtime.getptmty
110ms  2.35% 79.49%      110ms  2.35%  runtime.memclrNoHeapPointers
```
**Takeaway**

The pointer-chasing implementation shows performance degradation compared to compact implementation:
- **20.09%** of total execution spent on [madvise](https://man7.org/linux/man-pages/man2/madvise.2.html) syscall, indicating the GC is working hard to manage fragmented memory
- **40.60%** of total execution spent on `runtime.scanobject` (object scanning during GC marking)
- **32.85%** of total execution spent on algorithm (`main.bfs`)

In contrast, the compact implementation shows:
- **9.93%** spent on `runtime.madvise` (50% reduction)
- **5.30%** spent on `runtime.scanobject` (87% reduction)
- **75.50%** spent in the actual algorithm (`main.(*compactGraph).bfs`)

# Next
We started somewhere in the middle, then went up the stack (pun intended) now we're going to go a bit deeper and explore virtual memory, address translation, and memory fragmentation. Understanding these lower-level concepts should complete the picture.
