package synthia

type SignalWriter interface {
	WriteSignal(s []TimeDomain, samplerate float64)  error
}

type PCMWriter interface {
	WritePCM(p []int32, samplerate float64) error
}

type FrequencyWriter interface {
	WriteFreqDomain(s []FreqDomain) error
}

type Signal2PCMWriter struct {
	pcmWriter PCMWriter
}

func (w Signal2PCMWriter) WriteSignal(s []TimeDomain, samplerate float64)  error {
	return w.pcmWriter.WritePCM(TimeDomain2PCM(s), samplerate)
}

func (w Signal2PCMWriter) WritePCM(p []int32, samplerate float64) error {
	return w.pcmWriter.WritePCM(p, samplerate)
}

func NewSignal2PCMWriter(pcm PCMWriter) Signal2PCMWriter {
	return Signal2PCMWriter{pcm}
}