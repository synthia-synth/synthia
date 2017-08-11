package synthia

func summer(channels ...[]TimeDomain) []TimeDomain {
	maxSize := 0
	numChan := TimeDomain(len(channels))
	for _, ch := range channels {
		if len(ch) > maxSize {
			maxSize = len(ch)
		}
	}

	var sum = make([]TimeDomain, maxSize)
	for i := 0; i < maxSize; i++ { // current time
		for _, ch := range channels { //Each channel
			sum[i] +=  ch[i]/numChan
		}
	}

	return sum
}
