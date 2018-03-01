package probutil

import (
	"github.com/mwortsma/particle_systems2/matutil"
)

type InitialConditions func(matutil.Vec) float64

type NeighborTransition func(int, int, matutil.Vec) float64

// TODO: type LawTransition func(int, int, matutil.Vec) float64
