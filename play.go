package main

import (
	"github.com/gordonklaus/portaudio"
	"fmt"
)

func playTune(tune []int32, sampleRate float64) error {
	err := portaudio.Initialize()
	if err != nil {
		return err
	}
	defer portaudio.Terminate()
	buffer := make([]int32, len(tune))
	copy(buffer, tune)
	fmt.Printf("%v\n", len(buffer))
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
	err = stream.Write()
	if err != nil {
		return err
	}
	return nil
}
