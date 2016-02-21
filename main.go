package main

import (
	"fmt"
	"flag"
	"io/ioutil"
)

//go:generate -command yacc go tool yacc
//go:generate yacc -o lang.go -p "lang" lang.y

const DefaultSampleRate = 44100

var glsampleRate float64 = DefaultSampleRate

func usage(){
	fmt.Println("synthia: The synth which goes...")
	fmt.Println("Usage: synthia FILE")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1{
		usage()
		return
	}
	path := args[0]
	data, err := ioutil.ReadFile(path)
	if err != nil {
		println(err)
		return
	}
	langParse(&langLex{line: data})
	
	ast.Exec()
	tune := ast.Tune()
	
	//var filter = NewLowPassFilter(301)
	//tune = filter.Filter(tune)
	fmt.Printf("%v\n", len(tune))
	playTune(tune, glsampleRate)

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
