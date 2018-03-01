package mathutil

import (
	"github.com/mwortsma/particle_systems2/matutil"
	"math"
	"strconv"
)

// Minimum
func Min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
	return -1
}

// Binary Strings
func BinaryStrings(n int) []matutil.Vec {
	s := make([]matutil.Vec, int(math.Pow(2.0, float64(n))))
	for i := 0; i < int(math.Pow(2.0, float64(n))); i++ {
		b := []byte(string(strconv.FormatInt(int64(i), 2)))
		a := make(matutil.Vec, n)
		for j := 0; j < len(b); j++ {
			a[j] = int(b[len(b)-j-1] - 48)
		}
		s[i] = a
	}
	return s

}

// BinaryMats
func BinaryMats(r, c int) []matutil.Mat {
	l := int(math.Pow(2.0, float64(r*c)))
	s := make([]matutil.Mat, l)
	strings := BinaryStrings(r * c)
	for k, str := range strings {
		s[k] = matutil.Create(r, c)
		for i := 0; i < r; i++ {
			for j := 0; j < c; j++ {
				s[k][i][j] = str[c*i+j]
			}
		}
	}
	return s
}

// QStrings
func QStrings(n int, q int) []matutil.Vec {
	s := make([]matutil.Vec, int(math.Pow(float64(q), float64(n))))
	for i := 0; i < int(math.Pow(float64(q), float64(n))); i++ {
		b := []byte(string(strconv.FormatInt(int64(i), q)))
		a := make(matutil.Vec, n)
		for j := 0; j < len(b); j++ {
			a[j] = int(b[len(b)-j-1] - 48)
		}
		s[i] = a
	}
	return s
}

// Q Mats
func QMats(r, c int, q int) []matutil.Mat {
	l := int(math.Pow(float64(q), float64(r*c)))
	s := make([]matutil.Mat, l)
	strings := QStrings(r*c, q)
	for k, str := range strings {
		s[k] = matutil.Create(r, c)
		for i := 0; i < r; i++ {
			for j := 0; j < c; j++ {
				s[k][i][j] = str[c*i+j]
			}
		}
	}
	return s
}
