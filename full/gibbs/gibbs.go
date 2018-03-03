package gibbs

import (
	"github.com/mwortsma/particle_systems2/util/graphutil"
	"github.com/mwortsma/particle_systems2/util/mathutil"
	"github.com/mwortsma/particle_systems2/util/matutil"
	"github.com/mwortsma/particle_systems2/util/probutil"
)

func TimeDistr(
	T int,
	P func(matutil.Vec) float64,
	G graphutil.Graph,
	k int) probutil.TimeDistr {

	n := len(G)
	f := make([]float64, k)

	sum := 0.0
	for _, sigma := range mathutil.QStrings(n, k) {
		p := P(sigma)
		f[sigma[0]] += p
		sum += p
	}

	for k := range f {
		f[k] = f[k] / sum
	}

	t_array := make([]float64, T)
	for i := 0; i < T; i++ {
		t_array[i] = float64(i)
	}

	d := make([][]float64, T)
	for t := 0; t < T; t++ {
		d[t] = f
	}

	return probutil.TimeDistr{1, float64(T), k, d}
}

func FinalNeighborhoodDistr(
	P func(matutil.Vec) float64,
	G graphutil.Graph,
	k int) probutil.PathDistr {

	n := len(G)
	f := make(probutil.PathDistr)

	sum := 0.0
	for _, sigma := range mathutil.QStrings(n, k) {

		local := []int{sigma[0]}
		for _, j := range G[0] {
			local = append(local, sigma[j])
		}

		str := matutil.Vec(local).String()
		if _, ok := f[str]; !ok {
			f[str] = 0.0
		}
		p := P(sigma)
		f[str] += p
		sum += p
	}

	for k := range f {
		f[k] = f[k] / sum
	}

	return f
}
