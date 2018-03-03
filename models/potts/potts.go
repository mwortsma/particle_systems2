package potts

import (
	"github.com/mwortsma/particle_systems2/full/gibbs"
	"github.com/mwortsma/particle_systems2/full/mcmc"
	"github.com/mwortsma/particle_systems2/local"
	"github.com/mwortsma/particle_systems2/meanfield"
	"github.com/mwortsma/particle_systems2/util/graphutil"
	"github.com/mwortsma/particle_systems2/util/matutil"
	"github.com/mwortsma/particle_systems2/util/probutil"
	"golang.org/x/exp/rand"
	"math"
	//"fmt"
)

// General.
func getLawQ(n, d, k int, beta, J, h float64) probutil.LawTransition {
	// transition from s to s_new
	return func(s_new, s int, f probutil.Law) float64 {
		if s_new != s {
			diff_h := math.Max(
				float64(d)*float64(J)*(f[s]-f[s_new])+float64(h)*float64(s_new-s), 0.0)
			prob_transition := 1.0 / (float64(n) * float64(k-1))
			return math.Exp(-beta*diff_h) * prob_transition
		} else {
			prob := 1.0
			for i := 0; i < k; i++ {
				if i != s {
					diff_h := math.Max(
						float64(d)*float64(J)*(f[s]-f[i])+float64(h)*float64(i-s), 0.0)
					prob_transition := 1.0 / (float64(n) * float64(k-1))
					prob -= math.Exp(-beta*diff_h) * prob_transition
				}
			}
			return prob
		}
	}
}


func getNeighborQ(n, d, k int, beta, J, h float64) probutil.NeighborTransition {
	return probutil.GetNeighborTransition(getLawQ(n, d, k, beta, J, h), k)
}

/*
// A test.
func getNeighborQ(n, d, k int, beta, J, h float64) probutil.NeighborTransition {
	Ham := func(r, n int, v matutil.Vec) float64 {
		match := 0
		sum := r
		for _, val := range v {
			if val == r {
				match++
			}
			sum += val
		}
		return float64(sum) - float64(match)
	}



	Q_help := func(r int, s int, n int, v matutil.Vec) float64 {
		prob := 1.0/float64(n*(k-1))
		ham_r := Ham(r,n,v)
		ham_s := Ham(s,n,v)
		if ham_r > ham_s {
			prob *= math.Exp(-beta * (ham_r - ham_s))
		}
		return prob
	}

	// prob transition from s to r
	Q := func(r int, s int, v matutil.Vec) float64 {
		if r != s {
			return Q_help(r,s,n,v)
		} else {
			prob := 1.0
			for i := 0; i < k; i++ {
				if i != s {
					prob -= Q_help(i,s,n,v)
				}
			}
			return prob
		}
	}
	return Q
}
*/


func getNewStateFunc(n, d, k int, beta, J, h float64) mcmc.NewStateFunc {
	return func(v matutil.Vec, r *rand.Rand) ([]int, []int) {
		site := rand.Intn(n)
		newvalue := rand.Intn(k)
		for newvalue == v[site] {
			newvalue = rand.Intn(k)
		}
		return []int{site}, []int{newvalue}
	}
}


func H(sigma matutil.Vec, J, h float64, G graphutil.Graph) float64 {
	/*
	sum := 0.0
	for i := range G {
		sum += h * float64(sigma[i])
		for _, j := range G[i] {
			if sigma[i] != sigma[j] {
				sum += J
			}
		}
	}
	return sum
	*/
	match := 0
	n := len(G)
	sum := sigma[0]
	if sigma[0] == sigma[n-1] {
		match++
	}
	for i := 1; i < n; i++ {
		sum += sigma[i]
		if sigma[i-1] == sigma[i] {
			match++
		}
	}
	return float64(sum) - float64(match)
}

func getP(beta, J, h float64, G graphutil.Graph) func(matutil.Vec) float64 {
	return func(sigma matutil.Vec) float64 {
		return math.Exp(-beta * H(sigma, J, h, G))
	}
}

func getTransitionProbFunc(
	n, d, k int,
	beta, J, h float64,
	G graphutil.Graph) mcmc.TransitionProbFunc {
	return func(
		sigma matutil.Vec,
		sites, newvals []int) float64 {
		tau := make([]int, n)
		copy(tau, sigma)
		for j, site := range sites {
			tau[site] = newvals[j]
		}
		return math.Exp(-beta * (H(tau, J, h, G) - H(sigma, J, h, G)))
	}
}

// MCMC -- Ring.
func MCMCRingFinalNeighborhoodDistr(
	T, k int,
	beta, J, h float64,
	nu probutil.InitDistr,
	steps int,
	n int,
	d int) probutil.PathDistr {

	return mcmc.FinalNeighborhoodDistr(
		T, getNewStateFunc(n, d, k, beta, J, h),
		getTransitionProbFunc(n, d, k, beta, J, h, graphutil.Ring(n)),
		nu, graphutil.Ring(n), n, steps, d)
}

func MCMCRingTimeDistr(
	T, k int,
	beta, J, h float64,
	nu probutil.InitDistr,
	steps int,
	n int,
	d int) probutil.TimeDistr {

	return mcmc.TimeDistr(
		T, getNewStateFunc(n, d, k, beta, J, h),
		getTransitionProbFunc(n, d, k, beta, J, h, graphutil.Ring(n)),
		nu, steps, n, k)
}

func MCMCRingPathDistr(
	T, k int,
	beta, J, h float64,
	nu probutil.InitDistr,
	steps int,
	n int,
	d int) probutil.PathDistr {

	return mcmc.PathDistr(
		T, getNewStateFunc(n, d, k, beta, J, h),
		getTransitionProbFunc(n, d, k, beta, J, h, graphutil.Ring(n)),
		nu, steps, n)
}

// Gibbs -- Ring.
func GibbsRingTimeDistr(
	T int,
	beta, J, h float64,
	n, k int) probutil.TimeDistr {
	return gibbs.TimeDistr(
		T, getP(beta, J, h, graphutil.Ring(n)),
		graphutil.Ring(n), k)
}

func GibbsRingFinalNeighborhoodDistr(
	beta, J, h float64,
	n, k int) probutil.PathDistr {
	return gibbs.FinalNeighborhoodDistr(
		getP(beta, J, h, graphutil.Ring(n)),
		graphutil.Ring(n), k)
}

// Local
func LocalFinalNeighborhoodDistr(
	T, tau, d, k, n int,
	beta, J, h float64,
	nu probutil.InitFunc) probutil.PathDistr {

	return local.FinalNeighborhoodDistr(
		T, tau, d,
		getNeighborQ(n, d, k, beta, J, h),
		nu, k)
}

func LocalTimeDistr(
	T, tau, d, k, n int,
	beta, J, h float64,
	nu probutil.InitFunc) probutil.TimeDistr {

	return local.TimeDistr(
		T, tau, d,
		getNeighborQ(n, d, k, beta, J, h),
		nu, k)
}

func LocalPathDistr(
	T, tau, d, k, n int,
	beta, J, h float64,
	nu probutil.InitFunc) probutil.PathDistr {

	return local.PathDistr(
		T, tau, d,
		getNeighborQ(n, d, k, beta, J, h),
		nu, k)
}

// Mean Field
func MeanFieldFinalNeighborhoodDistr(
	T, d, k, n int,
	beta, J, h float64,
	nu probutil.InitDistr) probutil.PathDistr {

	return meanfield.FinalNeighborhoodDistr(
		T, getLawQ(n, d, k, beta, J, h),
		nu, k, d)
}

func MeanFieldTimeDistr(
	T, d, k, n int,
	beta, J, h float64,
	nu probutil.InitDistr) probutil.TimeDistr {

	return meanfield.TimeDistr(
		T, getLawQ(n, d, k, beta, J, h),
		nu, k)
}

func MeanFieldPathDistr(
	T, d, k, n int,
	beta, J, h float64,
	nu probutil.InitDistr) probutil.PathDistr {

	return meanfield.PathDistr(
		T, getLawQ(n, d, k, beta, J, h),
		nu, k)
}
