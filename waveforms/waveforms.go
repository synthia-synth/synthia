package waveforms

import (
	"math"
)

type Wave func(float64) float64

func Null(x float64) float64 {
	return 0.0
}

func Sin(x float64) float64 {
	xMod := math.Remainder(x, 1) + 0.5
	return math.Sin(math.Pi * 2 * xMod)
}

func Tri(x float64) float64 {
	xMod := math.Remainder(x, 1) + 0.5
	if xMod < 0.5 {
		return 1. - xMod*4.
	} else {
		return (xMod*4. - 3.)
	}
}

func Sqr(x float64) float64 {
	xMod := math.Remainder(x, 1) + 0.5
	//fmt.Printf("%v ", xMod)
	if xMod < 0.45 {
		return 1
	} else {
		return -1
	}
}

func Saw(x float64) float64 {
	return math.Remainder((x-0.5)*2, 2)
}
