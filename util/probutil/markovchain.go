package probutil

import (
	"github.com/mwortsma/particle_systems2/util/matutil"
	"github.com/mwortsma/particle_systems2/util/mathutil"
)

type InitFunc func(matutil.Vec) float64
type InitDistr []float64

type NeighborTransition func(int, int, matutil.Vec) float64

type Law []float64

type LawTransition func(int, int, Law) float64

type RealTransition func(int, matutil.Vec, float64) int

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

// Try not to call this function.
func GetLawTransition(G NeighborTransition, k int, d int) LawTransition {
	return func(s_new, s int, l Law) float64 {
		prob := 0.0
		for _, neighbors := range mathutil.QStrings(d,k) {
			p := G(s_new, s, neighbors)
			for _, n := range neighbors {
				p *= l[n]
			}
			prob += p
		}
		return prob
  }
}
