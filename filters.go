package synthia

import (
	"math"
	"github.com/synthia-synth/synthia/domains"
)

func adsr1(signal []domains.Time) []domains.Time {
	const riseSamples, decaySamples, releaseSamples = 200, 400, 200
	const sustainLvl = .7
	var baseNeed = int(riseSamples + decaySamples + releaseSamples)

	//Constant rise 0->1
	if baseNeed > len(signal) { //Shrink ranges
		return signal
	}

	var stoppingPoint = riseSamples
	var curve = 0.0

	for i := 0; i < stoppingPoint; i++ { //Rise (lin)
		curve = float64(i) / float64(riseSamples)
		signal[i] = signal[i] * domains.Time(curve)
	}
	stoppingPoint = riseSamples + decaySamples
	var lambda = float64(decaySamples) / math.Log(sustainLvl)
	for i := riseSamples; i < stoppingPoint; i++ { //Decay (exp)
		curve = math.Exp(float64(i-riseSamples+1) / lambda)
		signal[i] = signal[i] * domains.Time(curve)
	}
	stoppingPoint = len(signal) - releaseSamples
	for i := riseSamples + decaySamples; i < stoppingPoint; i++ { //Sustain (const)
		signal[i] = domains.Time(float64(signal[i]) * sustainLvl)
	}
	stoppingPoint = len(signal)
	lambda = float64(releaseSamples) / math.Log(sustainLvl)
	for i := len(signal) - releaseSamples; i < stoppingPoint; i++ { //Release (exp)
		curve = sustainLvl * math.Exp(float64(i-len(signal)-releaseSamples+1)/lambda)
		signal[i] = signal[i] * domains.Time(curve)
	}

	return signal
}
