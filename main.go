package main
import(
	"github.com/draringi/synthia/waveforms"
	)


//go:generate -command yacc go tool yacc
//go:generate yacc -o lang.go -p "lang" lang.y


import ("fmt")

const DefaultSampleRate = 44100

func main() {
	var sampleRate float64 = DefaultSampleRate
	setBPM(110)
	song := genTwinkle()
	var tune []int32
	toneGenerator := NewToneGenerator(sampleRate,waveforms.Sin)
	for i, n := range(song) {
		fmt.Printf("%v\n", i)
		tune = append(tune, n.GenerateTone(toneGenerator)...)
	}
	playTune(tune, sampleRate)
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
