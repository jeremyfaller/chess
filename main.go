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

	perft = flag.Int("perft", 0, "run perft then quit")
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

	if *perft != 0 {
		b := New()
		res := b.Perft(*perft, PerftConfig{Quiet: true})
		fmt.Println(res)
		return
	}

	u := NewUCI()
	if err := u.Run(); err != nil {
		log.Fatal(err)
	}
}
