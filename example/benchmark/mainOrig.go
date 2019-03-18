package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"strings"
	"github.com/quan8/cofgra"
)

func main4() {
	fmt.Println("mainOrig")

	graph := graff.NewEventGraph()

	filename := "12015/edges_12015.edges"
	content, err := ioutil.ReadFile(filename)
	if err != nil {
	    //Do something
	}
	lines := strings.Split(string(content), "\n")

	for _,s := range lines {
		if len(s) != 0 {
			tokens := strings.Split(string(s), " ")
			graph.AddEdge(tokens[0], tokens[1])
		}
	}


	dfs, err := graph.DFSSort()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Depth-first search sort: ", dfs)

	coffGra := graph.CoffmanGrahamSorter(WIDTH)

	coffmanGraham, err := coffGra.OrigSort();
	if err != nil {
		log.Fatalln(err)
	}

//	fmt.Println("\n1st ===   Coffman-Graham sort (W=",WIDTH,"): ", coffmanGraham)
	for i := 0; i <len(coffmanGraham); i++  {
		fmt.Println("level=",i,", len=",len(coffmanGraham[i])," ", coffmanGraham[i]);
	}
}
