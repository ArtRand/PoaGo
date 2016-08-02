package PoaGo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"os"
	"bufio"
)

func TestNodeConstruct(t *testing.T) {
	n := NodeConstruct(0, "A")
	if n.base != "A" {
		t.Errorf("Got wrong base, got %v, should be A", n.base)

	}
	if n.id != 0 {
		t.Errorf("Got wrong node id, got %v should be 0", n.id)
	}
	assert.Equal(t, n.id, 0, "Got wrong node ID")
}

func TestNode_String(t *testing.T) {
	n := NodeConstruct(0, "A")
	s := fmt.Sprint(n)
	if s != "(0 : A)" {
		t.Errorf("Printing error, got %v", s)
	}
}

func TestNode_AddInEdge(t *testing.T) {
	n := NodeConstruct(0, "A")
	n.AddInEdge(2, "seq")
	inDeg := n.InDegree()
	if inDeg != 1 {
		t.Errorf("Indegree error got %v should be 1", inDeg)
	}
}

func TestNode_AddOutEdge(t *testing.T) {
	n := NodeConstruct(0, "A")
	n.AddOutEdge(2, "seq")
	inDeg := n.OutDegree()
	if inDeg != 1 {
		t.Errorf("Indegree error got %v should be 1", inDeg)
	}
}

func TestEdgeConstruct(t *testing.T) {
	e := EdgeConstruct(0, 1)
	if e.inNodeID != 0 {
		t.Errorf("In node id error, got %v, should be 0", e.inNodeID)
	}
	if e.outNodeID != 1 {
		t.Errorf("Out node id error, got %v, should be 1", e.outNodeID)
	}
	if len(e.labels) != 0 {
		t.Error("Initialized edge with non-empty set")
	}
}

func TestEdge_AddLabel(t *testing.T) {
	e := EdgeConstruct(0, 1)
	e.AddLabel("seq")
	if !checkForLabel(e.labels, "seq") {
		t.Error("AddLabel error, did not contain addition")
	}
}

func TestPoaGraphConstruct(t *testing.T) {
	g := PoaGraphConstruct()
	if g.nextNodeId != 0 || g.nbNodes != 0 || g.needSort != false {
		t.Error("PoaGraphConstruct error")
	}
}

func TestPoaGraph_AddNode(t *testing.T) {
	g := PoaGraphConstruct()
	nid := g.AddNode("A")
	if g.nextNodeId != 1 || g.nbNodes != 1 || g.needSort != true {
		t.Error("AddNode error")
	}
	if nid != 0 {
		t.Errorf("Node id error, got %v should be 0", nid)
	}
	hasNode := checkForNode(g, 0)
	if !hasNode {
		t.Error("checkForNode Error, should have found node")
	}
	hasNode = checkForNode(g, 1)
	if hasNode {
		t.Error("checkForNode Error, shoud NOT have found node")
	}
}

func TestPoaGraph_AddEdge(t *testing.T) {
	g := PoaGraphConstruct()
	nid1 := g.AddNode("A")
	nid2 := g.AddNode("C")
	g.AddEdge(nid1, nid2, "seq")

	if g.nbEdges != 1 {
		t.Error("number of edges error")
	}
	if g.needSort != true {
		t.Error("Needs sort error")
	}
	if len(g.nodeDict[nid1].outEdges) != 1 {
		t.Errorf("node 0 out edge error, shold be length 1, got %v", len(g.nodeDict[nid1].inEdges))
	}
	if len(g.nodeDict[nid1].inEdges) != 0 {
		t.Errorf("node 0 in edge error, shold be length 0, got %v", len(g.nodeDict[nid1].inEdges))
	}

	if len(g.nodeDict[nid2].outEdges) != 0 {
		t.Errorf("node 1 out edge error, shold be length 0, got %v", len(g.nodeDict[nid1].inEdges))
	}
	if len(g.nodeDict[nid2].inEdges) != 1 {
		t.Errorf("node 1 in edge error, shold be length 1, got %v", len(g.nodeDict[nid1].inEdges))
	}
	g.AddEdge(nid1, nid2, "seq2")

	assert.True(t, len(g.nodeDict[nid1].outEdges[1].labels) == 2, "Didn't add edge label")
}

func TestPoaGraph_AddBaseSequence(t *testing.T) {
	g := PoaGraphConstruct()
	seq := "ACGTACG"
	label := "seq"
	f, l := g.AddBaseSequence(seq, label, true)
	if f != 0 {
		t.Errorf("incorrect start id, should be 0, got %v", f)
	}
	if l != 6 {
		t.Errorf("incorrect last id, should be 6, got %v", l)
	}
}

func TestPoaGraph_TopoSort(t *testing.T) {
	g := PoaGraphConstruct()
	seq := "ACGTACG"
	label := "seq"
	_, _ = g.AddBaseSequence(seq, label, true)
	g.TopoSort()
	c := g.testSort()
	if !c {
		t.Error("Topological sort error on base sequence only")
	}

	//   0 1 2 3 4 5 6
	//   A C G T A C G
	//       \ /
	//        N
	//        7
	next := g.AddNode("N")
	g.AddEdge(2, next, "addition")
	g.AddEdge(next, 4, "addition")
	g.TopoSort()
	c = g.testSort()
	if !c {
		t.Error("Topological sort error with branch")
	}
}

func TestDpMatrixConstruct(t *testing.T) {
	m := DpMatrixConstruct(4, 3)
	if m.lY != 3 {
		t.Errorf("DpMatrix construct error, lY should be 3, got %v", m.lY)
	}
	if m.lX != 4 {
		t.Errorf("DpMatrix construct error, lX should be 4, got %v", m.lX)
	}

	for x := 0; x < m.lX; x++ {
		for y := 0; y < m.lY; y++ {
			v := m.GetValue(x, y)
			if v != 0 {
				t.Error("DpMatrix value not initialized to zero")
			}
		}
	}
}

func TestDpMatrixConstructFull(t *testing.T) {
	val := 33.3
	m := DpMatrixConstructFull(5, 4, val)
	for x := 0; x < m.lX; x++ {
		for y := 0; y < m.lY; y++ {
			v := m.GetValue(x, y)
			if v != val {
				t.Error("DpMatrix value not initialized to zero")
			}
		}
	}
}

func TestDpMatrix_GetAndSetValue(t *testing.T) {
	m := DpMatrixConstruct(3, 3)
	ok := m.SetValue(1, 1, 99.0)
	if ok != nil {
		t.Error("DpMatrix.SetValue error")
	}

	v := m.GetValue(1, 1)
	if v != 99.0 {
		t.Error("DpMatrix.GetValue error")
	}
}

func TestDpMatrix_WhereMax(t *testing.T) {
	m := DpMatrixConstruct(10, 8)
	_ = m.SetValue(0, 0, 99)
	x, y, mx := m.WhereMax()
	if x != 0 || y != 0 || mx != 99 {
		t.Errorf("DpMatrix WhereMax error, should be 0, 0, 99 got %v %v %v", x, y, mx)
	}
	_ = m.SetValue(m.lX-1, m.lY-1, 100)
	x, y, mx = m.WhereMax()
	if x != m.lX-1 || y != m.lY-1 || mx != 100 {
		t.Errorf("DpMatrix WhereMax error, should be 10, 8, 100 got %v %v %v", x, y, mx)
	}
	_ = m.SetValue(5, 3, 101)
	x, y, mx = m.WhereMax()
	if x != 5 || y != 3 || mx != 101 {
		t.Errorf("DpMatrix WhereMax error, should be 5, 3, 101 got %v %v %v", x, y, mx)
	}
}

func TestMoveOptionConstruct(t *testing.T) {
	o := MoveOptionConstruct(0.0, 1, 0, "MATCH")
	if o.score != 0.0 {
		t.Error("MoveOption score error")
	}
	if o.backGraphIdx != 1 {
		t.Error("MoveOption backGraphIdx error")
	}
	if o.backSeqIdx != 0 {
		t.Error("MoveOption backSeqIdx error")
	}
	if o.moveType != "MATCH" {
		t.Error("MoveOption moveType error")
	}
}

func TestMoveOption_GetScore(t *testing.T) {
	m := MoveOptionConstruct(0.0, 1, 0, "MATCH")
	score := m.GetScore()
	if score != 0 {
		t.Error("MoveOption.GetScore error")
	}
}

func TestMaxMoveOption(t *testing.T) {
	m := MoveOptionConstruct(0.0, 1, 0, "MATCH")
	n := MoveOptionConstruct(1.0, 0, 1, "INSERT")
	l := make([]*MoveOption, 0)
	l = append(l, m, n)
	o, ok := MaxMoveOption(l)
	if ok != nil {
		t.Error("MaxMoveOption internal error (ok != nil)")
	}
	if o.score != 1.0 {
		t.Error("MoveOption score error")
	}
	if o.backGraphIdx != 0 {
		t.Error("MoveOption backGraphIdx error")
	}
	if o.backSeqIdx != 1 {
		t.Error("MoveOption backSeqIdx error")
	}
	if o.moveType != "INSERT" {
		t.Error("MoveOption moveType error")
	}
}

func Test_compareEdgeScores(t *testing.T) {
	tup1 := []int{1, 1, 0}
	tup2 := []int{0, 1, 0}
	assert.True(t, compareEdgeScores(tup1, tup2))
	tup2 = []int{2, 1, 0}
	assert.False(t, compareEdgeScores(tup1, tup2))
	assert.False(t, compareEdgeScores(tup1, tup1))
	tup2 = []int{0, 2, 0}
	assert.True(t, compareEdgeScores(tup1, tup2))
	tup2 = []int{1, 2, 0}
	assert.False(t, compareEdgeScores(tup1, tup2))
}

func TestAlignStringToGraph(t *testing.T) {
	g := PoaGraphConstruct()
	seq := "ACGT"
	label := "base"
	_, _ = g.AddBaseSequence(seq, label, true)
	g.TopoSort()
	result := g.testSort()
	if !result {
		t.Error("TestSort failed")
	}

	aln := PairwiseAlignmentParameters{
		matchScore:     4,
		mismatchScore:  -2,
		openGapScore:   -4,
		extendGapScore: -2,
	}

	pA := AlignStringToGraph(g, &aln, "ACT", "new")
	expectedIdxs := []int{0, 1, -1, 2}
	expectedMatches := []int{0, 1, 2, 3}

	assert.True(t, arrEqual(pA.stringIdxs, expectedIdxs))
	assert.True(t, arrEqual(pA.matches, expectedMatches))

	g.AddSequenceAlignment(pA)
	seqNames, alignmentStrings := g.GenerateAlignmentStrings()

	assert.True(t, len(seqNames) == 3)
	assert.True(t, seqNames[0] == "base")
	assert.True(t, seqNames[1] == "new")
	assert.True(t, seqNames[2] == "Consensus0")

	assert.True(t, len(alignmentStrings) == 3)
	assert.True(t, alignmentStrings[0] == "ACGT")
	assert.True(t, alignmentStrings[1] == "AC-T")
	assert.True(t, alignmentStrings[2] == "ACGT")
}

func TestAlignStringToGraph2(t *testing.T) {
	g := PoaGraphConstruct()
	seq := "PKMIVRPQKNETV"
	label := "seq1"
	_, _ = g.AddBaseSequence(seq, label, true)
	g.TopoSort()
	result := g.testSort()
	if !result {
		t.Error("TestSort failed")
	}

	aln := PairwiseAlignmentParameters{
		matchScore:     4,
		mismatchScore:  -2,
		openGapScore:   -4,
		extendGapScore: -2,
	}

	pA := AlignStringToGraph(g, &aln, "THKMLVRNETIM", "seq2")
	g.AddSequenceAlignment(pA)
	_, alignmentStrings := g.GenerateAlignmentStrings()

	assert.True(t, len(alignmentStrings) == 3)
	assert.True(t, alignmentStrings[0] == "--PKMIVRPQKNETV--" || alignmentStrings[0] == "--PKMIVRPQKNET--V")
	assert.True(t, alignmentStrings[1] == "TH-KMLVR---NET-IM" || alignmentStrings[1] == "TH-KMLVR---NETIM-")
	assert.True(t, alignmentStrings[2] == "TH-KMLVRPQKNET-IM" || alignmentStrings[2] == "TH-KMIVRPQKNETIM-" ||
		alignmentStrings[2] == "TH-KMLVRPQKNETIM-" || alignmentStrings[2] == "TH-KMIVRPQKNET-IM")

}

func TestFqReader_Iter(t *testing.T) {
	fH, ok := os.Open("../examples/example1.fa")
	assert.True(t, ok == nil, "Error opening file")
	defer fH.Close()

	var fqr = FqReader{Reader: bufio.NewReader(fH)}
	r, done := fqr.Iter()

	assert.True(t, !done)
	assert.True(t, r.Name == "seq1")
	assert.True(t, r.Seq == "PKMIVRPQKNETV")

	r, done = fqr.Iter()
	assert.True(t, !done)
	assert.True(t, r.Name == "seq2")
	assert.True(t, r.Seq == "THKMLVRNETIM")

	_, done = fqr.Iter()
	assert.True(t, done)

}
