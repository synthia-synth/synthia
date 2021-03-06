package synthia

import (
	"errors"
	"github.com/synthia-synth/synthia/waveforms"
	"github.com/synthia-synth/synthia/domains"
)

var (
	depth     int
	functions = map[string]function{
		"setBPM": (setBPMWrapper)(setBPM),
	}
	instruments map[string]ToneGenerator = map[string]ToneGenerator{}
)

var ast AST

type function interface {
	Name() string
	Exec([]expression)
}

type setBPMWrapper func(float64)

func (f setBPMWrapper) Name() string {
	return "setBPM"
}

func (f setBPMWrapper) Exec(args []expression) {
	if len(args) != 1 {
		panic("Ah! Panic!!!")
	}
	e := args[0].(intExp)
	callable := (func(float64))(f)
	callable(float64(e))
}

type astStream struct {
	instructions []instruction
	label        string
	tune         []domains.Time
}

func (s *astStream) Header(samplerate float64) {
	var tune []domains.Time
	for _, i := range s.instructions {
		i.Exec(samplerate)
		tune = append(tune, i.(*methodCall).tune...)
	}
	s.tune = tune
}

type instruction interface {
	Exec(samplerate float64)
}

type methodCall struct {
	obj       *object
	method    string
	arguments []expression
	tune      []domains.Time
}

func (m *methodCall) Exec(samplerate float64) {
	gen := instruments[m.obj.label]
	if gen == nil {
		panic("Not a valid generator!!!")
	}
	timing := m.arguments[1].(*timingExpression)
	play := m.arguments[0]
	switch play.(type) {
	case *noteExpression:
		noteInfo := play.(*noteExpression)
		n := NewNote(noteInfo.note, noteInfo.octave, noteInfo.accidental, timing.timing, timing.modifier)
		m.tune = n.GenerateTone(gen)
	case *chordExpression:
		chordInfo := play.(*chordExpression)
		var notes [][]domains.Time
		for _, n := range chordInfo.notes {
			print(n)
			note := NewNote(n.note, n.octave, n.accidental, timing.timing, timing.modifier)
			notes = append(notes, note.GenerateTone(gen))
		}
		m.tune = summer(notes...)
	}
}

type expression interface {
	Type() string
}

type chordExpression struct {
	notes []*noteExpression
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

func (f *functionCall) Exec(samplerate float64) {
	callable := functions[f.label]
	callable.Exec(f.arguments)
}

func (s *functionCall) Header(samplerate float64) {
	s.Exec(samplerate)
}

type instrumentInstance struct {
	label string
	inst  instrument
}

func (i *instrumentInstance) Exec(samplerate float64) {
	var generator ToneGenerator
	switch t := i.inst.(type) {
	case *tone:
		generator = NewWavetoneGenerator(samplerate, t.wave)
	case *simulation:
		generator = t.simulator(samplerate)
	default:
		panic("Unknown instrument type")
	}
	
	instruments[i.label] = generator
}

func (s *instrumentInstance) Header(samplerate float64) {
	s.Exec(samplerate)
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

type simulation struct {
	name string
	simulator func(samplerate float64) ToneSimulator
}

func (s *simulation) Name() string {
	return s.name
}

func (s *simulation) Type() string {
	return "Simulator"
}

type instrumentModule map[string]instrument

var (
	sinwave    = &tone{wave: waveforms.Sin, name: "sin"}
	triwave    = &tone{wave: waveforms.Tri, name: "triangle"}
	sawwave    = &tone{wave: waveforms.Saw, name: "saw"}
	sqrwave    = &tone{wave: waveforms.Sqr, name: "square"}
	nullwave   = &tone{wave: waveforms.Null, name: "null"}
	plucker = &simulation{simulator: NewPlucker, name: "pluck"}
	toneLookup = map[string]instrument{
		sinwave.name:  sinwave,
		triwave.name:  triwave,
		sawwave.name:  sawwave,
		sqrwave.name:  sqrwave,
		nullwave.name: nullwave,
	}
	simLookup = map[string]instrument{
		plucker.name: plucker,
	}
	instroModules = map[string]instrumentModule{
		"tone": toneLookup,
		"sim":simLookup,
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
	accidentalLookup = map[string]Accidental{
		"natural":     AccidentalNatural,
		"sharp":       AccidentalSharp,
		"flat":        AccidentalFlat,
		"doublesharp": AccidentalDoubleSharp,
		"doubleflat":  AccidentalDoubleFlat,
		"halfsharp":   AccidentalHalfSharp,
		"halfflat":    AccidentalHalfFlat,
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

type header interface {
	Header(samplerate float64)
}

type AST []header

func (a AST) Exec(samplerate float64) {
	headers := ([]header)(a)
	for _, h := range headers {
		h.Header(samplerate)
	}
}

func (a AST) Tune() []domains.Time {
	var tunes [][]domains.Time
	headers := ([]header)(a)
	for _, h := range headers {
		switch h.(type) {
		case *astStream:
			tunes = append(tunes, h.(*astStream).tune)
		default:
		}
	}
	return summer(tunes...)
}
