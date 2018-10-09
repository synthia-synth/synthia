package domains

type Time float64

type Frequency complex128

func TimeDomain2PCM(signal []Time) []int32 {
	pcm := make([]int32, len(signal))
	for i, val := range signal {
		pcm[i] = int32(val)
	}
	return pcm
}
