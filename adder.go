package main

func summer(channels ...[]int32) []int32{
	var maxSize = 0
	var numChan = size(channels)
	for i, ch := range channels {
		if size(ch(i)) > maxSize{
			maxSize = size(ch(i))
		}
	}

	var sum = make([]int32, maxSize)
	for i=0; i<maxSize; i++{ // current time
		for ch := range channels{ //Each channel
			sum[i] = sum[i] + ch[i]/numChan
		}
	}

	return sum
}