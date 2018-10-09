package synthia

import "github.com/synthia-synth/synthia/domains"

func summer(channels ...[]domains.Time) []domains.Time {
	maxSize := 0
	numChan := domains.Time(len(channels))
	for _, ch := range channels {
		if len(ch) > maxSize {
			maxSize = len(ch)
		}
	}

	var sum = make([]domains.Time, maxSize)
	for i := 0; i < maxSize; i++ { // current time
		for _, ch := range channels { //Each channel
			sum[i] += ch[i] / numChan
		}
	}

	return sum
}
