package main

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

func (t *ToneGenerator) Tone(freq, seconds float64) []int32{
	return nil
}
