package synthia

//go:generate -command yacc go tool yacc
//go:generate yacc -o lang.go -p "lang" lang.y

import (
	"io/ioutil"
	"github.com/synthia-synth/synthia/domains"
)

func FileToTune(path string, samplerate float64) ([]domains.Time, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return []domains.Time{}, err
	}
	langParse(&langLex{line: data})
	ast.Exec(samplerate)
	tune := ast.Tune()
	return tune, nil
}
