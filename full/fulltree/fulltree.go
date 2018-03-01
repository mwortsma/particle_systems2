package fulltree

import (
	"fmt"
	"github.com/mwortsma/particle_systems2/util/matutil"
	"github.com/mwortsma/particle_systems2/util/probutil"
	"golang.org/x/exp/rand"
	"time"
)

type node struct {
	children []*node
	parent   *node
	state    matutil.Vec
	is_leaf  bool
	is_root  bool
}

func Realization(
  T int,
  d int,
  Q probutil.NeighborTransition,
  nu probutil.InitDistr,
  k int) matutil.Vec {
	// Ger random number to be used throughout
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	// create tree
	var root node
	root.createNode(T, d, nu, &node{}, T-1, true, r)

	for t := 1; t < T; t++ {
		// transition will be called for the whole tree recursively
		root.transition(t, d, Q, k, r)
	}

	return root.state
}

func TimeDistr(
  T int,
  d int,
  Q probutil.NeighborTransition,
  nu probutil.InitDistr,
  k int,
  steps int) probutil.TimeDistr {

	t_array := make([]float64, T)
	for i := 0; i < T; i++ {
		t_array[i] = float64(i)
	}

	f := func() ([]float64, matutil.Vec) {
		return t_array, Realization(T, d, Q, nu, k)
	}
	return probutil.GetTimeDistrSync(f, 1, float64(T), 2, steps)
}

func PathDistr(
  T int,
  d int,
  Q probutil.NeighborTransition,
  nu probutil.InitDistr,
  k int,
	steps int) probutil.PathDistr {

	f := func() fmt.Stringer {
		return Realization(T, d, Q, nu, k)
	}
	return probutil.GetPathDistrSync(f, steps)
}

// Helpers
func (n *node) createNode(
	T int,
	d int,
	nu probutil.InitDistr,
	parent *node,
	depth int,
	is_root bool,
	r *rand.Rand) {

	// set parent
	n.is_root = is_root
	if !n.is_root {
		n.parent = parent
	}
	// create children
	if depth == 0 {
		n.is_leaf = true
	} else {
		n.children = make([]*node, d-1)
		for c := 0; c < d-1; c++ {
			var child node
			child.createNode(T, d, nu, n, depth-1, false, r)
			n.children[c] = &child
		}
		if n.is_root {
			var child node
			child.createNode(T, d, nu, n, depth-1, false, r)
			n.children = append(n.children, &child)
		}
	}
	// create state
	n.state = make(matutil.Vec, T)

	// TODO initial conditions
}

func (n *node) transition(
  t int,
  d int,
  Q probutil.NeighborTransition,
  k int,
  r *rand.Rand) {

	n.state[t] = n.state[t-1]

  neighbors := make([]int, 0)
  if !n.is_root {
    neighbors = append(neighbors, n.parent.state[t-1])
  }
  for _, c := range n.children {
    neighbors = append(neighbors, c.state[t-1])
  }

  // TODO Sample from Q

	// call transition on children
	for _, c := range n.children {
		c.transition(t, d, Q, k, r)
	}
}
