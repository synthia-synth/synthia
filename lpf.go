package main
import ("math")

type LowPassFilter struct {
	window []float64
}

func NewLowPassFilter( bandwidth int64 ) *LowPassFilter {

	//this is basically a low pass filter based on the Hamming (Hanning?) window
	//it doesn't really work too well

	const a = 0.50
	var inverseBandwidth = 1./(float64(bandwidth+1)*.5)
	var filter = new(LowPassFilter)
	filter.window = make([]float64, bandwidth)
	var sum float64 = 0
	for i:=0; i<int(bandwidth); i++ {
		filter.window[i] = (a + (1.-a)*math.Cos(math.Pi*(float64(i)-float64(bandwidth-1)*0.5)*inverseBandwidth))*0.25
		sum += filter.window[i]
	}
	/*
	//this is to make sure that the lowPassed will never be louder than original
	//oops, doesn't work
	var inverseSum = 1./sum
	for i:=0; i<int(bandwidth); i++ {
		filter.window[i] *= inverseSum
	}
	*/
	return filter
}

func (filter LowPassFilter) Filter( input []int32 ) []int32{
	var output = make([]int32, len(input))
	for i:=0; i<len(output); i++ {
		var accumulator float64 = 0
		for j:=0; j<len(filter.window); j++ {
			if i-len(filter.window)+1+j >= 0 {
				sample := float64(input[i-len(filter.window)+1+j])/float64(2<<31)
				accumulator += filter.window[j] * sample
			}
		}
		output[i] = int32(accumulator*float64(2<<31))
	}
	return output
}
	

