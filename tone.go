package synthia

import (
	"github.com/synthia-synth/synthia/waveforms"
)

type ToneGenerator struct {
	sampleRate float64
	step       float64 //delta t
	wave       waveforms.Wave
}

func NewToneGenerator(sampleRate float64, wave waveforms.Wave) *ToneGenerator {
	if sampleRate < 1 {
		return nil
	}
	t := new(ToneGenerator)
	t.sampleRate = sampleRate
	t.step = 1. / sampleRate
	t.wave = wave
	return t
}

//Generates a wave
func (t *ToneGenerator) Tone(freq, seconds float64, vol int32) []TimeDomain {
	var synthArray = make([]TimeDomain, int(seconds*t.sampleRate))
	delta := freq * t.step

	for i := 0; i < len(synthArray); i++ {
		synthArray[i] = TimeDomain(t.wave(float64(i)*delta) * float64(vol))

	}
	return synthArray
}
