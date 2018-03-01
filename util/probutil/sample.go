package probutil

func Sample(d []float64, r float64) int {
	s := 0.
	for k, v := range d {
		if s += v; s > r {
			return k
		}
	}
	return -1
}
