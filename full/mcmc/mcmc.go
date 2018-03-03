package mcmc

import (
	"fmt"
	"golang.org/x/exp/rand"
	"time"

	"github.com/mwortsma/particle_systems2/util/graphutil"
	"github.com/mwortsma/particle_systems2/util/matutil"
	"github.com/mwortsma/particle_systems2/util/probutil"
)

type NewStateFunc func(matutil.Vec, *rand.Rand) ([]int, []int)
type TransitionProbFunc func(matutil.Vec, []int, []int) float64

func Realization(
	T int,
	NewState NewStateFunc,
	TransitionProb TransitionProbFunc,
	nu probutil.InitDistr,
	n int) matutil.Mat {

	X := matutil.Create(T, n)

	// Ger random number to be used throughout
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	// Initial conditions.
	for i := 0; i < n; i++ {
		X[0][i] = probutil.Sample(nu, r.Float64())
	}

	for t := 1; t < T; t++ {

		sigma := X[t-1]
		sites, newvals := NewState(sigma, r)
		// TODO: double check this slice logic
		copy(X[t], X[t-1])
		if r.Float64() < TransitionProb(sigma, sites, newvals) {
			for j, site := range sites {
				X[t][site] = newvals[j]
			}
		}

	}
	return X
}

func FinalNeighborhoodDistr(
	T int,
	NewState NewStateFunc,
	TransitionProb TransitionProbFunc,
	nu probutil.InitDistr,
	G graphutil.Graph,
	n int,
	steps int,
	d int) probutil.PathDistr {

	f := func() fmt.Stringer {
		X := Realization(T, NewState, TransitionProb, nu, n)
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
	NewState NewStateFunc,
	TransitionProb TransitionProbFunc,
	nu probutil.InitDistr,
	steps int,
	n int,
	k int) probutil.TimeDistr {

	t_array := make([]float64, T)
	for i := 0; i < T; i++ {
		t_array[i] = float64(i)
	}

	f := func() ([]float64, matutil.Vec) {
		X := Realization(T, NewState, TransitionProb, nu, n)
		return t_array, X.Col(0)
	}
	return probutil.GetTimeDistrSync(f, 1, float64(T), k, steps)
}

func PathDistr(
	T int,
	NewState NewStateFunc,
	TransitionProb TransitionProbFunc,
	nu probutil.InitDistr,
	steps int,
	n int) probutil.PathDistr {

	f := func() fmt.Stringer {
		X := Realization(T, NewState, TransitionProb, nu, n)
		return X.Col(0)
	}
	return probutil.GetPathDistrSync(f, steps)
}
