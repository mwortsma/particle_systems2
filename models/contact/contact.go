package contact

import 	(
  //"github.com/mwortsma/particle_systems2/full/fullgraph"
  //"github.com/mwortsma/particle_systems2/full/fulltree"
  "github.com/mwortsma/particle_systems2/util/probutil"
  "github.com/mwortsma/particle_systems2/util/matutil"
)

// General.
func getNeighborQ(p,q float64, d int) probutil.NeighborTransition {
  // transition from s to k
  return func(k, s int, v matutil.Vec) float64 {
  		if s == 1 {
  			if k == 1 {
  				return 1-q
  			} else {
  				return q
  			}
  		} else {
  			sum_neighbors := 0
  			for i := 0; i < len(v); i++ {
  				sum_neighbors += v[i]
  			}
  			if k == 1 {
  				return (p/float64(d))*float64(sum_neighbors)
  			} else {
  				return 1-(p/float64(d))*float64(sum_neighbors)
  			}
  		}
  		return 0.0
  	}
}

func getLawQ(p,q float64, d int) probutil.LawTransition {
  // transition from s to k
  return func(k, s int, f probutil.Law) float64 {
  		if s == 1 {
  			if k == 1 {
  				return 1-q
  			} else {
  				return q
  			}
  		} else {
  			if k == 1 {
  				return p*f[1]
  			} else {
  				return 1-p*f[1]
  			}
  		}
  		return 0.0
  	}
}

// Full
