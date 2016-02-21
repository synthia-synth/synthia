package main
import(
	"math"
	"github.com/draringi/synthia/waveforms"
	)

type ToneGenerator struct {
	sampleRate float64
	step float64 //delta t
	wave waveforms.Wave
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

//Generates a square wave
func (t *ToneGenerator) Tone(freq, seconds float64, vol int32) []int32{
	var synthArray = make([]int32, int(seconds*t.sampleRate)) //duration/step = dur*sR
	var period, samplesPerPeriod = period(freq)

	for i:=0; i < len(synthArray); i++{
		synthArray[i]=int64(synthint32(wave(float64(i)/frameRate)*vol);
	}
	return synthArray
}

func period(freq float64) (float64, int){
	return 1.0/freq, int(period/t.step)
}
