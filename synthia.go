package synthia

//go:generate -command yacc go tool yacc
//go:generate yacc -o lang.go -p "lang" lang.y

import "io/ioutil"

func FileToTune(path string, samplerate float64) ([]TimeDomain, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return []TimeDomain{}, err
	}
	langParse(&langLex{line: data})
	ast.Exec(samplerate)
	tune := ast.Tune()
	return tune, nil
}
