package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
)

func main() {

	// TODO: Clean

	dtpp_rec := flag.Bool("dtpp_rec", false, "rec")
	dtpp_mcmc_byt := flag.Bool("dtpp_mcmc_byt", false, "rec")
	dtpp_mcmc_end := flag.Bool("dtpp_mcmc_end", false, "rec")
	dtpp_rec_end := flag.Bool("dtpp_rec_end", false, "rec")
	dtpp_rec_full_end := flag.Bool("dtpp_rec_full_end", false, "rec")
	dtsir_rec := flag.Bool("dtsir_rec", false, "...")
	dtsir_rec_full := flag.Bool("dtsir_rec_full", false, "...")
	dtsir_rec_end := flag.Bool("dtsir_rec_end", false, "---")
	dtsir_full_continuous := flag.Bool("dtsir_full_continuous", false, "---")
	dtsir_full_tree := flag.Bool("dtsir_full_tree", false, "---")
	dtsir_full_end := flag.Bool("dtsir_full_end", false, "---")

	d := flag.Int("d", 2, "degree of a noe")
	n := flag.Int("n", 10, "number of nodes")
	k := flag.Int("k", 3, "states")
	T := flag.Int("T", 2, "time horizon. T>0")
	beta := flag.Float64("beta", 1.5, "temp inverse")
	tau := flag.Int("tau", -1, "how many steps to look back")
	p := flag.Float64("p", 2.0/3.0, "infection rate")
	q := flag.Float64("q", 1.0/3.0, "recovery rate")

	steps := flag.Int("steps", 2, "steps")

	init := []float64{0.8, 0.2, 0}

	var file_str string
	flag.StringVar(&file_str, "file", "", "where to save the distribution.")

	flag.Parse()

	var distr probutil.GenDistr

	switch {

	case *dtpp_rec:
		distr = dtpp_local.Run(*T, *tau, *d, *beta, *k, *n)

	case *dtpp_mcmc_byt:
		distr = dtpp_local.MCMC_byt(*T, *tau, *d, *beta, *k, *n, *steps)
	case *dtpp_mcmc_end:
		distr = dtpp_local.MCMC_end(*T, *tau, *d, *beta, *k, *n, *steps)

	case *dtpp_rec_end:
		distr = dtpp_local.EndRun(*T, *tau, *d, *beta, *k, *n)
	case *dtpp_rec_full_end:
		distr = dtpp_local.FullEndRun(*T, *tau, *d, *beta, *k, *n)

	case *dtsir_full_tree:
		distr = dtsir_full.RegTreeTypicalDistr(*T, *d, *p, *q, init, *steps)

	case *dtsir_full_continuous:
		distr = dtsir_full.RegTreeTDistr(*T, *d, *p, *q, init, *steps)

	case *dtsir_full_end:
		distr = dtsir_full.RegTreeEndDistr(*T, *d, *p, *q, init, *steps)

	case *dtsir_rec_full:
		distr = dtsir_local.FullRun(*T, *tau, *d, *p, *q, init)

	case *dtsir_rec:
		distr = dtsir_local.Run(*T, *tau, *d, *p, *q, init)

	case *dtsir_rec_end:
		distr = dtsir_local.EndRun(*T, *tau, *d, *p, *q, init)

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
