package main

import (
//	"math"
	"math/rand"
)

func whitenoise(samples int) []int32 {
	sampleSpace := make([]int32, samples)
	for i := 0; i < samples; i++ {
		sampleSpace[i] = int32((rand.Float32()*2-1)*(2 << 31 - 1))
	}
	return sampleSpace
}

func pluck(freq, time float64) []int32 {
	sampleSize := int(glsampleRate / freq)
	loops := int(freq*time)
	sound := whitenoise(sampleSize)
	lastIndex := sampleSize-1
	out := make([]int32, loops*sampleSize)
	k := 0
	for i:=0; i < loops; i++ {
		for j, _ := range(sound) {
			if j==lastIndex {
				sound[j] = int32(0.996*float64(sound[j]+sound[0])/2)
			} else {
				sound[j] = (sound[j]+sound[j+1])/2
			}
			out[k] = sound[j]
			k++
		}
	}
	return out
}
