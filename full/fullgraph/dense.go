package fullgraph

import (
	"github.com/mwortsma/particle_systems2/util/graphutil"
	"github.com/mwortsma/particle_systems2/util/probutil"
)

// Graph specific
func DenseTimeDistr(
  T int,
  Q probutil.RealTransition,
  nu probutil.InitDistr,
  k int,
  steps int,
  n int) probutil.TimeDistr {

  return TimeDistr(T,n-1,Q,nu,k,steps,graphutil.Complete(n))
}

func DensePathDistr(
  T int,
  Q probutil.RealTransition,
  nu probutil.InitDistr,
  k int,
  steps int,
  n int) probutil.PathDistr {

  return PathDistr(T,n-1,Q,nu,k,steps,graphutil.Complete(n))
}
