package mcmc

import (
  "fmt"
	"golang.org/x/exp/rand"
	"time"

	"github.com/mwortsma/particle_systems2/util/graphutil"
	"github.com/mwortsma/particle_systems2/util/matutil"
	"github.com/mwortsma/particle_systems2/util/probutil"
)

func Realization(
	T int,
	NewState func(v matutil.Vec, r *rand.Rand) matutil.Vec,
  E func(v matutil.Vec, G graphutil.Graph) float64,
  transitionProb func(eTau float64, eSigma float64) float64,
	nu probutil.InitDistr,
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

    sigma := X[t-1]
    eSigma := E(sigma,G)
		tau := NewState(sigma,r)
    eTau := E(tau,G)

    if eTau < eSigma || r.Float64() < transitionProb(eTau, eSigma) {
      X[t] = tau
    } else {
      X[t] = sigma
    }

	}
	return X
}


func FinalNeighborhoodDistr(
	T int,
  NewState func(v matutil.Vec, r *rand.Rand) matutil.Vec,
  E func(v matutil.Vec, G graphutil.Graph) float64,
  transitionProb func(eTau float64, eSigma float64) float64,
	nu probutil.InitDistr,
	G graphutil.Graph,
  steps int,
  d int) probutil.PathDistr {

	f := func() fmt.Stringer {
		X := Realization(T, NewState, E, transitionProb, nu, G)
		v := []int{X[0][T-1]}
		for i, j := range G[0] {
			if d > 0 && i >= d {
				break
			}
			v = append(v, X[T-1][j])
		}
		return matutil.Vec(v)
	}
	return probutil.GetPathDistrSync(f, steps)
}

func TimeDistr(
  T int,
  NewState func(v matutil.Vec, r *rand.Rand) matutil.Vec,
  E func(v matutil.Vec, G graphutil.Graph) float64,
  transitionProb func(eTau float64, eSigma float64) float64,
	nu probutil.InitDistr,
  steps int,
	G graphutil.Graph,
  k int) probutil.TimeDistr {

	t_array := make([]float64, T)
	for i := 0; i < T; i++ {
		t_array[i] = float64(i)
	}

	f := func() ([]float64, matutil.Vec) {
		X := Realization(T, NewState, E, transitionProb, nu, G)
		return t_array, X.Col(0)
	}
	return probutil.GetTimeDistrSync(f, 1, float64(T), k, steps)
}

func PathDistr(
  T int,
  NewState func(v matutil.Vec, r *rand.Rand) matutil.Vec,
  E func(v matutil.Vec, G graphutil.Graph) float64,
  transitionProb func(eTau float64, eSigma float64) float64,
	nu probutil.InitDistr,
  steps int,
	G graphutil.Graph) probutil.PathDistr {

	f := func() fmt.Stringer {
		X := Realization(T, NewState, E, transitionProb, nu, G)
		return X.Col(0)
	}
	return probutil.GetPathDistrSync(f, steps)
}
