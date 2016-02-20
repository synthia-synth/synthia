package main
import("math")


type ToneGenerator struct {
	sampleRate float64
	step float64 //delta t
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

func (t *ToneGenerator) Tone(freq, seconds float64, vol int32) []int32{
	return t.SinTone(freq, seconds, vol)
}

//Generates a sin wave of a certain freq for a specific volume
func (t *ToneGenerator) SinTone(freq, seconds float64, vol int32) []int32{
	var synthArray = make([]int32, int(seconds*t.sampleRate)) //duration/step = dur*sR
	for i:=0; i < len(synthArray); i++{
		synthArray[i] = int32(float64(vol)*math.Sin(freq *2* math.Pi * float64(i) * t.step))
	}
	return synthArray
}

//Generates a square wave
func (t *ToneGenerator) SquareTone(freq, seconds float64, vol int32) []int32{
	var synthArray = make([]int32, int(seconds*t.sampleRate)) //duration/step = dur*sR
	var period, samplesPerPeriod = period(freq)

	for i:=0; i < len(synthArray); i++{
		if i%samplesPerPeriod < samplesPerPeriod/2 { //Width: first half = 1, second = 0
			synthArray[i] = vol
		} else {
			synthArray[i] = -vol
		}
	}
	return synthArray
}

//Generates a saw wave
func (t *ToneGenerator) SawTone(freq, seconds float64, vol int32) []int32{
	var synthArray = make([]int32, int(seconds*t.sampleRate)) //duration/step = dur*sR
	var period, samplesPerPeriod = period(freq)

	//from -vol to vol linearly for one period
	var ramp = 2*vol/samplesPerPeriod
	var val = -vol
	for i:=0; i < len(synthArray); i++{
		if val < samplesPerPeriod{
			val = val + ramp
		} else{ //Reset
			val = -vol
		}
		synthArray[i] = val
	}
	return synthArray
}

//Generates a triangle wave
func (t *ToneGenerator) TriTone(freq, seconds float64, vol int32) []int32{
	var synthArray = make([]int32, int(seconds*t.sampleRate)) //duration/step = dur*sR
	var period, samplesPerPeriod = period(freq)


	for i:=0; i < len(synthArray); i++{ //Rise then fall

	}
	return synthArray
}


func period(freq float64) (float64, int){
	return 1.0/freq, int(period/t.step)
}