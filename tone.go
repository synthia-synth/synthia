package main
import("fmt", "math")


type ToneGenerator struct {
	sampleRate float64
	step float64
}

func NewToneGenerator(sampleRate float64) *Tone {
	if sampleRate < 1 {
		return nil
	}
	t := new(ToneGenerator)
	t.sampleRate = sampleRate
	t.step = 1/sampleRate
	return t
}

//Max's function - Given freq and delta t return an array for each step
//Call Tone with freq and sec (floats) - return array
func (t *ToneGenerator) Tone(freq, seconds float64) []int32{
	var synthArray = make([]int32, int(duration*sampleRate)) //duration/step = dur*sR
	for i:=0; i < len(synthArray); i++{
		synthArray[i] = sin(freq * math.PI * i * step)
	}
	return synthArray
}
