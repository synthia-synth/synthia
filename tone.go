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
	var period = period(freq)

	//How many samples are there in a period
	var samplesPerPeriod = int(period/t.step)

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
	//Find one period
	//Width: first half = 1, second = 0
	//Repeat n.m times
	for i:=0; i < len(synthArray); i++{
		synthArray[i] = int32(float64(vol)*math.Sin(freq *2* math.Pi * float64(i) * t.step))
	}
	return synthArray
}

//Generates a triangle wave
func (t *ToneGenerator) TriTone(freq, seconds float64, vol int32) []int32{
	var synthArray = make([]int32, int(seconds*t.sampleRate)) //duration/step = dur*sR
	//Find one period
	//Width: first half = 1, second = 0
	//Repeat n.m times
	for i:=0; i < len(synthArray); i++{
		synthArray[i] = int32(float64(vol)*math.Sin(freq *2* math.Pi * float64(i) * t.step))
	}
	return synthArray
}


func period(freq float64) float64{
	return 1.0/freq 
}