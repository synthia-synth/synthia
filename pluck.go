package synthia

import (
	//	"math"
	"math/rand"
)

func whitenoise(samples int) []TimeDomain {
	sampleSpace := make([]TimeDomain, samples)
	for i := 0; i < samples; i++ {
		sampleSpace[i] = TimeDomain(rand.Float64())
	}
	return sampleSpace
}

type Plucker struct {
	samplerate float64
}

func NewPlucker(sampleRate float64) ToneSimulator {
	return &Plucker{samplerate: sampleRate}
}

func (p *Plucker) SetSampleRate(sampleRate float64) {
	p.samplerate = sampleRate
}

func (p *Plucker) pluck(freq, time float64) []TimeDomain {
	sampleSize := int(p.samplerate / freq)
	loops := int(freq * time)
	sound := whitenoise(sampleSize)
	lastIndex := sampleSize - 1
	out := make([]TimeDomain, loops*sampleSize)
	k := 0
	for i := 0; i < loops; i++ {
		for j, _ := range sound {
			if j == lastIndex {
				sound[j] = (0.996 * (sound[j] + sound[0]) / 2)
			} else {
				sound[j] = (sound[j] + sound[j+1]) / 2
			}
			out[k] = sound[j]
			k++
		}
	}
	return out
}

func (p *Plucker) Play(freq, seconds float64, vol int32) []TimeDomain {
	note := p.pluck(freq, seconds)
	for i, val := range note {
		note[i] = TimeDomain(vol) * val
	}
	return note
}
