package fullgraph

import (
	"fmt"
	"github.com/mwortsma/particle_systems2/util/graphutil"
	"github.com/mwortsma/particle_systems2/util/matutil"
	"github.com/mwortsma/particle_systems2/util/probutil"
	"golang.org/x/exp/rand"
	"time"
)

func Realization(
  T int,
  d int,
  Q probutil.NeighborTransition,
  nu probutil.InitDistr,
  k int,
  G graphutil.Graph) matutil.Mat {

	n := len(G)
	X := matutil.Create(T, n)

	// Ger random number to be used throughout
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	// Initial conditions.
	for i := 0; i < n; i++ {
		X[0][i] = probutil.Sample(nu, r.Float64())
	}

	for t := 1; t < T; t++ {
		for i := 0; i < n; i++ {
			X[t][i] = X[t-1][i]
			// get the sum of the neighbors
			neighbors := make([]int, 0)
			for j := 0; j < len(G[i]); j++ {
				neighbors = append(neighbors, X[t-1][G[i][j]])
			}

      // TODO Sample from Q

		}
	}
	return X
}

func TimeDistr(
  T int,
  d int,
  Q probutil.NeighborTransition,
  nu probutil.InitDistr,
  k int,
  steps int,
  G graphutil.Graph) probutil.TimeDistr {

	t_array := make([]float64, T)
	for i := 0; i < T; i++ {
		t_array[i] = float64(i)
	}

	f := func() ([]float64, matutil.Vec) {
    X := Realization(T, d, Q, nu, k, G)
		return t_array, X.Col(0)
	}
	return probutil.GetTimeDistrSync(f, 1, float64(T), 2, steps)
}

func PathDistr(
  T int,
  d int,
  Q probutil.NeighborTransition,
  nu probutil.InitDistr,
  k int,
  steps int,
  G graphutil.Graph) probutil.PathDistr {

	f := func() fmt.Stringer {
    X := Realization(T, d, Q, nu, k, G)
		return X.Col(0)
	}
	return probutil.GetPathDistrSync(f, steps)
}
