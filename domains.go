package synthia

type TimeDomain float64

type FreqDomain complex128

func TimeDomain2PCM(signal []TimeDomain) []int32 {
	pcm := make([]int32, len(signal))
	for i, val := range(signal) {
		pcm[i] = int32(val)
	}
	return pcm
}