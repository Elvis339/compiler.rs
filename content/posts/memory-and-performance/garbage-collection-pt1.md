+++
date = '2025-05-25T10:36:27+04:00'
draft = false
title = 'Memory and Performance: Garbage Collection Pt. 1'
tags = ['go', 'rust', 'software-performance', 'profiling', 'memory', 'arm', 'cpu']
categories = ['software', 'software-performance']
series = ['memory-and-performance']
description = "Profiling and tinkering with memory trough the lens of Go's experimental Green Tea garbage-collector. Experimenting with tools available to gain performance oriented signals."
+++

In the previous blog post, we learned the importance of memory latency and how to reduce its impact. Here we're going to apply some of our learning and analyze parallel memory aware experimental Go's garbage collector - [green tea](https://github.com/golang/go/issues/73581)

# What is garbage-collector?
Garbage collector's main idea is to free  programmer from manually allocating and freeing objects. Currently Go programming language [uses](https://github.com/ardanlabs/gotraining/tree/master/reading#garbage-collection) concurrent tri-color mark-and-sweep collector ([animation](https://spin.atomicobject.com/2014/09/03/visualizing-garbage-collection-algorithms)), which we're going to explore and compare against experimental [Green Tea](https://github.com/golang/go/issues/73581) garbage collector.

We're mainly going to focus on Go's implementation of mark-and-sweep algorithm and for that we need to sort out termionlogy.
### Terminology
- **Roots** are objects accesible by a program (things in stack, global variables, etc.)
- **Live** objects that are reacahble from the roots by following pointers
- **Dead** objects that are inaccessible and can be recycled

Go is currently using mark-and-sweep garbage collection algorithm which is a Graph based algorithm. The core idea behind general mark-and-sweep garbage collection is that the objects and pointers form a directed graph.

The work is divided in two stages as the name suggests: marking and sweeping.
- The mark stage starts from the root and using bread-first search marks all of the live objects.
- The sweep stage scans memory and frees unmarked objects because they're not accessible anymore from the program i.e. dead.

> Disclaimer: this is an oversimplification of the algorithm. If interested, read more [here](https://tip.golang.org/doc/gc-guide).

# High level overview of Go's garbage colector
As previously mentioned the collector runs through three phases of work.
1. Mark Setup - STW (Stop the World)
2. Marking - Concurrent
3. Mark Termination - STW

## Mark Setup - Stop the World
Mark setup must stop the world, and that means stopping every goroutine application is running.

Let’s say we have 4 application goroutines running before the start of a collection. Each of those 4 goroutines must be stopped. The only way to do that is for the collector to watch and wait for each goroutine to make a function call. 
Function calls guarantee the goroutines are at a safe point to be stopped. What happens if one of those goroutines doesn’t make a function call but the others do?

The Go team solved [this](https://github.com/golang/go/issues/10958) by making tight loops preemptible. I highly recommend reading this discussion.

## Concurrent Marking
Now the collector can start the Marking phase. The Marking phase consists of marking values in heap memory that are still in-use. 
This work starts by inspecting the stacks for all existing goroutines to find root pointers to heap memory. 
Then the collector must traverse the heap memory graph from those root pointers. It's worth noting that for marking process collector allocates one goroutine.

## Mark Termination - Stop the World
Mark Termination runs after marking various cleanup tasks are performed, and the next collection goal is calculated. Once the collection is finished, every goroutine can be used by the application.

# Green Tea Analysis

## What is Green Tea?
Green Tea is an experimental parallel marking garbage-collection algorithm that is memory-aware in that it endeavours to process objects close to another.

## Why?
Go's current garbage collector is not considering the memory location of the objects that are being processed, exhibiting extremely poor spatial locality by jumping between completely different parts of memory being oblivious to topology.   
The same memory access anti-pattern that made our [pointer-based BST 53% slower than its array-based counterpart](https://blog.compiler.rs/posts/memory-and-performance/latency/#memory-latency).

As a result, on average 85% of the garbage collector's time is spent in the core loop of this graph flood—the scan loop—and >35% of CPU cycles in the scan loop are spent solely [stalled on memory accesses](https://blog.compiler.rs/posts/memory-and-performance/latency/#pipeline-hazards-and-memory-interaction).

## How?
Green Tea seems to be changing the [marking phase](https://blog.compiler.rs/posts/memory-and-performance/garbage-collection/#concurrent-marking).  
Instead of the current approach where concurrent markers scan individual objects the new algorithm scans spans.

```
Go GC: obj1 → obj2 → obj3 (scattered)
Green Tea:  [span: obj1, obj2, obj3, obj4] (contiguous)
```  

### What is span?
> A span is always some multiple of 8 KiB, always aligned to 8 KiB, and consists entirely of objects of one size.  

**Why 8KiB size and alignment?**
Aligning to page boundaries enables efficient virtual memory management and reduces [TLB misses](https://en.wikipedia.org/wiki/Translation_lookaside_buffer). (a topic we're going to explore in the following blog post).

**Why multiple of 8 KiB?**
It seems like this decision is made to keep address arithmetic simple. Like in the previous [blog post](https://blog.compiler.rs/posts/memory-and-performance/latency/#assembly-reveals-the-truth).

Additionally, this approach endorses cache-friendly access patterns when scanning and predictable memory layouts.

# Next
In the next blog post, we'll compare Go's current garbage collector with the experimental Green Tea algorithm—not to prove one superior, but to explore how to identify performance signals and act on them.