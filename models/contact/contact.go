package contact

import 	(
  "github.com/mwortsma/particle_systems2/full/fullgraph"
  "github.com/mwortsma/particle_systems2/full/fulltree"
  "github.com/mwortsma/particle_systems2/util/probutil"
)

// General.
func getLawQ(p,q float64) probutil.LawTransition {
  // transition from s to s_new
  return func(s_new, s int, f probutil.Law) float64 {
  		if s == 1 {
  			if s_new == 1 {
  				return 1-q
  			} else {
  				return q
  			}
  		} else {
  			if s_new == 1 {
  				return p*f[1]
  			} else {
  				return 1-p*f[1]
  			}
  		}
  		return 0.0
  	}
}

func getNeighborQ(p,q float64) probutil.NeighborTransition {
  return probutil.GetNeighborTransition(getLawQ(p,q),2)
}

// Dense.
func DenseTimeDistr(
  T int,
  p,q float64,
  nu probutil.InitDistr,
  steps int,
  n int) probutil.TimeDistr {

  return fullgraph.DenseTimeDistr(T,getNeighborQ(p,q),nu,2,steps,n)

}

func DensePathDistr(
  T int,
  p,q float64,
  nu probutil.InitDistr,
  steps int,
  n int) probutil.PathDistr {

  return fullgraph.DensePathDistr(T,getNeighborQ(p,q),nu,2,steps,n)
}

// Tree
func TreeTimeDistr(
  T int,
  p,q float64,
  d int,
  nu probutil.InitDistr,
	steps int) probutil.TimeDistr {

  return fulltree.TimeDistr(T,d,getNeighborQ(p,q),nu,2,steps)
}

func TreePathDistr(
  T int,
  p,q float64,
  d int,
  nu probutil.InitDistr,
	steps int) probutil.PathDistr {

  return fulltree.PathDistr(T,d,getNeighborQ(p,q),nu,2,steps)
}

// Local
