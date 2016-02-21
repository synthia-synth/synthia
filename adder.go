package main

func summer(channels ...[]int32) []int32 {
	maxSize := 0
	numChan := int32(len(channels))
	for _, ch := range channels {
		if len(ch) > maxSize {
			maxSize = len(ch)
		}
	}

	var sum = make([]int32, maxSize)
	for i := 0; i < maxSize; i++ { // current time
		for _, ch := range channels { //Each channel
			sum[i] = sum[i] + ch[i]/numChan
		}
	}

	return sum
}
