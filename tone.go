package main
import(
	"github.com/draringi/synthia/waveforms"
	)

type ToneGenerator struct {
	sampleRate float64
	step float64 //delta t
	wave waveforms.Wave
}

func NewToneGenerator(sampleRate float64, wave waveforms.Wave) *ToneGenerator {
	if sampleRate < 1 {
		return nil
	}
	t := new(ToneGenerator)
	t.sampleRate = sampleRate
	t.step = 1/sampleRate
	t.wave = wave
	return t
}

//Generates a square wave
func (t *ToneGenerator) Tone(freq, seconds float64, vol int32) []int32{
	var synthArray = make([]int32, int(seconds*t.sampleRate))
	

	for i:=0; i < len(synthArray); i++ {
		synthArray[i]=int32(t.wave(float64(i)*t.step)*float64(vol));

	}
	return synthArray
}
