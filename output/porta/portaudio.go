package porta

import (
	"github.com/gordonklaus/portaudio"
)

const (
	DefaultBufferSize = 4096
)

type PAOutput struct {
	buffer []int32
}

func NewPAOutput(bufferSize int32) (*PAOutput) {
	paOutput := new(PAOutput)
	paOutput.buffer = make([]int32, bufferSize)
	return paOutput
}

func (p *PAOutput)WritePCM(pcm []int32, samplerate float64) error {
	err := portaudio.Initialize()
	if err != nil {
		return err
	}
	defer portaudio.Terminate()
	stream, err := portaudio.OpenDefaultStream(0, 1, samplerate, len(p.buffer), &p.buffer)
	if err != nil {
		return err
	}
	defer stream.Close()
	stream.Start()
	for i:= 0; i < len(pcm); i += len(p.buffer) {
		end := i + len(p.buffer)
		if end > len(pcm) {
			copy(p.buffer, pcm[i:])
		} else {
			copy(p.buffer, pcm[i:end])
		}
		err = stream.Write()
		if err != nil {
			return err
		}
	}
	return nil
}