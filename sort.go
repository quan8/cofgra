package graff

import (
	"fmt"
	"errors"
)

// Errors relating to the DFSSorter.
var (
	ErrCyclicGraph = errors.New("The graph cannot be cyclic")
)

// DFSSorter topologically sorts a directed graph's nodes based on the
// directed edges between them using the Depth-first search algorithm.
type DFSSorter struct {
	graph      *DirectedGraph
	sorted     []Node
	visiting   map[Node]bool
	discovered map[Node]bool
}

// NewDFSSorter returns a new DFS sorter.
func NewDFSSorter(graph *DirectedGraph) *DFSSorter {
	return &DFSSorter{
		graph: graph,
	}
}

func (s *DFSSorter) init() {
	s.sorted = make([]Node, 0, s.graph.NodeCount())
	s.visiting = make(map[Node]bool)
	s.discovered = make(map[Node]bool, s.graph.NodeCount())
}

// Sort returns the sorted nodes.
func (s *DFSSorter) Sort() ([]Node, error) {
	s.init()

	// > while there are unmarked nodes do
	for _, node := range s.graph.Nodes() {
		if err := s.visit(node); err != nil {
			return nil, err
		}
	}

	// as the nodes were appended to the slice for performance reasons,
	// rather than prepended as correctly stated by the algorithm,
	// we need to reverse the sorted slice
	for i, j := 0, len(s.sorted)-1; i < j; i, j = i+1, j-1 {
		s.sorted[i], s.sorted[j] = s.sorted[j], s.sorted[i]
	}

	return s.sorted, nil
}

// See https://en.wikipedia.org/wiki/Topological_sorting#Depth-first_search
func (s *DFSSorter) visit(node Node) error {
	// > if n has a permanent mark then return
	if discovered, ok := s.discovered[node]; ok && discovered {
		return nil
	}
	// > if n has a temporary mark then stop (not a DAG)
	if visiting, ok := s.visiting[node]; ok && visiting {
		return ErrCyclicGraph
	}

	// > mark n temporarily
	s.visiting[node] = true

	// > for each node m with an edge from n to m do
	for _, outgoing := range s.graph.OutgoingEdges(node) {
		if err := s.visit(outgoing); err != nil {
			return err
		}
	}

	s.discovered[node] = true
	delete(s.visiting, node)

	s.sorted = append(s.sorted, node)
	return nil
}

// DFSSort returns the graph's nodes in topological order based on the
// directed edges between them using the Depth-first search algorithm.
func (g *DirectedGraph) DFSSort() ([]Node, error) {
	sorter := NewDFSSorter(g)
	return sorter.Sort()
}

// Errors relating to the CoffmanGrahamSorter.
var (
	ErrDependencyOrder = errors.New("The topological dependency order is incorrect")
)

// CoffmanGrahamSorter sorts a graph's nodes into a sequence of levels,
// arranging so that a node which comes after another in the order is
// assigned to a lower level, and that a level never exceeds the width.
// See https://en.wikipedia.org/wiki/Coffman–Graham_algorithm
type CoffmanGrahamSorter struct {
	graph *DirectedGraph
	width int

	layers [][]Node
	levels map[Node]int
	level int
}

// NewCoffmanGrahamSorter returns a new Coffman-Graham sorter.
func NewCoffmanGrahamSorter(graph *DirectedGraph, width int) *CoffmanGrahamSorter {
	layers := make([][]Node, 0)
	levels := make(map[Node]int, 0)
	level := 0

	return &CoffmanGrahamSorter{
		graph: graph,
		width: width,
		layers : layers,
		levels : levels,
		level : level,
	}
}


// Sort returns the sorted nodes.
// This version tries to optimize for directed graph (not reverse graph)
func (s *CoffmanGrahamSorter) Sort() ([][]Node, error) {
	// create a copy of the graph and remove transitive edges
	//fmt.Println("-start-- layers", s.layers)

	reduced := s.graph.Copy()
	reduced.RemoveTransitives()

	// topologically sort the graph nodes
	nodes, err := reduced.DFSSort()
	if err != nil {
		return nil, err
	}

	//fmt.Println("--- DFS nodes", nodes);
	layers := s.layers
	levels := s.levels
	level := s.level

	//fmt.Println("--- level", level)
	//fmt.Println("--- levels", levels, "layers", layers)

	for _, node := range nodes {
		//fmt.Println("--- node", node)
		_, ok := levels[node]
		if ok {
			// if already assigned a level, dont need to assign again
			continue
		}

		fmt.Println("--- node  -- not assigned", node)

		dependantLevel := -1
		for _, dependant := range reduced.IncomingEdges(node) {

			level, ok := levels[dependant]
			//fmt.Println("--- dependant", node, "level", level)

			if !ok {
				return nil, ErrDependencyOrder
			}
			if level > dependantLevel {
				dependantLevel = level
			}

			//fmt.Println("--- dependantLevel", dependantLevel)
		}

		level := -1

		// find the first ununfilled layer outgoing the dependent layer
		// skip this if the dependent layer is the last

		//fmt.Println("--- dependantLevel=", dependantLevel, ", len(layers) - 1 = ", len(layers) - 1)

		if dependantLevel < len(layers)-1 {
			//fmt.Println("--- len(layers)-1=", len(layers)-1)

			for i := dependantLevel + 1; i < len(layers); i++ {
				// ensure the layer doesn't exceed the desired width
				if len(layers[i]) < s.width {
					level = i
					break
				}
			}
		}

		//fmt.Println("HEEERR level ", level);
		// create a new layer new none was found
		if level == -1 {
			//fmt.Println("HEEERR level ", level, "layers", layers);

			layers = append(layers, make([]Node, 0, 1))
			level = len(layers) - 1
		}

		//fmt.Println("--== level ", level);

		layers[level] = append(layers[level], node)
		levels[node] = level
	}

	s.level = level
	//fmt.Println("--== level ", level);
	s.layers = layers;
	//fmt.Println("ending--- levels", levels, "layers", layers)

	return layers, nil
}

// Sort returns the sorted nodes.
// This version is orginal impl for directed graph (not reverse graph)
func (s *CoffmanGrahamSorter) OrigSort() ([][]Node, error) {
	// create a copy of the graph and remove transitive edges
	reduced := s.graph.Copy()
	reduced.RemoveTransitives()

	// topologically sort the graph nodes
	nodes, err := reduced.DFSSort()
	if err != nil {
		return nil, err
	}

	layers := make([][]Node, 0)
	levels := make(map[Node]int, len(nodes))

	for _, node := range nodes {
		dependantLevel := -1
		for _, dependant := range reduced.IncomingEdges(node) {
			level, ok := levels[dependant]
			if !ok {
				return nil, ErrDependencyOrder
			}
			if level > dependantLevel {
				dependantLevel = level
			}
		}

		level := -1
		// find the first ununfilled layer outgoing the dependent layer
		// skip this if the dependent layer is the last
		if dependantLevel < len(layers)-1 {
			for i := dependantLevel + 1; i < len(layers); i++ {
				// ensure the layer doesn't exceed the desired width
				if len(layers[i]) < s.width {
					level = i
					break
				}
			}
		}
		// create a new layer new none was found
		if level == -1 {
			layers = append(layers, make([]Node, 0, 1))
			level = len(layers) - 1
		}

		layers[level] = append(layers[level], node)
		levels[node] = level
	}

	return layers, nil
}


func (g *DirectedGraph) CoffmanGrahamSorter(width int) (*CoffmanGrahamSorter) {
	sorter := NewCoffmanGrahamSorter(g, width)
	return sorter
}

// CoffmanGrahamSort sorts the graph's nodes into a sequence of levels,
// arranging so that a node which comes after another in the order is
// assigned to a lower level, and that a level never exceeds the specified width.
func (g *DirectedGraph) CoffmanGrahamSort(width int) ([][]Node, error) {
	sorter := NewCoffmanGrahamSorter(g, width)
	return sorter.Sort()
}

