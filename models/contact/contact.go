package contact

import 	(
  "github.com/mwortsma/particle_systems2/full/fullgraph"
  "github.com/mwortsma/particle_systems2/full/fulltree"
  "github.com/mwortsma/particle_systems2/local"
  "github.com/mwortsma/particle_systems2/util/matutil"
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
  	}
}

func getRealQ(p,q float64) probutil.RealTransition {
  // transition from s to s_new
  return func(prev int, neighbors matutil.Vec, r float64) int {
  		if prev == 1 {
        // transition back with probability q
  			if r < q {
  				return 0
  			} else {
  				return 1
  			}
  		} else {
  			f := probutil.Freq(neighbors, 2)
        if r < p*f[1] {
          return 1
        } else {
          return 0
        }
  		}
  	}
}

func getNeighborQ(p,q float64) probutil.NeighborTransition {
  return probutil.GetNeighborTransition(getLawQ(p,q),2)
}

// Dense.
func DenseFinalNeighborhoodDistr(
  T int,
  p,q float64,
  nu probutil.InitDistr,
  steps int,
  n int,
  d int) probutil.PathDistr {

  return fullgraph.DenseFinalNeighborhoodDistr(T,getRealQ(p,q),nu,2,steps,n,d)
}

func DenseTimeDistr(
  T int,
  p,q float64,
  nu probutil.InitDistr,
  steps int,
  n int) probutil.TimeDistr {

  return fullgraph.DenseTimeDistr(T,getRealQ(p,q),nu,2,steps,n)

}

func DensePathDistr(
  T int,
  p,q float64,
  nu probutil.InitDistr,
  steps int,
  n int) probutil.PathDistr {

  return fullgraph.DensePathDistr(T,getRealQ(p,q),nu,2,steps,n)
}

// Tree
func TreeFinalNeighborhoodDistr(
  T int,
  p,q float64,
  d int,
  nu probutil.InitDistr,
	steps int,
  depth int) probutil.PathDistr {

  return fulltree.FinalNeighborhoodDistr(T,d,getRealQ(p,q),nu,2,steps, depth)
}

func TreeTimeDistr(
  T int,
  p,q float64,
  d int,
  nu probutil.InitDistr,
	steps int,
  depth int) probutil.TimeDistr {

  return fulltree.TimeDistr(T,d,getRealQ(p,q),nu,2,steps, depth)
}

func TreePathDistr(
  T int,
  p,q float64,
  d int,
  nu probutil.InitDistr,
	steps int,
  depth int) probutil.PathDistr {

  return fulltree.PathDistr(T,d,getRealQ(p,q),nu,2,steps,depth)
}

// Local
func LocalFinalNeighborhoodDistr(
	T int,
	tau int,
	d int,
	p,q float64,
	nu probutil.InitFunc) probutil.PathDistr {

  return local.FinalNeighborhoodDistr(T,tau,d,getNeighborQ(p,q),nu,2)
}


func LocalTimeDistr(
	T int,
	tau int,
	d int,
	p,q float64,
	nu probutil.InitFunc) probutil.TimeDistr {

  return local.TimeDistr(T,tau,d,getNeighborQ(p,q),nu,2)
}

func LocalPathDistr(
	T int,
	tau int,
	d int,
	p,q float64,
	nu probutil.InitFunc) probutil.PathDistr {

  return local.PathDistr(T,tau,d,getNeighborQ(p,q),nu,2)
}
