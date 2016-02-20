package main

//go:generate -command yacc go tool yacc
//go:generate yacc -o lang.go -p "lang" lang.y

const DefaultSampleRate = 44100

func main() {
	var sampleRate float64 = DefaultSampleRate
	setBPM(60)
	var song []*Note
	song = append(song, NewNote(A, 4, AccidentalNatural, Breve, NormalLength))
	song = append(song, NewNote(C, 4, AccidentalNatural, SemiQuaver, NormalLength))
	song = append(song, NewNote(B, 4, AccidentalNatural, SemiQuaver, NormalLength))
	song = append(song, NewNote(C, 4, AccidentalNatural, SemiQuaver, NormalLength))
	var tune []int32
	
}
