package gs

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
	Q probutil.RealTransition,
	nu probutil.InitDistr,
	k int,
	G graphutil.Graph,
  tlim int) matutil.Mat {

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
			// get the sum of the neighbors
			neighbors := make([]int, 0)
			for j := 0; j < len(G[i]); j++ {
				neighbors = append(neighbors, X[t-1][G[i][j]])
			}
			X[t][i] = Q(X[t-1][i], neighbors, r.Float64())
		}
	}


  for tstep := 0; tstep < tlim; tstep++ {
    for i := 0; i < n; i++ {
      //X[0][i] = probutil.Sample(nu, r.Float64())
      for t := 1; t < T; t++ {
  			neighbors := make([]int, 0)
  			for j := 0; j < len(G[i]); j++ {
  				neighbors = append(neighbors, X[t-1][G[i][j]])
  			}
  			X[t][i] = Q(X[t-1][i], neighbors, r.Float64())
  		}
  	}
  }
	return X
}


func FinalNeighborhoodDistr(
	T int,
	d int,
	Q probutil.RealTransition,
	nu probutil.InitDistr,
	k int,
	steps int,
	G graphutil.Graph,
  tlim int) probutil.PathDistr {

	f := func() fmt.Stringer {
		X := Realization(T, Q, nu, k, G, tlim)
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
	Q probutil.RealTransition,
	nu probutil.InitDistr,
	k int,
	steps int,
	G graphutil.Graph,
  tlim int) probutil.TimeDistr {

	t_array := make([]float64, T)
	for i := 0; i < T; i++ {
		t_array[i] = float64(i)
	}

	f := func() ([]float64, matutil.Vec) {
		X := Realization(T, Q, nu, k, G, tlim)
		return t_array, X.Col(0)
	}
	return probutil.GetTimeDistrSync(f, 1, float64(T), k, steps)
}

func PathDistr(
	T int,
	Q probutil.RealTransition,
	nu probutil.InitDistr,
	k int,
	steps int,
	G graphutil.Graph,
  tlim int) probutil.PathDistr {

	f := func() fmt.Stringer {
		X := Realization(T, Q, nu, k, G, tlim)
		return X.Col(3)
	}
	return probutil.GetPathDistrSync(f, steps)
}
