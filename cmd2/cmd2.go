package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mwortsma/particle_systems2/models/contact"
	"github.com/mwortsma/particle_systems2/models/potts"
	"github.com/mwortsma/particle_systems2/models/sir"
	"github.com/mwortsma/particle_systems2/util/probutil"
	"github.com/mwortsma/particle_systems2/util/graphutil"
	"io/ioutil"
	"strings"
)

func main() {

	// Get initial conditions.
	var nu string
	flag.StringVar(&nu, "nu", "[0.3,0.3,0.4]", "Initial Conditions, i.e. [0.3,0.7]")

	var graph string
	flag.StringVar(&graph, "graph", "ring", "Graph type: ring, dense, random")

	// File
	var file_str string
	flag.StringVar(&file_str, "file", "", "save location")

	// Params
	d := flag.Int("d", -1, "degree of a noe")
	n := flag.Int("n", 10, "number of nodes")
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
	tlim := flag.Int("tlim", 100, "for gibbs sampling")
	erp := flag.Float64("erp", 0.1, "random graph edge probability")
	seed := flag.Int("seed", 10, "for shared random seed")


	// Contact process
	contact_graph_path := flag.Bool("contact_graph_path", false, "")
	contact_graph_time := flag.Bool("contact_graph_time", false, "")
	contact_graph_end := flag.Bool("contact_graph_end", false, "")

	contact_tree_path := flag.Bool("contact_tree_path", false, "")
	contact_tree_time := flag.Bool("contact_tree_time", false, "")
	contact_tree_end := flag.Bool("contact_tree_end", false, "")

	contact_local_path := flag.Bool("contact_local_path", false, "")
	contact_local_time := flag.Bool("contact_local_time", false, "")
	contact_local_end := flag.Bool("contact_local_end", false, "")

	contact_meanfield_path := flag.Bool("contact_meanfield_path", false, "")
	contact_meanfield_time := flag.Bool("contact_meanfield_time", false, "")
	contact_meanfield_end := flag.Bool("contact_meanfield_end", false, "")

	contact_gs_path := flag.Bool("contact_gs_path", false, "")
	contact_gs_time := flag.Bool("contact_gs_time", false, "")
	contact_gs_end := flag.Bool("contact_gs_end", false, "")

	// SIR process
	sir_graph_path := flag.Bool("sir_graph_path", false, "")
	sir_graph_time := flag.Bool("sir_graph_time", false, "")
	sir_graph_end := flag.Bool("sir_graph_end", false, "")

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

	G := graphutil.GetGraph(graph,*n,*erp,*seed)

	var distr probutil.Distr

	switch {

	// Contact Process
	case *contact_graph_path:
		distr = contact.GraphPathDistr(*T, *p, *q, init, *steps, G)
	case *contact_graph_time:
		distr = contact.GraphTimeDistr(*T, *p, *q, init, *steps, G)
	case *contact_graph_end:
		distr = contact.GraphFinalNeighborhoodDistr(
			*T, *p, *q, init, *steps, *d, G)

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

	case *contact_gs_path:
		distr = contact.GSGraphPathDistr(*T, *p, *q, init, *steps, G, *tlim)
	case *contact_gs_time:
		distr = contact.GSGraphTimeDistr(*T, *p, *q, init, *steps, G, *tlim)
	case *contact_gs_end:
		distr = contact.GSGraphFinalNeighborhoodDistr(
			*T, *p, *q, init, *steps, *d, G, *tlim)

	// SIR Process
	case *sir_graph_path:
		distr = sir.GraphPathDistr(*T, *p, *q, init, *steps, G)
	case *sir_graph_time:
		distr = sir.GraphTimeDistr(*T, *p, *q, init, *steps, G)
	case *sir_graph_end:
		distr = sir.GraphFinalNeighborhoodDistr(*T, *p, *q, init, *steps, *d, G)

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
		distr = potts.MCMCFinalNeighborhoodDistr(
			*T, *k, *beta, *J, *h, init, *steps, *n, *d, G)
	case *potts_mcmc_time:
		distr = potts.MCMCTimeDistr(
			*T, *k, *beta, *J, *h, init, *steps, *n, *d, G)
	case *potts_mcmc_path:
		distr = potts.MCMCPathDistr(
			*T, *k, *beta, *J, *h, init, *steps, *n, *d, G)

	case *potts_gibbs_end:
		distr = potts.GibbsFinalNeighborhoodDistr(
			*beta, *J, *h, *k, G)
	case *potts_gibbs_time:
		distr = potts.GibbsTimeDistr(
			*T, *beta, *J, *h, *k, G)

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

	// fmt.Println(distr)

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
