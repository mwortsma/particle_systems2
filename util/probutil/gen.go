package probutil

import (
	"github.com/mwortsma/particle_systems2/util/matutil"
)

func Freq(v matutil.Vec, k int) Law {
	s := make([]float64, k)
	denom := float64(len(v))
	for _, vi := range v {
		s[vi] += 1.0 / denom
	}
	return s
}
