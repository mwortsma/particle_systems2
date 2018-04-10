package graphutil

import (
	"bytes"
	"fmt"
	"golang.org/x/exp/rand"
)

type Graph [][]int

func GetGraph(graph string, n int, erp float64, seed int) Graph {
	switch graph {
	case "ring":
		return Ring(n)
	case "complete":
		return Complete(n)
	case "ER":
		return ER(n, erp, seed)
	}
	return Ring(n)
}

// Ring makes a ring Graph with n nodes.
func Ring(n int) Graph {
	G := make(Graph, n)
	for i := 0; i < n; i++ {
		G[i] = []int{(i - 1 + n) % n, (i + 1) % n}
	}
	return G
}

// Complete makes a complete Graph with n nodes.
func Complete(n int) Graph {
	G := make(Graph, n)
	for i := 0; i < n; i++ {
		G[i] = make([]int, 0)
		for j := 0; j < n; j++ {
			if j != i {
				G[i] = append(G[i], j)
			}
		}
	}
	return G
}

// ER random graph.
func ER(n int, erp float64, seed int) Graph {
	r := rand.New(rand.NewSource(uint64(seed)))
	G := make(Graph, n)
	for i := 0; i < n; i++ {
		G[i] = make([]int, 0)
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if r.Float64() < erp {
				G[i] = append(G[i], j)
				G[j] = append(G[j], i)
			}
		}
	}
	fmt.Println("Node 0 has degree", len(G[0]))
	return G
}

// String gets a description of G
func (G Graph) Print() {
	var buffer bytes.Buffer
	for node, neighbors := range G {
		buffer.WriteString(fmt.Sprintf("Node %d: %v\n", node, neighbors))
	}
	fmt.Println(buffer.String())
}
