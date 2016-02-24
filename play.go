package main

import (
	"github.com/Orion90/portaudio"
)

const playBufferSize = 8192

func playTune(tune []int32, sampleRate float64) error {
	err := portaudio.Initialize()
	if err != nil {
		return err
	}
	defer portaudio.Terminate()
	buffer := make([]int32, playBufferSize)
	stream, err := portaudio.OpenDefaultStream(0, 1, sampleRate, len(buffer), &buffer)
	if err != nil {
		return err
	}
	defer stream.Close()
	err = stream.Start()
	if err != nil {
		return err
	}
	defer stream.Stop()
	for i := 0; i < len(tune); i += len(buffer) {
		end := i + playBufferSize
		if end > len(tune) {
			copy(buffer, tune[i:])
		} else {
			copy(buffer, tune[i:end])
		}
		err = stream.Write()
		if err != nil {
			return err
		}
	}
	return nil
}
