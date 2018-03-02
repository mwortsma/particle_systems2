package meanfield

import (
	//"fmt"
	//"github.com/mwortsma/particle_systems2/util/mathutil"
	//"github.com/mwortsma/particle_systems2/util/matutil"
	"github.com/mwortsma/particle_systems2/util/probutil"
	//"math"
)

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
