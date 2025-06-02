package main

import (
	"flag"
	"fmt"
	"syscall"
	"time"
)

const mb = 1024 * 1024

func getPageFaults() int {
	var rusage syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &rusage); err != nil {
		panic(err)
	}
	return int(rusage.Minflt)
}

func pageFault(data []byte) {
	for i := 0; i < len(data); i += 4096 {
		data[i] = 1
	}
}

func noPageFault(data []byte) {
	for i := 0; i < len(data); i += 4096 {
		data[i] = 2
	}
}

// go:inline
func log(pageFaults int, startTime time.Time) {
	since := time.Since(startTime)
	fmt.Println("PageFaults:", pageFaults)
	fmt.Println("Time elapsed:", since)
}

// make run EXEC=vmem ARGS="-v [fault or call without ARGS]"
func main() {
	version := flag.String("v", "all", "Version: fault")
	flag.Parse()

	data := make([]byte, mb)

	switch *version {
	case "fault":
		fmt.Println("Usage: With Page Fault")

		f0 := getPageFaults()
		start := time.Now()
		pageFault(data)
		log(f0, start)
		return
	default:
		f0 := getPageFaults()
		fmt.Println("Usage: No Page Fault")
		pageFault(data)

		f1 := getPageFaults()
		start := time.Now()
		noPageFault(data)
		log(f1-f0, start)
	}
}
