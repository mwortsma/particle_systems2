package fullgraph

import (
	"github.com/mwortsma/particle_systems2/util/graphutil"
	"github.com/mwortsma/particle_systems2/util/probutil"
)

// Untested.

// Graph specific
func RingFinalNeighborhoodDistr(
	T int,
	Q probutil.RealTransition,
	nu probutil.InitDistr,
	k int,
	steps int,
	n int) probutil.PathDistr {

	return FinalNeighborhoodDistr(T, 2, Q, nu, k, steps, graphutil.Ring(n))
}

func RingTimeDistr(
	T int,
	Q probutil.RealTransition,
	nu probutil.InitDistr,
	k int,
	steps int,
	n int) probutil.TimeDistr {

	return TimeDistr(T, 2, Q, nu, k, steps, graphutil.Ring(n))
}

func RingPathDistr(
	T int,
	Q probutil.RealTransition,
	nu probutil.InitDistr,
	k int,
	steps int,
	n int) probutil.PathDistr {

	return PathDistr(T, 2, Q, nu, k, steps, graphutil.Ring(n))
}
