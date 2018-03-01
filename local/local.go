package local

import (
	"fmt"
	"github.com/mwortsma/particle_systems2/mathutil"
	"github.com/mwortsma/particle_systems2/matutil"
	"github.com/mwortsma/particle_systems2/probutil"
	"math"
)

/////////////////////////////////////////////////////////////////////
/////									The Main Algorithm										    /////
/////////////////////////////////////////////////////////////////////

func getJointAndTransitionKernel(
	t int,
	tau int,
	d int,
	Q probutil.NeighborTransition,
	j probutil.PathDistr,
	c probutil.Conditional,
	k int) (probutil.PathDistr, probutil.PathDistr) {

	fmt.Println("Obtaining joint at", t)

	jnew := make(probutil.PathDistr)

	p := make(probutil.PathDistr)

	l := mathutil.Min(tau, t) + 1
	r := mathutil.Min(tau, t-1) + 1

	prev_vals := mathutil.QMats(r, d+1, k)
	new_vals := mathutil.QStrings(d+1, k)
	other_children_vals := mathutil.QStrings(d-1, k)

	for _, prev := range prev_vals {
		prob_prev := j[prev.String()]
		for _, new_val := range new_vals {
			full := append(prev, new_val)
			trimmed := full[len(full)-l:]
			trimmed_str := trimmed.String()
			lastrow := prev[len(prev)-1]
			prob := 1.0
			prob *= Q(new_val[0], lastrow[0], lastrow[1:])
			for i := 1; i < d+1; i++ {
				sum_prob := 0.0
				for _, other_children := range other_children_vals {
					hist := prev.Cols([]int{i, 0}).String()
					sum_prob += c[hist][other_children.String()] *
						Q(new_val[i], lastrow[i], append(other_children, lastrow[0]))
				}
				prob *= sum_prob
			}
			p[full.String()] = prob
			prob *= prob_prev
			if _, ok := jnew[trimmed_str]; !ok {
				jnew[trimmed_str] = prob
			} else {
				jnew[trimmed_str] += prob
			}
		}
	}

	return jnew, p
}

func getConditional(
	t int,
	tau int,
	d int,
	jt probutil.PathDistr,
	k int) probutil.Conditional {
	fmt.Println("Obtaining Conditional at", t)

	ct := make(probutil.Conditional)

	l := mathutil.Min(tau, t) + 1

	history_vals := mathutil.QMats(l, 2, k)
	children_vals := mathutil.QMats(l, d-1, k)

	for _, history := range history_vals {
		hist_str := history.String()
		ct[hist_str] = make(probutil.PathDistr)
		denom := 0.0
		for _, children := range children_vals {
			full := matutil.Concat(history, children)
			denom += jt[full.String()]
		}
		// Important
		if denom == 0 {
			continue
		}
		for _, children := range children_vals {
			lastrow := matutil.Vec(children[l-1]).String()
			// TODO: debug
			if _, ok := ct[hist_str][lastrow]; !ok {
				ct[hist_str][lastrow] = 0
			}
			full := matutil.Concat(history, children)
			ct[hist_str][lastrow] += jt[full.String()] / denom
		}
	}

	return ct
}

/////////////////////////////////////////////////////////////////////
/////									  Shortened & Scalable								    /////
/////////////////////////////////////////////////////////////////////

// Gets the distriution over the local neighborhood at the end
func FinalNeighborhoodDistr(
	T int,
	tau int,
	d int,
	Q probutil.NeighborTransition,
	nu probutil.InitialConditions,
	k int) probutil.PathDistr {

	j := tauApprox(T, tau, d, Q, nu, k)

	f := make(probutil.PathDistr)

	for k, prob := range j {
		mat := matutil.StringToMat(k)
		v := matutil.Vec(mat[len(mat)-1])
		str := v.String()
		if _, ok := f[str]; !ok {
			f[str] = prob
		} else {
			f[str] += prob
		}
	}

	return f
}

// Returns the last final joint distribution.
func tauApprox(
	T int,
	tau int,
	d int,
	Q probutil.NeighborTransition,
	nu probutil.InitialConditions,
	k int) probutil.PathDistr {

	if tau < 0 {
		tau = math.MaxInt32
	}

	j := make(probutil.PathDistr)

	var c probutil.Conditional

	for _, init := range mathutil.QStrings(d+1, k) {
		j[matutil.Mat([][]int{init}).String()] = nu(init)
	}

	for t := 1; t < T; t++ {

		c = getConditional(t-1, tau, d, j, k)

		j, _ = getJointAndTransitionKernel(t, tau, d, Q, j, c, k)

	}

	return j
}

/////////////////////////////////////////////////////////////////////
/////									     For Every Time				    				    /////
/////////////////////////////////////////////////////////////////////

func TimeDistr(
	T int,
	tau int,
	d int,
	Q probutil.NeighborTransition,
	nu probutil.InitialConditions,
	k int) probutil.TimeDistr {

	js, _ := tauApproxForEachT(T, tau, d, Q, nu, k)

	f := make([][]float64, T)
	t := 0
	for _, j := range js {
		f[t] = make([]float64, k)
		for k, prob := range j {
			mat := matutil.StringToMat(k)
			v := matutil.Vec(mat.Col(0))
			f[t][v[len(v)-1]] += prob
		}
		t += 1
	}

	return probutil.TimeDistr{1, float64(T), k, f}
}

// This is equivelant to above except for that returns j
// for all t and c for all t
func tauApproxForEachT(
	T int,
	tau int,
	d int,
	Q probutil.NeighborTransition,
	nu probutil.InitialConditions,
	k int) ([]probutil.PathDistr, []probutil.PathDistr) {

	if tau < 0 {
		tau = math.MaxInt32
	}

	j := make([]probutil.PathDistr, T)
	p := make([]probutil.PathDistr, T)
	j[0] = make(probutil.PathDistr)

	var c probutil.Conditional

	for _, init := range mathutil.QStrings(d+1, k) {
		j[0][matutil.Mat([][]int{init}).String()] = nu(init)
	}

	for t := 1; t < T; t++ {

		j[t] = make(probutil.PathDistr)
		p[t] = make(probutil.PathDistr)

		c = getConditional(t-1, tau, d, j[t-1], k)

		j[t], p[t] = getJointAndTransitionKernel(t, tau, d, Q, j[t-1], c, k)

	}

	fmt.Println("Exiting")

	return j, p
}

/////////////////////////////////////////////////////////////////////
/////									     For Path Distr				    				    /////
/////////////////////////////////////////////////////////////////////

func fullPathProbability(
	T int,
	tau int,
	d int,
	j []probutil.PathDistr,
	p []probutil.PathDistr,
	state matutil.Mat) float64 {

	prob := 1.0
	t := T - 1
	for len(state) > tau+1 {
		rel_state := state[len(state)-(tau+2):]
		prob *= p[t][rel_state.String()]
		t = t - 1
		state = state[0 : len(state)-1]
	}
	prob *= j[t][state.String()]
	return prob

}

func PathDistr(
	T int,
	tau int,
	d int,
	Q probutil.NeighborTransition,
	nu probutil.InitialConditions,
	k int) probutil.PathDistr {

	j_array, p_array := tauApproxForEachT(T, tau, d, Q, nu, k)

	f := make(probutil.PathDistr)

	if tau < 0 {
		tau = math.MaxInt32
	}

	states := mathutil.QMats(T, d+1, k)
	for _, state := range states {
		path := state.Col(0).String()
		if _, ok := f[path]; !ok {
			f[path] = 0.0
		}
		f[path] += fullPathProbability(T, tau, d, j_array, p_array, state)
	}

	return f
}
