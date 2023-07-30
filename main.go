package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"runtime/trace"
)

var (
	cpuProfile   = flag.String("cpuprofile", "", "filename where we should write the cpu profile")
	memProfile   = flag.String("memprofile", "", "filename where we should write the mem profile")
	traceProfile = flag.String("traceprofile", "", "filename where we should write trace output")
	depth        = flag.Int("depth", 6, "perft depth")
)

func main() {
	flag.Parse()
	if len(*cpuProfile) != 0 {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if len(*memProfile) != 0 {
		f, err := os.Create(*memProfile)
		if err != nil {
			log.Fatal(err)
		}
		defer pprof.WriteHeapProfile(f)
	}
	if len(*traceProfile) != 0 {
		f, err := os.Create(*traceProfile)
		if err != nil {
			log.Fatal(err)
		}
		trace.Start(f)
		defer trace.Stop()
	}

	b := New()
	fmt.Printf("perft(%d) = %d\n", *depth, b.Perft(*depth))
}
