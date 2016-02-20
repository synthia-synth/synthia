package main

import (
	"math"
)

type NoteName int

const (
	A NoteName = 9
	B          = 11
	C          = 0
	D          = 2
	E          = 4
	F          = 5
	G          = 7
)

type Accidental float64

const (
	AccidentalNatural     Accidental = 0
	AccidentalSharp                  = 1
	AccidentalFlat                   = -1
	AccidentalDoubleSharp            = 2
	AccidentalDoubleFlat             = -2
	AccidentalHalfSharp              = 0.5
	AccidentalHalfFlat               = -.5
)

type NoteLen float64

const (
	Breve                      NoteLen = 8
	SemiBreve                  NoteLen = 4
	Minim                      NoteLen = 2
	Crotchet                   NoteLen = 1
	Quaver                     NoteLen = 1 / 2
	SemiQuaver                 NoteLen = 1 / 4
	DemiSemiQuaver             NoteLen = 1 / 8
	HemiDemiSemiQuaver         NoteLen = 1 / 16
	SemiHemiDemiSemiQuaver     NoteLen = 1 / 32
	DemiSemiHemiDemiSemiQuaver NoteLen = 1 / 64
)

type LenModifier float64

const (
	Dotted LenModifier = 1.5
)

const (
	referencePoint = 4*12 + int(A)
	referenceFreq  = 440
	freqStep       = 1.059463094359
)

var bpm, beatLength float64

type Note struct {
	note           NoteName
	accidental     Accidental
	length         NoteLength
	lengthModifier LenModifier
	octave         int
}

func setBPM(newBPM float64) {
	bpm = newBPM
	beatLength = 60 / bpm
}

func lengthToDuration(len NoteLen, modifier LenModifier) float64 {
	return float64(len) * beatLength * float64(modifier)
}

func (n *Note) Frequency() float64 {
	n := int(n.note)
	diff := float64(n - referencePoint)
	diff += float64(n.accidental)
	return referenceFreq * math.Pow(freqStep, diff)
}
