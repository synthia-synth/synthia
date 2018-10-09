package synthia

import (
	"github.com/synthia-synth/synthia/waveforms"
	"github.com/synthia-synth/synthia/domains"
)

type ToneGenerator interface {
	Play(freq, seconds float64, vol int32) []domains.Time
}

type ToneSimulator interface {
	ToneGenerator
	SetSampleRate(sampleRate float64)
}

type WaveToneGenerator struct {
	sampleRate float64
	step       float64 //delta t
	wave       waveforms.Wave
}

func NewWavetoneGenerator(sampleRate float64, wave waveforms.Wave) *WaveToneGenerator {
	if sampleRate < 1 {
		return nil
	}
	t := new(WaveToneGenerator)
	t.sampleRate = sampleRate
	t.step = 1. / sampleRate
	t.wave = wave
	return t
}

//Generates a wave
func (t *WaveToneGenerator) Play(freq, seconds float64, vol int32) []domains.Time {
	var synthArray = make([]domains.Time, int(seconds*t.sampleRate))
	delta := freq * t.step

	for i := 0; i < len(synthArray); i++ {
		synthArray[i] = domains.Time(t.wave(float64(i)*delta) * float64(vol))

	}
	return synthArray
}
