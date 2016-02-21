package main

import (
	"github.com/draringi/synthia/waveforms"
)

type astStream struct {
	instructions []instruction
	label string
}

type instruction interface {
	Exec()
}

type methodCall struct {
	obj *object
	method string
	arguments []expression
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

type noteExpression struct {
	note NoteName
	octave int
	accidental Accidental
}

type tonePlayMethod struct {
	timing *timingExpression
	notes expression
}

func (e *noteExpression) Type() string {
	return "Note"
}

type timingExpression struct {
	timing NoteLen
	modifier LenModifier
}

func (e *timingExpression) Type() string {
	return "Timing"
}

type object struct {
	label string
}

func (m *methodCall) Exec() {
	
}

type instrumentInstance struct {
	label string
	inst instrument
}

type instrument interface {
	Name() string
	Type() string
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

var (
	sinwave *tone = &tone{wave: waveforms.Sin, name: "sin"}
	triwave *tone = &tone{wave: waveforms.Tri, name: "triangle"}
)
