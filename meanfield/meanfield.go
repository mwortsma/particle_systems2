package meanfield

import (
	//"fmt"
	"github.com/mwortsma/particle_systems2/util/mathutil"
	//"github.com/mwortsma/particle_systems2/util/matutil"
	"github.com/mwortsma/particle_systems2/util/probutil"
	//"math"
)

// Main Algorithm
func getLaw(
  T int,
  Q probutil.LawTransition,
  nu probutil.InitDistr,
  k int) [][]float64 {

  law := make([][]float64, T)
  law[0] = make([]float64, k)
  for i, v := range nu {
    law[0][i] = v
  }

  for t := 1; t < T; t++ {
    law[t] = make([]float64, k)
    for i := 0; i < k; i++ {
      for j := 0; j < k; j++ {
        law[t][i] += Q(i,j,law[t-1])*law[t-1][j]
      }
    }
  }
  return law
}


// Gets the distriution over the local neighborhood at the end
func FinalNeighborhoodDistr(
  T int,
  Q probutil.LawTransition,
  nu probutil.InitDistr,
  k int,
  d int) probutil.PathDistr {

	law := getLaw(T,Q,nu,k)
  f := law[T-1]

	distr := make(probutil.PathDistr)

	for _, v := range mathutil.QStrings(d,k) {
    prob := 1.0
    for _, vi := range v {
      prob *= f[vi]
    }
    distr[v.String()] = prob
	}
  return distr
}

func TimeDistr(
	T int,
	Q probutil.LawTransition,
	nu probutil.InitDistr,
	k int) probutil.TimeDistr {

	return probutil.TimeDistr{1, float64(T), k, getLaw(T,Q,nu,k)}
}

func PathDistr(
	T int,
	Q probutil.LawTransition,
	nu probutil.InitDistr,
  k int) probutil.PathDistr {

	law := getLaw(T,Q,nu,k)
	f := make(probutil.PathDistr)

	for _, path := range mathutil.QStrings(T,k) {
    prob := nu[path[0]]
    for t := 1; t < T; t++ {
      prob *= Q(path[t],path[t-1],law[t-1])
    }
    f[path.String()] = prob
  }

	return f
}
