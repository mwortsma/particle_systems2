package fullgraph

import (
	"github.com/mwortsma/particle_systems2/util/graphutil"
	"github.com/mwortsma/particle_systems2/util/probutil"
)

// Graph specific
func DenseTimeDistr(
  T int,
  d int,
  Q probutil.NeighborTransition,
  nu probutil.InitDistr,
  k int,
  steps int,
  n int,
  G graphutil.Graph) probutil.TimeDistr {

  return TimeDistr(T,d,Q,nu,k,steps,graphutil.Complete(n))
}

func DensePathDistr(
  T int,
  d int,
  Q probutil.NeighborTransition,
  nu probutil.InitDistr,
  k int,
  steps int,
  n int,
  G graphutil.Graph) probutil.PathDistr {

  return PathDistr(T,d,Q,nu,k,steps,graphutil.Complete(n))
}
