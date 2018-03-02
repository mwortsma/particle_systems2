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
	Q probutil.RealTransition,
	nu probutil.InitDistr,
	k int,
	depth int) node {
	// Ger random number to be used throughout
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	if depth < 0 {
		depth = T - 1
	}
	// create tree
	var root node
	root.createNode(T, d, nu, &node{}, depth, true, r)

	for t := 1; t < T; t++ {
		// transition will be called for the whole tree recursively
		root.transition(t, d, Q, k, r)
	}

	return root
}

func FinalNeighborhoodDistr(
	T int,
	d int,
	Q probutil.RealTransition,
	nu probutil.InitDistr,
	k int,
	steps int,
	depth int) probutil.PathDistr {

	f := func() fmt.Stringer {
		x := Realization(T, d, Q, nu, k, depth)
		v := []int{x.state[T-1]}
		for _, c := range x.children {
			v = append(v, c.state[T-1])
		}
		return matutil.Vec(v)
	}
	return probutil.GetPathDistrSync(f, steps)
}

func TimeDistr(
	T int,
	d int,
	Q probutil.RealTransition,
	nu probutil.InitDistr,
	k int,
	steps int,
	depth int) probutil.TimeDistr {

	t_array := make([]float64, T)
	for i := 0; i < T; i++ {
		t_array[i] = float64(i)
	}

	f := func() ([]float64, matutil.Vec) {
		x := Realization(T, d, Q, nu, k, depth)
		return t_array, x.state
	}
	return probutil.GetTimeDistrSync(f, 1, float64(T), 2, steps)
}

func PathDistr(
	T int,
	d int,
	Q probutil.RealTransition,
	nu probutil.InitDistr,
	k int,
	steps int,
	depth int) probutil.PathDistr {

	f := func() fmt.Stringer {
		x := Realization(T, d, Q, nu, k, depth)
		return x.state
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
	n.state[0] = probutil.Sample(nu, r.Float64())
}

func (n *node) transition(
	t int,
	d int,
	Q probutil.RealTransition,
	k int,
	r *rand.Rand) {

	neighbors := make([]int, 0)
	if !n.is_root {
		neighbors = append(neighbors, n.parent.state[t-1])
	}
	for _, c := range n.children {
		neighbors = append(neighbors, c.state[t-1])
	}

	n.state[t] = Q(n.state[t-1], neighbors, r.Float64())

	// call transition on children
	for _, c := range n.children {
		c.transition(t, d, Q, k, r)
	}
}
