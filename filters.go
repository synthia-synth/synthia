package main

import("math")

func adsr1 (signal []int32) []int32{
	const riseSamples, decaySamples, releaseSamples = 20, 40, 10
	const sustainLvl = .8
	var baseNeed = int(riseSamples + decaySamples + releaseSamples)

	//Constant rise 0->1
	if baseNeed > len(signal){ //Shrink ranges
		return signal
	}

	var stoppingPoint = riseSamples
	var curve = 0.0

	for i=0; i<stoppingPoint; i++{ //Rise (lin)
		curve = float32(i)/float32(riseSamples)
		signal[i] = signal[i] * int32(curve)
	}
	stoppingPoint = riseSamples + stoppingPointdecaySamples
	var lambda = float32(decaySamples)/Math.Ln(sustainLvl)
	for i=riseSamples; i<stoppingPoint; i++{ //Decay (exp)
		curve = Math.Exp(float32(i-riseSamples+1)/lambda)
		signal[i] = signal[i] * int32(curve)
	}
	stoppingPoint = len(signal)-releaseSamples
	for i=riseSamples+decaySamples; i<stoppingPoint; i++{ //Sustain (const)
		signal[i] = signal[i] * sustainLvl
	}
	stoppingPoint = len(signal)
	lambda = float32(releaseSamples)/Math.Ln(sustainLvl)
	for i=len(singal)-releaseSamples; i<stoppingPoint; i++{ //Release (exp)
		curve = sustainLvl*Math.Exp(float32(i-len(signal)-releaseSamples+1)/lambda)
		signal[i] = signal[i] * int32(curve)
	}

	return signal
}