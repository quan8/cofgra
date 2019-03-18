package graff

import (
//	"fmt"
)


// CoffmanGrahamSorter sorts a graph's nodes into a sequence of levels,
// arranging so that a node which comes after another in the order is
// assigned to a lower level, and that a level never exceeds the width.
// See https://en.wikipedia.org/wiki/Coffman–Graham_algorithm
type OptimizedCoffmanGrahamSorter struct {
	graph *DirectedGraph
	width int

	layers [][]Node
	levels map[Node]int
	level int
}

// Sort returns the sorted nodes.
// This version is to optimize for reverse graph (not the original directed graph)
func (s *OptimizedCoffmanGrahamSorter) EventSort() ([][]Node, error) {
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

//	layers := make([][]Node, 0)
//	levels := make(map[Node]int, len(nodes))

	layers := s.layers
	levels := s.levels
	//level := s.level

//	fmt.Println("---SSS level", level)
//	fmt.Println("---SSS levels", levels, "layers", layers)

	maxLevel := -1

	for _, node := range nodes {
		//fmt.Println("--- node", node)

		_, ok := levels[node]
		if ok {
			// if already assigned a level, dont need to assign again
			continue
		}

		//fmt.Println("--- node  -- not assigned", node)

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
//			fmt.Println("--- dependantLevel", dependantLevel)
		}

		level := -1

		// find the first ununfilled layer outgoing the dependent layer
		// skip this if the dependent layer is the last

//		fmt.Println("--- dependantLevel=", dependantLevel, ", len(layers) - 1 = ", len(layers) - 1)

		if dependantLevel < len(layers)-1 {
//			fmt.Println("--- len(layers)-1=", len(layers)-1)

			for i := dependantLevel + 1; i < len(layers); i++ {
				// ensure the layer doesn't exceed the desired width
				if len(layers[i]) < s.width {
					level = i
					break
				}
			}
		}

//		fmt.Println("HEEERR level ", level);

		// create a new layer new none was found
		if level == -1 {
			layers = append(layers, make([]Node, 0, 1))
			level = len(layers) - 1

//			fmt.Println("CREATING NEW level=", level, " , layers=", layers);
		}

//		fmt.Println("- LOOOOOOOP --- level ", level);

		layers[level] = append(layers[level], node)
		levels[node] = level

		if maxLevel == -1 || maxLevel < level {
			maxLevel = level
		}
	}

	s.level = maxLevel;
	s.layers = layers;
	s.levels = levels;

//	fmt.Println("ENDDDING level=", s.level);
//	fmt.Println("ENDDDING --- levels", s.levels,
//		 "\n layers", s.layers)

	return layers, nil
}

func (g *DirectedGraph) OptimizedCoffmanGrahamSorter(width int) (*OptimizedCoffmanGrahamSorter) {
	sorter := NewOptimizedCoffmanGrahamSorter(g, width)
	return sorter
}

// NewCoffmanGrahamSorter returns a new Coffman-Graham sorter.
func NewOptimizedCoffmanGrahamSorter(graph *DirectedGraph, width int) *OptimizedCoffmanGrahamSorter {
	layers := make([][]Node, 0)
	levels := make(map[Node]int, 0)
	level := 0

	return &OptimizedCoffmanGrahamSorter{
		graph: graph,
		width: width,
		layers : layers,
		levels : levels,
		level : level,
	}
}
