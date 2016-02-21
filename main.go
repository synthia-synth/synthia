package main

import (
	"fmt"
	"github.com/draringi/synthia/waveforms"
)

//go:generate -command yacc go tool yacc
//go:generate yacc -o lang.go -p "lang" lang.y

const DefaultSampleRate = 44100

func main() {
	var sampleRate float64 = DefaultSampleRate
	setBPM(110)
	song := genTwinkle()
	var tune []int32
	toneGenerator := NewToneGenerator(sampleRate, waveforms.Saw)
	for _, n := range song {
		tune = append(tune, n.GenerateTone(toneGenerator)...)
	}

	var testFilter = NewLowPassFilter(11)
	fmt.Printf("%v\n", testFilter.window)
	var filter = NewLowPassFilter(301)
	tune = filter.Filter(tune)
	fmt.Printf("%v\n", len(tune))
	if playTune(tune, sampleRate) == nil {
		fmt.Printf("all ended correctly\n")
	}

}

func genTwinkle() []*Note {
	var song []*Note
	song = append(song, NewNote(C, 5, AccidentalNatural, Crotchet, NormalLength))
	song = append(song, NewNote(C, 5, AccidentalNatural, Crotchet, NormalLength))
	song = append(song, NewNote(G, 5, AccidentalNatural, Crotchet, NormalLength))
	song = append(song, NewNote(G, 5, AccidentalNatural, Crotchet, NormalLength))
	song = append(song, NewNote(A, 5, AccidentalNatural, Crotchet, NormalLength))
	song = append(song, NewNote(A, 5, AccidentalNatural, Crotchet, NormalLength))
	song = append(song, NewNote(G, 5, AccidentalNatural, Minim, NormalLength))
	song = append(song, NewNote(C, 5, AccidentalNatural, Crotchet, NormalLength))
	song = append(song, NewNote(A, 5, AccidentalNatural, Crotchet, NormalLength))
	song = append(song, NewNote(G, 5, AccidentalNatural, Crotchet, NormalLength))
	song = append(song, NewNote(B, 5, 0.1, Crotchet, NormalLength))
	song = append(song, NewNote(B, 5, 0.2, Crotchet, NormalLength))
	song = append(song, NewNote(B, 5, 0.3, Crotchet, NormalLength))
	song = append(song, NewNote(B, 5, 0.4, Crotchet, NormalLength))
	song = append(song, NewNote(B, 5, 0.5, Crotchet, NormalLength))
	song = append(song, NewNote(B, 5, 0.6, Crotchet, NormalLength))
	song = append(song, NewNote(B, 5, 0.7, Crotchet, NormalLength))
	song = append(song, NewNote(B, 5, 0.8, Crotchet, NormalLength))
	song = append(song, NewNote(B, 5, 0.9, Crotchet, NormalLength))
	song = append(song, NewNote(B, 5, 1.0, Crotchet, NormalLength))
	song = append(song, NewNote(B, 5, 1.1, Crotchet, NormalLength))
	song = append(song, NewNote(B, 5, 1.2, Crotchet, NormalLength))
	return song
}

func genFast() []*Note {
	var song []*Note
	song = append(song, NewNote(C, 5, AccidentalNatural, HemiDemiSemiQuaver, NormalLength))
	song = append(song, NewNote(A, 5, AccidentalNatural, HemiDemiSemiQuaver, NormalLength))
	song = append(song, NewNote(G, 5, AccidentalNatural, HemiDemiSemiQuaver, NormalLength))
	song = append(song, NewNote(B, 5, AccidentalNatural, HemiDemiSemiQuaver, NormalLength))
	song = append(song, NewNote(A, 5, AccidentalNatural, HemiDemiSemiQuaver, NormalLength))
	song = append(song, NewNote(C, 5, AccidentalNatural, HemiDemiSemiQuaver, NormalLength))
	song = append(song, NewNote(G, 5, AccidentalNatural, HemiDemiSemiQuaver, NormalLength))
	song = append(song, NewNote(C, 5, AccidentalNatural, HemiDemiSemiQuaver, NormalLength))
	song = append(song, NewNote(A, 5, AccidentalNatural, HemiDemiSemiQuaver, NormalLength))
	song = append(song, NewNote(G, 5, AccidentalNatural, HemiDemiSemiQuaver, NormalLength))
	song = append(song, NewNote(B, 5, AccidentalNatural, HemiDemiSemiQuaver, NormalLength))
	song = append(song, NewNote(A, 5, AccidentalNatural, HemiDemiSemiQuaver, NormalLength))
	song = append(song, NewNote(C, 5, AccidentalNatural, HemiDemiSemiQuaver, NormalLength))
	song = append(song, NewNote(G, 5, AccidentalNatural, HemiDemiSemiQuaver, NormalLength))
	return song
}
