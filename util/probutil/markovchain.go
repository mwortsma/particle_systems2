package probutil

import (
	"github.com/mwortsma/particle_systems2/util/matutil"
)

type InitFunc func(matutil.Vec) float64
type InitDistr []float64

type NeighborTransition func(int, int, matutil.Vec) float64

type Law []float64

type LawTransition func(int, int, Law) float64

func GetInitFunc(nu InitDistr) InitFunc {
	return func(v matutil.Vec) float64 {
		p := 1.0
		for _, vi := range v {
			p *= nu[vi]
		}
		return p
	}
}

func GetNeighborTransition(G LawTransition, k int) NeighborTransition {
	return func(s_new, s int, v matutil.Vec) float64 {
		return G(s_new, s, Freq(v,k))
  }
}
