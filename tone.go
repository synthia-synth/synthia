package main
import("math")


type ToneGenerator struct {
	sampleRate float64
	step float64
}

func NewToneGenerator(sampleRate float64) *ToneGenerator {
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
func (t *ToneGenerator) Tone(freq, seconds float64, vol int32) []int32{
	var synthArray = make([]int32, int(seconds*t.sampleRate)) //duration/step = dur*sR
	for i:=0; i < len(synthArray); i++{
		synthArray[i] = int32(float64(vol)*math.Sin(freq *2* math.Pi * float64(i) * t.step))
	}
	return synthArray
}
