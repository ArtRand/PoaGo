package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	PoaGo "github.com/ArtRand/PoaGo/lib"
)

func check(ok error, msg string) {
	if ok != nil {
		panic(msg)
	} else {
		return
	}
}

const (
	VERSION  = "NOTSET"
	REVISION = "NOTSET"
)

func main() {
	inFile := flag.String("f", "", "file location")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	version := flag.Bool("version", false, "print the version and revision")

	flag.Parse()

	if *version {
		fmt.Printf("Version: %s", VERSION)
		fmt.Printf("Revision: %s", REVISION)
		os.Exit(0)
	}

	fH, ok := os.Open(*inFile)
	check(ok, fmt.Sprintf("Error opening file %v", *inFile))
	defer fH.Close()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	fqr := PoaGo.FqReader{Reader: bufio.NewReader(fH)}

	// get the first sequence
	r, done := fqr.Iter()
	// add it to the graph
	g := PoaGo.PoaGraphConstruct()
	_, _ = g.AddBaseSequence(r.Seq, r.Name, true)

	aln := PoaGo.PairwiseAlignmentParametersConstruct(4.0, -2.0, -4.0, -2.0)

	for {
		r, done = fqr.Iter()
		if done {
			break
		}
		pA := PoaGo.AlignStringToGraph(g, aln, r.Seq, r.Name)
		g.AddSequenceAlignment(pA)
	}

	seqNames, alnStrings := g.GenerateAlignmentStrings()

	for i := 0; i < len(seqNames); i++ {
		fmt.Printf("%-12s\t%-6s\n", seqNames[i], alnStrings[i])
	}

	return
}
