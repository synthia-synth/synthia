package synthia

import (
	"math"
	"github.com/synthia-synth/synthia/domains"
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
	Quaver                     NoteLen = 1.0 / 2
	SemiQuaver                 NoteLen = 1.0 / 4
	DemiSemiQuaver             NoteLen = 1.0 / 8
	HemiDemiSemiQuaver         NoteLen = 1.0 / 16
	SemiHemiDemiSemiQuaver     NoteLen = 1.0 / 32
	DemiSemiHemiDemiSemiQuaver NoteLen = 1.0 / 64
)

type LenModifier float64

const (
	Dotted       LenModifier = 1.5
	NormalLength LenModifier = 1
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
	length         NoteLen
	lengthModifier LenModifier
	octave         int
	vol            int32
}

func setBPM(newBPM float64) {
	bpm = newBPM
	beatLength = 60 / bpm
}

func lengthToDuration(len NoteLen, modifier LenModifier) float64 {
	return float64(len) * beatLength * float64(modifier)
}

func (n *Note) Frequency() float64 {
	note := int(n.note) + 12*n.octave
	diff := float64(note - referencePoint)
	diff += float64(n.accidental)
	return referenceFreq * math.Pow(freqStep, diff)
}

func NewNote(note NoteName, octave int, accidental Accidental, length NoteLen, lengthModifier LenModifier) *Note {
	n := new(Note)
	n.note = note
	n.octave = octave
	n.accidental = accidental
	n.length = length
	n.lengthModifier = lengthModifier
	n.vol = 1 << 30
	return n
}

func (n *Note) GenerateTone(generator ToneGenerator) []domains.Time {
	return adsr1(generator.Play(n.Frequency(), lengthToDuration(n.length, n.lengthModifier), n.vol))
}
