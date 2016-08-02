package PoaGo

import (
	"errors"
	"fmt"
	"math"
)

type DpMatrix struct {
	lX     int
	lY     int
	matrix [][]float64
}

func DpMatrixConstruct(lX, lY int) *DpMatrix {
	matrix := make([][]float64, lY)

	for i := 0; i < lY; i++ {
		inner := make([]float64, lX)
		matrix[i] = inner
	}

	return &DpMatrix{lX: lX, lY: lY, matrix: matrix}
}

func DpMatrixConstructFull(lX, lY int, value float64) *DpMatrix {
	matrix := make([][]float64, lY)

	for i := 0; i < lY; i++ {
		inner := make([]float64, lX)
		for j := 0; j < lX; j++ {
			inner[j] = value
		}
		matrix[i] = inner
	}

	return &DpMatrix{lX: lX, lY: lY, matrix: matrix}
}

func (self *DpMatrix) String() string {
	s := fmt.Sprintf("%v\n", self.matrix[0])
	for i := 1; i < self.lY; i++ {
		s += fmt.Sprintf("%v\n", self.matrix[i])
	}
	s += fmt.Sprintf("%v x %v\n", self.lX, self.lY)
	return s
}

func (self *DpMatrix) checkCoords(x, y int) error {
	if x >= self.lX || y >= self.lY {
		e := fmt.Sprintf("DpMatrix coordinates out of bounds x: %v y: %v", x, y)
		return errors.New(e)
	} else {
		return nil
	}
}

func (self *DpMatrix) SetValue(x, y int, value float64) error {
	ok := self.checkCoords(x, y)
	if ok != nil {
		e := fmt.Sprintf("DpMatrix.SetValue Failed %v", ok)
		//panic(e)
		fmt.Println(e)
		return ok
	}
	row := self.matrix[y]
	row[x] = value
	return ok
}

func (self *DpMatrix) GetValue(x, y int) float64 {
	ok := self.checkCoords(x, y)
	if ok != nil {
		e := fmt.Sprintf("DpMatrix.GetValue failed %v", ok)
		//panic(e)
		fmt.Println(e)
		return math.NaN()
	}
	row := self.matrix[y]
	value := row[x]
	return value
}

func (self *DpMatrix) WhereMax() (int, int, float64) {
	m := math.Inf(-1)
	X := 0
	Y := 0
	for y := 0; y < self.lY; y++ {
		row := self.matrix[y]
		for x := 0; x < self.lX; x++ {
			if row[x] >= m {
				m = row[x]
				X = x
				Y = y
			}
		}
	}
	return X, Y, m
}

type MoveOption struct {
	score        float64
	backGraphIdx int
	backSeqIdx   int
	moveType     string
}

func MoveOptionConstruct(score float64, backGraphIdx, backSeqIdx int, move string) *MoveOption {
	return &MoveOption{score: score, backGraphIdx: backGraphIdx, backSeqIdx: backSeqIdx, moveType: move}
}

func (self MoveOption) GetScore() float64 {
	return self.score
}

func (self MoveOption) String() string {
	return fmt.Sprintf("(%v %v %v %v)", self.score, self.backGraphIdx, self.backSeqIdx, self.moveType)
}

// returns the MoveOption with the max score
func MaxMoveOption(options []*MoveOption) (*MoveOption, error) {
	score := math.Inf(-1) // init to neg infinity
	bestIdx := -1
	for i, opt := range options {
		thisScore := opt.GetScore()
		if thisScore > score {
			score = thisScore
			bestIdx = i
		}
	}

	var ok error

	if bestIdx < 0 {
		ok = errors.New("Didn't find max")
	} else {
		ok = nil
	}

	return options[bestIdx], ok
}

type PairwiseAlignmentParameters struct {
	matchScore     float64
	mismatchScore  float64
	openGapScore   float64
	extendGapScore float64
}

func PairwiseAlignmentParametersConstruct(matchScore, mismatchScore, openGapScore, extendGapScore float64) *PairwiseAlignmentParameters {
	return &PairwiseAlignmentParameters{matchScore: matchScore, mismatchScore: mismatchScore, openGapScore: openGapScore, extendGapScore: extendGapScore}
}

func (self *PairwiseAlignmentParameters) MatchBases(c1, c2 string) float64 {
	if c1 == c2 {
		return self.matchScore
	} else {
		return self.mismatchScore
	}
}

type PairwiseAlignment struct {
	stringIdxs []int
	matches    []int
	sequence   string
	label      string
}

func PairwiseAlignmentConstruct(strIdxs, matches []int, sequence, label string) *PairwiseAlignment {
	return &PairwiseAlignment{stringIdxs: strIdxs, matches: matches, sequence: sequence, label: label}
}

func MakeNodeIndexMaps(g *PoaGraph) (map[int]int, map[int]int) {
	IdToIndex := make(map[int]int)
	IndexToId := make(map[int]int)

	if g.needSort {
		g.TopoSort()
		g.testSort()
	}

	for i, n := range g.nodeList {
		node := g.nodeDict[n]
		IdToIndex[node.id] = i
		IndexToId[i] = node.id
	}

	return IdToIndex, IndexToId
}

func prevIndices(n *Node, IdtoIndex map[int]int) []int {
	prev := make([]int, 0)
	for k := range n.inEdges {
		prev = append(prev, IdtoIndex[k])
	}

	if len(prev) == 0 {
		prev = append(prev, -1)
	}
	return prev
}

func DoTraceBack(scores, backSeqMatrix, backGrphMatrix *DpMatrix, IndexToId map[int]int) ([]int, []int) {
	besti, bestj, _ := scores.WhereMax()
	matches := make([]int, 0)
	strIndexs := make([]int, 0)

	for (scores.GetValue(besti, bestj) > 0) && !(bestj == 0 && besti == 0) {
		nexti := int(backGrphMatrix.GetValue(besti, bestj))
		nextj := int(backSeqMatrix.GetValue(besti, bestj))
		curStrIdx := besti - 1
		curNodeIdx := IndexToId[bestj-1]
		if nextj != besti {
			strIndexs = append(strIndexs, curStrIdx)
		} else {
			strIndexs = append(strIndexs, -1)
		}
		if nexti != bestj {
			matches = append(matches, curNodeIdx)
		} else {
			matches = append(matches, -1)
		}
		bestj = nexti
		besti = nextj
	}
	intArrayReverse(strIndexs)
	intArrayReverse(matches)

	return strIndexs, matches
}

func AlignStringToGraph(g *PoaGraph, aln *PairwiseAlignmentParameters, sequence, label string) *PairwiseAlignment {
	IdToIndex, IndexToId := MakeNodeIndexMaps(g)

	lX := len(sequence)
	lY := g.nbNodes

	scores := DpMatrixConstruct(lX+1, lY+1)
	backGrphMatrix := DpMatrixConstruct(lX+1, lY+1)
	backSeqMatrix := DpMatrixConstruct(lX+1, lY+1)
	insertCostMatrix := DpMatrixConstructFull(lX+1, lY+1, aln.openGapScore)
	deleteCostMatrix := DpMatrixConstructFull(lX+1, lY+1, aln.openGapScore)

	for i, nodeIdx := range g.nodeList {
		node := g.nodeDict[nodeIdx]
		pbase := node.base

		for j, sbase := range sequence {
			candidates := make([]*MoveOption, 0) // could be optimized if I know the pred nodes already
			insScore := scores.GetValue(j, i+1) + insertCostMatrix.GetValue(j, i+1)
			insertOption := MoveOptionConstruct(insScore, i+1, j, "INSERT")
			candidates = append(candidates, insertOption)
			previousIdxs := prevIndices(node, IdToIndex)
			for _, predIdx := range previousIdxs {
				// handle the matches
				matchScore := scores.GetValue(j, predIdx+1) + aln.MatchBases(pbase, string(sbase))
				matchOption := MoveOptionConstruct(matchScore, predIdx+1, j, "MATCH")
				// handle the deletes
				deleteScore := scores.GetValue(j+1, predIdx+1) + deleteCostMatrix.GetValue(j+1, predIdx+1)
				deleteOption := MoveOptionConstruct(deleteScore, predIdx+1, j+1, "DELETE")
				candidates = append(candidates, matchOption, deleteOption)
			}
			maxMove, _ := MaxMoveOption(candidates)
			scores.SetValue(j+1, i+1, maxMove.score)
			backGrphMatrix.SetValue(j+1, i+1, float64(maxMove.backGraphIdx))
			backSeqMatrix.SetValue(j+1, i+1, float64(maxMove.backSeqIdx))
			if maxMove.moveType == "INSERT" {
				insertCostMatrix.SetValue(j+1, i+1, aln.extendGapScore)
			}
			if maxMove.moveType == "DELETE" {
				deleteCostMatrix.SetValue(j+1, i+1, aln.extendGapScore)
			}
			if scores.GetValue(j+1, i+1) < 0 {
				scores.SetValue(j+1, i+1, 0)
				backGrphMatrix.SetValue(j+1, i+1, -1)
				backSeqMatrix.SetValue(j+1, i+1, -1)
			}
		}
	}

	strIdxs, matches := DoTraceBack(scores, backSeqMatrix, backGrphMatrix, IndexToId)

	pA := PairwiseAlignmentConstruct(strIdxs, matches, sequence, label)

	return pA
	//return strIdxs, matches
}
