+++
date = '2025-05-26T21:13:14+04:00'
draft = true
title = 'Memory and Performance: Garbage Collection Pt. 2'
tags = ['go', 'rust', 'software-performance', 'profiling', 'memory', 'arm', 'cpu']
categories = ['software', 'software-performance']
series = ['memory-and-performance']
description = "Profiling and tinkering with memory through the lens of Go's experimental Green Tea garbage-collector. Experimenting with tools available to gain performance oriented signals."
+++

For the sake of the article length, I won't walk through running every command here. Instead, I've documented everything to make it seamless for you to reproduce these results yourself. I encourage you to run the profiling tools and discover the patterns firsthand.

**All code, scripts, and detailed instructions are available [here](https://github.com/Elvis339/compiler.rs/blob/main/code/memory-and-performance/README.md#graph)**

# Recap

- In the first blog post, we explored the effects of thoughtful memory layout. The array-based BST was 53% faster than the pointer-based version despite using more memory, because sequential access leverages CPU cache lines and having predictable access patterns optimize hardware prefetchers.
- In the previous blog post, we learned what garbage collection is and discovered that Go's current GC suffers from the exact same problem: it spends 85% of its time chasing scattered pointers through memory, with over 35% of CPU cycles stalled on memory accesses - confirming our findings from the first post about spatial locality.
- Now, we're going to explore how to identify when GC is bottlenecking your algorithms using real profiling tools like pprof and Instruments, and demonstrate how Green Tea's span-based scanning approach could solve these memory access problems.

# Reading GC Traces

We have a [graph](https://github.com) implementation that creates 2M random nodes and performs BFS by visiting all nodes. 
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

While Green Tea solves the spatial locality problem at the GC level, we can still achieve significant improvements by designing your data structures and allocation patterns to be GC-friendly.
