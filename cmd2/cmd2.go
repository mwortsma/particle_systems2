package main

import (
	"encoding/json"
	"flag"
	"fmt"
  "strings"
	"io/ioutil"
	"github.com/mwortsma/particle_systems2/util/probutil"
	"github.com/mwortsma/particle_systems2/models/contact"


)

func main() {

  // Get initial conditions.
  var nu string
  flag.StringVar(&nu, "nu", "", "Initial Conditions, i.e. [0.3, 0.7]")

  var init []float64
  dec := json.NewDecoder(strings.NewReader(nu))
  err := dec.Decode(&init)
  if err != nil {
    fmt.Println("Nu not formatted correctly, e.g. [0.3, 0.7]")
    return
  }
	init_f := probutil.GetInitFunc(init)

	// File
	var file_str string
	flag.StringVar(&file_str, "file", "", "save location")

	// Params
	d := flag.Int("d", -1, "degree of a noe")
	n := flag.Int("n", -1, "number of nodes")
	// k := flag.Int("k", -1, "states")
	T := flag.Int("T", -1, "time horizon. T>0")
	depth := flag.Int("T", -1, "time horizon. T>0")
	// beta := flag.Float64("beta", 1.5, "temp inverse")
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

  flag.Parse()
	var distr probutil.Distr

	switch {

	// Contact Process
	case *contact_dense_path:
		distr = contact.DensePathDistr(*T,*p,*q,init,*steps,*n)
	case *contact_dense_time:
		distr = contact.DenseTimeDistr(*T,*p,*q,init,*steps,*n)
	case *contact_dense_end:
		distr = contact.DenseFinalNeighborhoodDistr(*T,*p,*q,init,*steps,*n,*d)

	case *contact_tree_path:
		distr = contact.TreePathDistr(*T,*p,*q,*d,init,*steps,*depth)
	case *contact_tree_time:
		distr = contact.TreeTimeDistr(*T,*p,*q,*d,init,*steps,*depth)
	case *contact_tree_end:
		distr = contact.TreeFinalNeighborhoodDistr(*T,*p,*q,*d,init,*steps,*depth)

	case *contact_local_path:
		distr = contact.LocalPathDistr(*T,*tau,*d,*p,*q,init_f)
	case *contact_local_time:
		distr = contact.LocalTimeDistr(*T,*tau,*d,*p,*q,init_f)
	case *contact_local_end:
		distr = contact.LocalFinalNeighborhoodDistr(*T,*tau,*d,*p,*q,init_f)
	}

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
