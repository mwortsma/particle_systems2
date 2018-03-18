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
  Qr probutil.RealTransition,
	Q probutil.NeighborTransition,
	nu probutil.InitDistr,
	k int,
	G graphutil.Graph,
  tlim int) matutil.Mat {

  j := 2
  n := len(G)
  numrotations := 50

	X := matutil.Create(T, n)

	// Ger random number to be used throughout
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	// Initial conditions.
  for t := 0; t < T; t++ {
  	for i := 0; i < n; i++ {
  	    X[t][i] = probutil.Sample(nu, r.Float64())
  	}
  }

  for rots := 0; rots < numrotations; rots++ {
    for tstep := 0; tstep < tlim; tstep++ {
      // choose t at random
      t := r.Intn(T)
      // calculate prob X[t][j] = 0
      if t > 0 && t < T-1 {
        prob0 := Q(0, X[t-1][j], []int{X[t-1][j-1], X[t-1][j+1]})
        prob0 *= Q(X[t+1][j], 0, []int{X[t][j-1], X[t][j+1]})
        prob0 *= Q(X[t+1][j+1], X[t][j+1], []int{0, X[t][j+2]})
        prob0 *= Q(X[t+1][j-1], X[t][j-1], []int{0, X[t][j-2]})

        prob1 := Q(1, X[t-1][j], []int{X[t-1][j-1], X[t-1][j+1]})
        prob1 *= Q(X[t+1][j], 1, []int{X[t][j-1], X[t][j+1]})
        prob1 *= Q(X[t+1][j+1], X[t][j+1], []int{1, X[t][j+2]})
        prob1 *= Q(X[t+1][j-1], X[t][j-1], []int{1, X[t][j-2]})

        var ratio float64
        if X[t][j] == 0 {
          ratio = prob1/prob0
        } else {
          ratio = prob0/prob1
        }
        if r.Float64() < ratio {
          X[t][j] = 1-X[t][j]
        }

      }
    }
    // Rotate
    for t := 0; t < T; t++ {
      X[t] = append([]int{X[t][n-1]}, X[t][:n-1]...)
    }
  }

	return X
}


func FinalNeighborhoodDistr(
	T int,
	d int,
  Qr probutil.RealTransition,
	Q probutil.NeighborTransition,
	nu probutil.InitDistr,
	k int,
	steps int,
	G graphutil.Graph,
  tlim int) probutil.PathDistr {

	f := func() fmt.Stringer {
		X := Realization(T, Qr, Q, nu, k, G, tlim)
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
  Qr probutil.RealTransition,
	Q probutil.NeighborTransition,
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
		X := Realization(T, Qr, Q, nu, k, G, tlim)
		return t_array, X.Col(0)
	}
	return probutil.GetTimeDistrSync(f, 1, float64(T), k, steps)
}

func PathDistr(
	T int,
  Qr probutil.RealTransition,
	Q probutil.NeighborTransition,
	nu probutil.InitDistr,
	k int,
	steps int,
	G graphutil.Graph,
  tlim int) probutil.PathDistr {
	f := func() fmt.Stringer {
		X := Realization(T, Qr, Q, nu, k, G, tlim)
		return X.Col(0)
	}
	return probutil.GetPathDistrSync(f, steps)
}
