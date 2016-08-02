package PoaGo

import (
	"fmt"
	"math"
	"strings"
)

// constants
const maxFraction float64 = 0.5

// Utils
func intArrayReverse(arr []int) {
	i := len(arr) - 1
	j := 0

	for i > j {
		temp := arr[i]
		arr[i] = arr[j]
		arr[j] = temp
		i -= 1
		j += 1
	}
}

func intArrayMin(arr []int) int {
	var m int = arr[0]
	for _, x := range arr[1:] {
		if x < m {
			m = x
		}
	}
	return m
}

func intArrayMax(arr []int) int {
	m := int(math.Inf(-1))
	for _, x := range arr {
		if m < x {
			m = x
		}
	}
	return m
}

func intArrayArgmax(arr []int) int {
	m := int(math.Inf(-1))
	idx := 0
	for i, x := range arr {
		if m < x {
			m = x
			idx = i
		}
	}
	return idx
}

func arrEqual(arr1, arr2 []int) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	// decide that nil != nil
	if arr1 == nil || arr2 == nil {
		return false
	}

	for i, _ := range arr1 {
		if arr1[i] != arr2[i] {
			return false
		}
	}
	return true
}

// Edge : Directed edge object
type Edge struct {
	inNodeID  int
	outNodeID int
	labels    []string
}

// returns true if the string label is in the slice labels
func checkForLabel(labels []string, label string) bool {
	for _, s := range labels {
		if s == label {
			return true
		}
	}
	return false
}

func (self Edge) String() string {
	edgeString := fmt.Sprintf("(%v) -> (%v)", self.inNodeID, self.outNodeID)
	if len(self.labels) == 0 {
		return edgeString
	}
	return edgeString + fmt.Sprintf("%v", self.labels)
}

func EdgeConstruct(inNodeID, outNodeID int) Edge {
	return Edge{inNodeID: inNodeID, outNodeID: outNodeID, labels: make([]string, 0)}
}

func (self *Edge) AddLabel(label string) {
	// check if we already have this label
	hasLabel := checkForLabel(self.labels, label)
	if !hasLabel {
		self.labels = append(self.labels, label)
	}
}

// Node : Vertex in a DAG
type Node struct {
	id        int           // the id of this node
	base      string        // the base associated with it
	inEdges   map[int]*Edge // key: incident node, value: edge containing label(s)
	outEdges  map[int]*Edge // key: incident node, value: edge containing label(s)
	alignedTo []int
}

func NodeConstruct(id int, base string) *Node {
	return &Node{id: id, base: base, alignedTo: make([]int, 0),
		inEdges: make(map[int]*Edge), outEdges: make(map[int]*Edge)}
}

func (self *Node) addEdge(neighborID int, label string, edgeSet map[int]*Edge) {
	// check if node already in adjacency map, if so, add the label. if not make the edge
	// and init the label
	_, check := edgeSet[neighborID]
	if check {
		edge := edgeSet[neighborID]
		edge.AddLabel(label)
		return
	} else { // if not, make a new edge and label
		edge := EdgeConstruct(neighborID, self.id)
		edge.AddLabel(label)
		edgeSet[neighborID] = &edge
		return
	}
}

func (self Node) String() string {
	return fmt.Sprintf("(%v : %v)", self.id, self.base)
}

func (self *Node) AddInEdge(neighborID int, label string) {
	self.addEdge(neighborID, label, self.inEdges)
}

func (self *Node) AddOutEdge(neighborID int, label string) {
	self.addEdge(neighborID, label, self.outEdges)
}

func (self Node) InDegree() int {
	return len(self.inEdges)
}

func (self Node) OutDegree() int {
	return len(self.outEdges)
}

func (self Node) NextNode(label string) int {
	nextId := -1
	for nid, edge := range self.outEdges {
		hasLabel := checkForLabel(edge.labels, label)
		if hasLabel {
			nextId = nid
		}
	}
	return nextId
}

func collectEdgeLabels(edge Edge, labelSet *[]string) {
	for _, label := range edge.labels {
		contains := checkForLabel(*labelSet, label)
		if !contains {
			*labelSet = append(*labelSet, label)
		}
	}
}

func (self Node) Labels() []string {
	labelSet := make([]string, 0)

	for i := range self.inEdges {
		edge := self.inEdges[i]
		collectEdgeLabels(*edge, &labelSet)
	}
	for i := range self.outEdges {
		edge := self.outEdges[i]
		collectEdgeLabels(*edge, &labelSet)
	}
	return labelSet
}

// POAGraph : A partial order graph for multiple sequence alignment
type PoaGraph struct {
	nextNodeId int
	nbNodes    int
	nbEdges    int
	nodeDict   map[int]*Node
	nodeList   []int
	needSort   bool
	labels     []string
	seqs       []string
	starts     []int
}

func PoaGraphConstruct() *PoaGraph {
	return &PoaGraph{
		nextNodeId: 0,
		nbNodes:    0,
		nbEdges:    0,
		nodeDict:   make(map[int]*Node),
		nodeList:   make([]int, 0),
		needSort:   false,
		labels:     make([]string, 0),
		seqs:       make([]string, 0),
		starts:     make([]int, 0)}
}

func checkForNode(g *PoaGraph, nodeId int) bool {
	_, check := g.nodeDict[nodeId]
	return check
}

func isFinished(finished *[]int, query int) bool {
	for _, x := range *finished {
		if x == query {
			return true
		}
	}
	return false
}

func (self *PoaGraph) AddNode(base string) int {
	// keep track of the idexing of the nodes
	nodeId := self.nextNodeId
	newNode := NodeConstruct(nodeId, base)
	self.nodeDict[nodeId] = newNode
	self.nodeList = append(self.nodeList, nodeId)
	self.nbNodes += 1
	self.nextNodeId += 1
	self.needSort = true
	return nodeId
}

func (self *PoaGraph) AddEdge(startId, endId int, label string) {
	if startId < 0 || endId < 0 {
		return
	}

	if !checkForNode(self, startId) {
		fmt.Errorf("Start node %v not in graph", startId)
	}
	if !checkForNode(self, endId) {
		fmt.Errorf("End node %v not in graph", startId)
	}

	// keep track of the number of edges already going from start->end
	oldNodeEdges := self.nodeDict[startId].OutDegree() + self.nodeDict[endId].InDegree()

	self.nodeDict[startId].AddOutEdge(endId, label)
	self.nodeDict[endId].AddInEdge(startId, label)

	newNodeEdges := self.nodeDict[startId].OutDegree() + self.nodeDict[endId].InDegree()

	if oldNodeEdges != newNodeEdges {
		self.nbEdges += 1
	}

	self.needSort = true
}

func (self *PoaGraph) AddBaseSequence(sequence string, label string, updateSequence bool) (int, int) {
	firstId, lastId := -1, -1
	needSort := self.needSort
	for _, base := range sequence {
		nodeId := self.AddNode(string(base))
		if firstId < 0 {
			firstId = nodeId
		}
		if lastId >= 0 {
			self.AddEdge(lastId, nodeId, label)
		}
		lastId = nodeId
	}
	self.needSort = needSort

	if updateSequence {
		self.seqs = append(self.seqs, sequence)
		self.labels = append(self.labels, label)
		self.starts = append(self.starts, firstId)
	}

	return firstId, lastId
}

func dfs(g *PoaGraph, start int, marked map[int]bool, onStack map[int]bool, finished *[]int) {
	marked[start] = true
	onStack[start] = true
	for neighbor, _ := range g.nodeDict[start].outEdges {
		if !marked[neighbor] {
			dfs(g, neighbor, marked, onStack, finished)
		}
		if onStack[neighbor] {
			fmt.Println("has cycle")
		}
	}
	onStack[start] = false
	*finished = append(*finished, start)
}

func (self *PoaGraph) TopoSort() {
	marked := make(map[int]bool)
	finished := make([]int, 0)
	onStack := make(map[int]bool)

	for _, n := range self.nodeList {
		if !isFinished(&finished, n) {
			dfs(self, n, marked, onStack, &finished)
		}
	}
	intArrayReverse(finished)
	self.nodeList = finished
	self.needSort = false
}

func (self *PoaGraph) testSort() bool {
	if len(self.nodeList) == 0 {
		return false
	}
	// keep track of node we've visited
	seenNodes := make(map[int]bool)

	for _, nodeIdx := range self.nodeList {
		node := self.nodeDict[nodeIdx]
		for inNeighbor, _ := range node.inEdges {
			_, check := seenNodes[inNeighbor]
			if !check {
				return false
			}
		}
		seenNodes[nodeIdx] = true
	}
	return true
}

func (self *PoaGraph) AddSequenceAlignment(pA *PairwiseAlignment) {
	validStringIdxs := make([]int, 0) // lookup how to resize
	// add all of the not-None (not -1) string indices

	strIdxs := pA.stringIdxs
	sequence := pA.sequence
	label := pA.label
	matches := pA.matches

	for _, si := range strIdxs {
		if si >= 0 {
			validStringIdxs = append(validStringIdxs, si)
		}
	}

	firstId, headId, tailId := -1, -1, -1

	startSeqIdx := validStringIdxs[0]
	endSeqIdx := validStringIdxs[len(validStringIdxs)-1]

	// if the new aligned sequence has 'ragged ends' that aren't aligned to the graph, add them
	if startSeqIdx > 0 {
		firstId, headId = self.AddBaseSequence(sequence[:startSeqIdx], label, false)
	}
	if endSeqIdx < len(sequence) {
		tailId, _ = self.AddBaseSequence(sequence[endSeqIdx+1:], label, false)
	}

	//
	for i, sIndex := range strIdxs {
		if sIndex < 0 {
			continue
		}

		base := string(sequence[sIndex])
		matchId := matches[i]

		var nodeId int

		switch {
		case matchId < 0:
			nodeId = self.AddNode(base)
		case self.nodeDict[matchId].base == base:
			nodeId = matchId
		default:
			otherAligns := self.nodeDict[matchId].alignedTo
			foundNode := -1
			// check if this base is aligned to a node that is connected to a matching base
			for _, otherNodeId := range otherAligns {
				if self.nodeDict[otherNodeId].base == base {
					foundNode = otherNodeId
				}
			}
			if foundNode < 0 {
				nodeId = self.AddNode(base)
				otherAligns = append(otherAligns, matchId)
				self.nodeDict[nodeId].alignedTo = append(self.nodeDict[nodeId].alignedTo, otherAligns...)
				for _, otherNodeId := range self.nodeDict[nodeId].alignedTo {
					self.nodeDict[otherNodeId].alignedTo = append(self.nodeDict[otherNodeId].alignedTo, nodeId)
				}
			} else {
				nodeId = foundNode
			}
		}
		self.AddEdge(headId, nodeId, label)
		headId = nodeId
		if firstId < 0 {
			firstId = headId
		}
	}
	self.AddEdge(headId, tailId, label)

	self.TopoSort()

	ok := self.testSort()
	if !ok {
		panic("AddSequenceAlignment: sort failed")
	}

	self.seqs = append(self.seqs, sequence)
	self.labels = append(self.labels, label)
	self.starts = append(self.starts, firstId)
}

func makeAlignmentColumnArray(nbCols int) []string {
	charList := make([]string, nbCols)
	for col := 0; col < nbCols; col++ {
		charList[col] = "-"
	}
	return charList
}

func (self *PoaGraph) GenerateAlignmentStrings() ([]string, []string) {
	columnIndex := make(map[int]int)
	currentColumn := 0

	for _, nodeIdx := range self.nodeList {
		node := self.nodeDict[nodeIdx]
		otherColumns := make([]int, 0)
		for _, other := range node.alignedTo {
			_, contains := columnIndex[other] // todo can be cleaned into one line?
			if contains {
				otherColumns = append(otherColumns, columnIndex[other])
			}
		}

		var foundIdx int
		if len(otherColumns) > 0 {
			foundIdx = intArrayMin(otherColumns)
		} else {
			foundIdx = currentColumn
			currentColumn += 1
		}
		columnIndex[node.id] = foundIdx
	}

	nColumns := currentColumn

	seqNames := make([]string, 0) // todo should know the length of these slices before hand
	alignmentStrings := make([]string, 0)

	for i, start := range self.starts {
		thisLabel := self.labels[i]
		seqNames = append(seqNames, thisLabel)
		curNodeId := start
		charList := makeAlignmentColumnArray(nColumns)

		for curNodeId >= 0 {
			node := self.nodeDict[curNodeId]
			charList[columnIndex[curNodeId]] = node.base
			curNodeId = node.NextNode(thisLabel)
		}

		alnString := strings.Join(charList, "")
		alignmentStrings = append(alignmentStrings, alnString)
	}

	consensusPaths, consensusBases, nbConsensus := self.AllConsensuses(maxFraction)
	for i := 0; i < nbConsensus; i++ {
		path := *consensusPaths[i]
		bases := *consensusBases[i]

		if len(path) != len(bases) {
			e := fmt.Sprintf("Consensus %v doesn't have correct length", i)
			panic(e)
		}

		charList := makeAlignmentColumnArray(nColumns)
		for col := 0; col < len(path); col++ {
			charList[columnIndex[path[col]]] = bases[col]
		}

		alnString := strings.Join(charList, "")
		alignmentStrings = append(alignmentStrings, alnString)
		seqNames = append(seqNames, fmt.Sprintf("Consensus%v", i))
	}

	//fmt.Println("step 3", alignmentStrings)

	return seqNames, alignmentStrings
}

// returns true of tup1 is greater than tup2, if they are equal, returns false
func compareEdgeScores(tup1, tup2 []int) bool {
	if len(tup1) != 3 || len(tup2) != 3 {
		panic("compareEdgeScores: tuples aren't the correct length")
	}
	for i := 0; i < 2; i++ {
		if tup1[i] == tup2[i] {
			continue
		}
		if tup1[i] > tup2[i] {
			return true
		}
	}
	return false
}

func (self *PoaGraph) consensus(exclusions []string) ([]int, []string, [][]string) {
	excludeLabels := make([]string, 0)
	if len(exclusions) > 0 {
		excludeLabels = exclusions
	}

	if self.needSort {
		self.TopoSort()
		ok := self.testSort()
		if !ok {
			panic("Consensus, sort failed")
		}
	}

	nodesInReverse := make([]int, len(self.nodeList))
	copy(nodesInReverse, self.nodeList)
	intArrayReverse(nodesInReverse)

	maxNodeId := intArrayMax(nodesInReverse) + 1
	nextInPath := make([]int, maxNodeId)
	for i := 0; i < len(nextInPath); i++ {
		nextInPath[i] = -1
	}
	scores := make([]int, maxNodeId)

	for _, nodeId := range nodesInReverse {
		bestWeightScoreEdge := []int{-1, -1, -1}

		for neighborId, edge := range self.nodeDict[nodeId].outEdges {
			weight := 0

			// go over the labels, the weight is the total count of labels that aren't in the 'exclude' list
			for _, label := range edge.labels {
				exclude := checkForLabel(excludeLabels, label)
				if !exclude {
					weight += 1
				}
			}

			weightScoreEdge := []int{weight, scores[neighborId], neighborId}
			if compareEdgeScores(weightScoreEdge, bestWeightScoreEdge) {
				bestWeightScoreEdge = weightScoreEdge
			}
		}
		scores[nodeId] = bestWeightScoreEdge[0] + bestWeightScoreEdge[1]
		nextInPath[nodeId] = bestWeightScoreEdge[2]
	}

	//fmt.Println("scores", scores)
	//fmt.Println("nextInPath", nextInPath)

	pos := intArrayArgmax(scores)
	path := make([]int, 0)
	bases := make([]string, 0)
	labels := make([][]string, 0)

	for pos >= 0 {
		path = append(path, pos)
		bases = append(bases, self.nodeDict[pos].base)
		labels = append(labels, self.nodeDict[pos].Labels())
		pos = nextInPath[pos]
	}

	//fmt.Println("path", path)
	//fmt.Println("bases", bases)
	//fmt.Println("labels", labels)

	return path, bases, labels
}

func (self *PoaGraph) AllConsensuses(maxFraction float64) ([]*[]int, []*[]string, int) {
	// containers for accumulating
	allPaths := make([]*[]int, 0)
	allBases := make([]*[]string, 0)
	exclusions := make([]string, 0)
	nbConsensus := 0

	for len(exclusions) < len(self.labels) {
		path, bases, labelLists := self.consensus(exclusions)
		if len(path) > 0 {
			allPaths = append(allPaths, &path)
			allBases = append(allBases, &bases)
			nbConsensus += 1
			labelCounts := make(map[string]int)
			// tally up all of the labels we've seen in this consensus
			for _, labelList := range labelLists {
				for _, label := range labelList {
					labelCounts[label] += 1
				}
			}

			if len(self.seqs) != len(self.labels) {
				panic("number of sequences != number of labels")
			}

			for i := 0; i < len(self.labels); i++ {
				label := self.labels[i]
				seq := self.seqs[i]
				_, contains := labelCounts[label]
				if contains {
					count := float64(labelCounts[label])
					if count >= maxFraction*float64(len(seq)) {
						exclusions = append(exclusions, label)
					}
				}
			}
		}
	}

	return allPaths, allBases, nbConsensus
}
