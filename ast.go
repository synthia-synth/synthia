package main

import (
	"errors"
	"github.com/draringi/synthia/waveforms"
)

var ()

type astStream struct {
	instructions []instruction
	label        string
}

type instruction interface {
	Exec()
}

type methodCall struct {
	obj       *object
	method    string
	arguments []expression
}

func (m *methodCall) Exec() {

}

type expression interface {
	Type() string
}

type chordExpression struct {
	notes []noteExpression
}

func (e *chordExpression) Type() string {
	return "Chord"
}

func (e *chordExpression) IsNotes() bool {
	return true
}

type noteExpression struct {
	note       NoteName
	octave     int
	accidental Accidental
}

func (e *noteExpression) Type() string {
	return "Note"
}

func (e *noteExpression) IsNotes() bool {
	return true
}

type tonePlayMethod struct {
	timing *timingExpression
	notes  interface {
		IsNotes() bool
	}
}

type timingExpression struct {
	timing   NoteLen
	modifier LenModifier
}

func (e *timingExpression) Type() string {
	return "Timing"
}

type object struct {
	label string
}

type intExp int

func (i intExp) Type() string {
	return "Integer"
}

type functionCall struct {
	label     string
	arguments []expression
}

func (f *functionCall) Exec() {

}

type instrumentInstance struct {
	label string
	inst  instrument
}

func (i *instrumentInstance) Exec() {

}

type instrument interface {
	Name() string
	Type() string
}

type voice struct {
	voiceData interface{}
}

type tone struct {
	wave waveforms.Wave
	name string
}

func (t *tone) Name() string {
	return t.name
}

func (t *tone) Type() string {
	return "ToneGenerator"
}

type instrumentModule map[string]instrument

var (
	sinwave    = &tone{wave: waveforms.Sin, name: "sin"}
	triwave    = &tone{wave: waveforms.Tri, name: "triangle"}
	sawwave    = &tone{wave: waveforms.Saw, name: "saw"}
	sqrwave    = &tone{wave: waveforms.Sqr, name: "square"}
	toneLookup = map[string]instrument{
		sinwave.name: sinwave,
		triwave.name: triwave,
		sawwave.name: sawwave,
		sqrwave.name: sqrwave,
	}
	instroModules = map[string]instrumentModule{
		"tone": toneLookup,
	}
	timingLookup = map[string]NoteLen{
		"breve":                      Breve,
		"semibreve":                  SemiBreve,
		"minim":                      Minim,
		"crotchet":                   Crotchet,
		"quaver":                     Quaver,
		"semiquaver":                 SemiQuaver,
		"demisemiquaver":             DemiSemiQuaver,
		"hemidemisemiquaver":         HemiDemiSemiQuaver,
		"semihemidemisemiquaver":     SemiHemiDemiSemiQuaver,
		"demisemihemidemisemiquaver": DemiSemiHemiDemiSemiQuaver,
	}
	noteLookup = map[string]NoteName{
		"A": A,
		"B": B,
		"C": C,
		"D": D,
		"E": E,
		"F": F,
		"G": G,
	}
)

func instrumentLookup(module, name string) (instrument, error) {
	m, exists := instroModules[module]
	if !exists {
		return nil, errors.New("Invalid module name")
	}
	i, exists := m[name]
	if !exists {
		return nil, errors.New("Invalid instrument name")
	}
	return i, nil
}
