package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
)

var (
	cpuProfile = flag.String("cpuprofile", "", "filename where we should write the cpu profile")
	memProfile = flag.String("memprofile", "", "filename where we should write the mem profile")
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
	println(*memProfile, flag.NFlag())
	if len(*memProfile) != 0 {
		f, err := os.Create(*memProfile)
		if err != nil {
			log.Fatal(err)
		}
		defer pprof.WriteHeapProfile(f)
	}
	b := New()
	b.Perft(5)
}
