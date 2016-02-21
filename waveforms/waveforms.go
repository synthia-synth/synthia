package waveforms
import ("math")

type Wave func(float64) float64

func Sin(x float64) float64 {
	xMod:= math.Remainder(x,1)
	return math.Sin(math.Pi*2*xMod)
}

func Tri(x float64) float64 {
	xMod:= math.Remainder(x,1)
	if xMod<0.5 {
  	return 1. - xMod*4.
  }
  else {
  	return (xMod*4. - 3.)
  }
}
