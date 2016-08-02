package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/ArtRand/PoaGo/lib"
	"os"
)

func check(ok error, msg string) {
	if ok != nil {
		panic(msg)
	} else {
		return
	}
}

func main() {
	inFile := flag.String("f", "", "file location")

	flag.Parse()

	fH, ok := os.Open(*inFile)
	check(ok, fmt.Sprintf("Error opening file %v", *inFile))
	defer fH.Close()

	fqr := PoaGo.FqReader{Reader: bufio.NewReader(fH)}

	// get the first sequence
	r, done := fqr.Iter()
	// add it to the graph
	g := PoaGo.PoaGraphConstruct()
	_, _ = g.AddBaseSequence(r.Seq, r.Name, true)

	aln := PoaGo.PairwiseAlignmentParametersConstruct(4.0, -2.0, -4.0, -2.0)

	for !done {
		r, done = fqr.Iter()
		pA := PoaGo.AlignStringToGraph(g, aln, r.Seq, r.Name)
		g.AddSequenceAlignment(pA)
	}

	seqNames, alnStrings := g.GenerateAlignmentStrings()

	for i := 0; i < len(seqNames); i++ {
		fmt.Printf("%-12s\t%-6s\n", seqNames[i], alnStrings[i])
	}

	return
}
