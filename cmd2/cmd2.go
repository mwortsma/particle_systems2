package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mwortsma/particle_systems2/models/contact"
	"github.com/mwortsma/particle_systems2/models/potts"
	"github.com/mwortsma/particle_systems2/models/sir"
	"github.com/mwortsma/particle_systems2/util/probutil"
	"io/ioutil"
	"strings"
)

func main() {

	// Get initial conditions.
	var nu string
	flag.StringVar(&nu, "nu", "[0.3,0.7]", "Initial Conditions, i.e. [0.3,0.7]")

	// File
	var file_str string
	flag.StringVar(&file_str, "file", "", "save location")

	// Params
	d := flag.Int("d", -1, "degree of a noe")
	n := flag.Int("n", -1, "number of nodes")
	k := flag.Int("k", 3, "states")
	T := flag.Int("T", -1, "time horizon. T>0")
	depth := flag.Int("depth", -1, "time horizon. T>0")
	beta := flag.Float64("beta", 1.5, "temp inverse")
	J := flag.Float64("J", -1.0, "potts param")
	h := flag.Float64("h", 1.0, "potts param")
	tau := flag.Int("tau", -1, "how many steps to look back")
	p := flag.Float64("p", 2.0/3.0, "infection rate")
	q := flag.Float64("q", 1.0/3.0, "recovery rate")
	steps := flag.Int("steps", -1, "for estimating probability")

	// Contact process
	contact_dense_path := flag.Bool("contact_dense_path", false, "")
	contact_dense_time := flag.Bool("contact_dense_time", false, "")
	contact_dense_end := flag.Bool("contact_dense_end", false, "")

	contact_tree_path := flag.Bool("contact_tree_path", false, "")
	contact_tree_time := flag.Bool("contact_tree_time", false, "")
	contact_tree_end := flag.Bool("contact_tree_end", false, "")

	contact_local_path := flag.Bool("contact_local_path", false, "")
	contact_local_time := flag.Bool("contact_local_time", false, "")
	contact_local_end := flag.Bool("contact_local_end", false, "")

	contact_meanfield_path := flag.Bool("contact_meanfield_path", false, "")
	contact_meanfield_time := flag.Bool("contact_meanfield_time", false, "")
	contact_meanfield_end := flag.Bool("contact_meanfield_end", false, "")

	// SIR process
	sir_dense_path := flag.Bool("sir_dense_path", false, "")
	sir_dense_time := flag.Bool("sir_dense_time", false, "")
	sir_dense_end := flag.Bool("sir_dense_end", false, "")

	sir_tree_path := flag.Bool("sir_tree_path", false, "")
	sir_tree_time := flag.Bool("sir_tree_time", false, "")
	sir_tree_end := flag.Bool("sir_tree_end", false, "")

	sir_local_path := flag.Bool("sir_local_path", false, "")
	sir_local_time := flag.Bool("sir_local_time", false, "")
	sir_local_end := flag.Bool("sir_local_end", false, "")

	sir_meanfield_path := flag.Bool("sir_meanfield_path", false, "")
	sir_meanfield_time := flag.Bool("sir_meanfield_time", false, "")
	sir_meanfield_end := flag.Bool("sir_meanfield_end", false, "")

	// Potts process
	potts_mcmc_path := flag.Bool("potts_mcmc_path", false, "")
	potts_mcmc_time := flag.Bool("potts_mcmc_time", false, "")
	potts_mcmc_end := flag.Bool("potts_mcmc_end", false, "")

	potts_gibbs_end := flag.Bool("potts_gibbs_end", false, "")
	potts_gibbs_time := flag.Bool("potts_gibbs_time", false, "")

	potts_local_path := flag.Bool("potts_local_path", false, "")
	potts_local_time := flag.Bool("potts_local_time", false, "")
	potts_local_end := flag.Bool("potts_local_end", false, "")

	potts_meanfield_path := flag.Bool("potts_meanfield_path", false, "")
	potts_meanfield_time := flag.Bool("potts_meanfield_time", false, "")
	potts_meanfield_end := flag.Bool("potts_meanfield_end", false, "")

	flag.Parse()

	var init []float64
	dec := json.NewDecoder(strings.NewReader(nu))
	err := dec.Decode(&init)
	if err != nil {
		fmt.Println("Nu not formatted correctly, e.g. [0.3, 0.7]")
		return
	}
	init_f := probutil.GetInitFunc(init)

	var distr probutil.Distr

	switch {

	// Contact Process
	case *contact_dense_path:
		distr = contact.DensePathDistr(*T, *p, *q, init, *steps, *n)
	case *contact_dense_time:
		distr = contact.DenseTimeDistr(*T, *p, *q, init, *steps, *n)
	case *contact_dense_end:
		distr = contact.DenseFinalNeighborhoodDistr(
			*T, *p, *q, init, *steps, *n, *d)

	case *contact_tree_path:
		distr = contact.TreePathDistr(*T, *p, *q, *d, init, *steps, *depth)
	case *contact_tree_time:
		distr = contact.TreeTimeDistr(*T, *p, *q, *d, init, *steps, *depth)
	case *contact_tree_end:
		distr = contact.TreeFinalNeighborhoodDistr(
			*T, *p, *q, *d, init, *steps, *depth)

	case *contact_local_path:
		distr = contact.LocalPathDistr(*T, *tau, *d, *p, *q, init_f)
	case *contact_local_time:
		distr = contact.LocalTimeDistr(*T, *tau, *d, *p, *q, init_f)
	case *contact_local_end:
		distr = contact.LocalFinalNeighborhoodDistr(*T, *tau, *d, *p, *q, init_f)

	case *contact_meanfield_path:
		distr = contact.MeanFieldPathDistr(*T, *p, *q, init)
	case *contact_meanfield_time:
		distr = contact.MeanFieldTimeDistr(*T, *p, *q, init)
	case *contact_meanfield_end:
		distr = contact.MeanFieldFinalNeighborhoodDistr(*T, *p, *q, init, *d)

	// SIR Process
	case *sir_dense_path:
		distr = sir.DensePathDistr(*T, *p, *q, init, *steps, *n)
	case *sir_dense_time:
		distr = sir.DenseTimeDistr(*T, *p, *q, init, *steps, *n)
	case *sir_dense_end:
		distr = sir.DenseFinalNeighborhoodDistr(*T, *p, *q, init, *steps, *n, *d)

	case *sir_tree_path:
		distr = sir.TreePathDistr(*T, *p, *q, *d, init, *steps, *depth)
	case *sir_tree_time:
		distr = sir.TreeTimeDistr(*T, *p, *q, *d, init, *steps, *depth)
	case *sir_tree_end:
		distr = sir.TreeFinalNeighborhoodDistr(*T, *p, *q, *d, init, *steps, *depth)

	case *sir_local_path:
		distr = sir.LocalPathDistr(*T, *tau, *d, *p, *q, init_f)
	case *sir_local_time:
		distr = sir.LocalTimeDistr(*T, *tau, *d, *p, *q, init_f)
	case *sir_local_end:
		distr = sir.LocalFinalNeighborhoodDistr(*T, *tau, *d, *p, *q, init_f)

	case *sir_meanfield_path:
		distr = sir.MeanFieldPathDistr(*T, *p, *q, init)
	case *sir_meanfield_time:
		distr = sir.MeanFieldTimeDistr(*T, *p, *q, init)
	case *sir_meanfield_end:
		distr = sir.MeanFieldFinalNeighborhoodDistr(*T, *p, *q, init, *d)

		// Potts Process
	case *potts_mcmc_end:
		distr = potts.MCMCRingFinalNeighborhoodDistr(
			*T, *k, *beta, *J, *h, init, *steps, *n, *d)
	case *potts_mcmc_time:
		distr = potts.MCMCRingTimeDistr(
			*T, *k, *beta, *J, *h, init, *steps, *n, *d)
	case *potts_mcmc_path:
		distr = potts.MCMCRingPathDistr(
			*T, *k, *beta, *J, *h, init, *steps, *n, *d)

	case *potts_gibbs_end:
		distr = potts.GibbsRingFinalNeighborhoodDistr(
			*beta, *J, *h, *n, *k)
	case *potts_gibbs_time:
		distr = potts.GibbsRingTimeDistr(
			*T, *beta, *J, *h, *n, *k)

	case *potts_local_path:
		distr = potts.LocalPathDistr(*T, *tau, *d, *k, *n, *beta, *J, *h, init_f)
	case *potts_local_time:
		distr = potts.LocalTimeDistr(*T, *tau, *d, *k, *n, *beta, *J, *h, init_f)
	case *potts_local_end:
		distr = potts.LocalFinalNeighborhoodDistr(
			*T, *tau, *d, *k, *n, *beta, *J, *h, init_f)

	case *potts_meanfield_path:
		distr = potts.MeanFieldPathDistr(*T, *d, *k, *n, *beta, *J, *h, init)
	case *potts_meanfield_time:
		distr = potts.MeanFieldTimeDistr(*T, *d, *k, *n, *beta, *J, *h, init)
	case *potts_meanfield_end:
		distr = potts.MeanFieldFinalNeighborhoodDistr(*T, *d, *k, *n, *beta, *J, *h, init)
	}

	fmt.Println(distr)

	b, err := json.Marshal(distr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Writing to file ...")

	if file_str != "" {
		err = ioutil.WriteFile(file_str, b, 0777)
		if err != nil {
			panic(err)
		}
	}
}
