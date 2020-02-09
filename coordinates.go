package gogm

// coordinate is the position of a graph entity within the graph
type coordinate struct {
	depth         int
	subgraphIndex int
	graphIndex    int
}
