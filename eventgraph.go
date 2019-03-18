package graff

// EventGraph is a graph supporting directed edges between nodes.
type EventGraph struct {
	*DirectedGraph
}

// NewDirectedGraph creates a graph of nodes with directed edges.
func NewEventGraph() *EventGraph {
	return &EventGraph{
		&DirectedGraph {
			graph: newGraph(),
			edges: newDirectedEdgeList(),
		},
	}
}

// Copy returns a clone of the directed graph.
func (g *EventGraph) Copy() *EventGraph {
	return &EventGraph{
		&DirectedGraph{
			graph: g.graph.Copy(),
			edges: g.edges.Copy(),
		},
	}
}

// AddEdge adds the edge to the graph.
func (g *EventGraph) AddEdge(from Node, to Node) {
	g.DirectedGraph.AddEdge(to, from);
}

// RemoveEdge removes the edge from the graph.
func (g *EventGraph) RemoveEdge(from Node, to Node) {
	g.DirectedGraph.RemoveEdge(to, from)
}

// EdgeExists checks whether the edge exists within the graph.
func (g *EventGraph) EdgeExists(from Node, to Node) bool {
	return g.DirectedGraph.EdgeExists(to, from)
}




