%{
package synthia
import (
	"log"
	"strconv"
	"bytes"
	"unicode/utf8"
)

var stm = false
var line int = 1
%}

%union {
	stream *astStream
	str string
	instructions []instruction
	inst instruction
	expr expression
	expressions []expression
	integer int
	instumentDef instrument
	headers []header
	headerField header
	note *noteExpression
	notes []*noteExpression
}

%token '(' ')' '{' '}' '[' ']' '.' '=' ';'
%token STREAM
%token <str> LABEL
%token <integer> NUM
%token <expr> TIMING
%token <note> NOTE

%type <stream> strm
%type <instructions> inst_list
%type <inst> inst
%type <expr> expr chord
%type <expressions> expr_list
%type <instumentDef> definition
%type <headers> top head_list
%type <headerField> head
%type <note> note_el
%type <notes> note_list

%%

top: head_list
	{ ast = AST($1) }
	

head_list:
	head ';'
	{
		$$ = []header{$1}
	}
	| head_list head ';'
	{
		$$ = append($1, $2)
	}
	
head:
	strm {
		stm = true
		$$ = $1
	}
	| LABEL '(' expr_list ')'
	{
		stm = true
		$$ = &functionCall{label: $1, arguments: $3}
	}
	| LABEL '=' definition
	{
		stm = true
		$$ = &instrumentInstance{label: $1, inst: $3}
	}

strm:
	STREAM LABEL '{' inst_list '}'
	{
		$$ = &astStream{label: $2, instructions: $4}
	}

inst_list:
	inst ';'
	{
		$$ = []instruction{$1}
	}
	| inst_list inst ';'
	{
		$$ = append($1, $2)
	}

inst:
	LABEL '.' LABEL '(' expr_list ')'
	{
		stm = true
		$$ = &methodCall{obj: &object{label: $1}, method: $3, arguments: $5}
	}
	| LABEL '(' expr_list ')'
	{
		stm = true
		$$ = &functionCall{label: $1, arguments: $3}
	}
	| LABEL '=' definition
	{
		stm = true
		$$ = &instrumentInstance{label: $1, inst: $3}
	}

definition:
	LABEL '.' LABEL
	{
		i, err := instrumentLookup($1, $3)
		if err != nil {
			log.Fatalf("Instrument lookup with (%s,%s) failed: %s\n", $1, $3, err)
		}
		$$ = i
	}

expr_list:
	expr
	{
		$$ = []expression{$1}
	}
	| expr_list ',' expr
	{
		$$ = append($1, $3)
	}

expr:
	note_el { $$ = $1 }
	| NUM { $$ = intExp($1) }
	| TIMING { $$ = $1 }
	| chord { $$ = $1 }
	
chord:
	'(' note_list ')' { $$ = &chordExpression{notes: $2} }
	
note_list:
	note_el { $$ = []*noteExpression{$1} }
	| note_list ',' note_el { $$ = append($1, $3) }
	
note_el:
	NOTE
	{
		$$ = $1
	}
%%
// The parser expects the lexer to return 0 on EOF.  Give it a name
// for clarity.
const eof = 0

// The parser uses the type <prefix>Lex as a lexer.  It must provide
// the methods Lex(*<prefix>SymType) int and Error(string).
type langLex struct {
	line []byte
	peek rune
}

// The parser calls this method to get each new token.  This
// implementation returns operators and NUM.
func (x *langLex) Lex(yylval *langSymType) int {
	for {
		c := x.next()
		switch c {
		case eof:
			return eof
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return x.num(c, yylval)
		case '+', '-', '*', '/', '(', ')', ';', '[', ']', '.', '{', '}', '=', ',':
			return int(c)

		// Recognize Unicode multiplication and division
		// symbols, returning what the parser expects.
		case '×':
			return '*'
		case '÷':
			return '/'
		case '\n':
			if stm {
				x.peek = ';'
				stm = false
			}
			line++
		case ' ', '\t', '\r':
		default:
			return x.label(c, yylval)
		}
	}
}

// Lex a number.
func (x *langLex) num(c rune, yylval *langSymType) int {
	add := func(b *bytes.Buffer, c rune) {
		if _, err := b.WriteRune(c); err != nil {
			log.Fatalf("WriteRune: %s", err)
		}
	}
	var b bytes.Buffer
	add(&b, c)
	L: for {
		c = x.next()
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			add(&b, c)
		default:
			break L
		}
	}
	if c != eof {
		x.peek = c
	}
	var err error
	yylval.integer, err = strconv.Atoi(b.String())
	if err != nil {
		log.Fatalf("bad number %q, error: %s", b.String(), err)
		return eof
	}
	return NUM
}

func (x *langLex) label(c rune, yylval *langSymType) int {
	add := func(b *bytes.Buffer, c rune) {
		if _, err := b.WriteRune(c); err != nil {
			log.Fatalf("WriteRune: %s", err)
		}
	}
	var b bytes.Buffer
	add(&b, c)
	L: for {
		c = x.next()
		switch c {
		case ' ', '.', '\t', '+', '-', '*', '/', '(', ')', '\n', '[', ']', '{', '}', '=', ',', eof:
			break L
		default:
			add(&b, c)
		}
	}
	if c != eof {
		x.peek = c
	}

	yylval.str = b.String()
	timing, isTiming := timingLookup[yylval.str]
	if isTiming {
		modifier := NormalLength
		c = x.next()
		if c == '.' {
			modifier = Dotted
		} else {
			x.peek = c
		}
		yylval.expr = &timingExpression{ timing: timing, modifier: modifier}
		return TIMING
	}
	note, isNote := noteLookup[yylval.str]
	if isNote {
		c = x.next()
		if c != '['{
			log.Fatalf("Invalid Note %s with no octave. Should be %s[n] where n is a single-digit integer\n", yylval.str, yylval.str)
		}
		c = x.next()
		var octave int
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			var err error
			octave, err = strconv.Atoi(string(c))
			if err != nil {
				log.Fatal("Critical Error: Integer is not an Integer!\n")
			}
		default:
			log.Fatalf("Invalid Note %s with no octave. Should be %s[n] where n is a single-digit integer\n", yylval.str, yylval.str)
		}
		c = x.next()
		if c != ']'{
			log.Fatalf("Invalid Note %s with no octave. Should be %s[n] where n is a single-digit integer\n", yylval.str, yylval.str)
		}
		acc := AccidentalNatural
		c = x.next()
		if c == '.' {
			var a bytes.Buffer
			LA: for {
				c = x.next()
				switch c {
				case ' ', '.', '\t', '+', '-', '*', '/', '(', ')', '\n', '[', ']', '{', '}', ',', '=', eof:
					break LA
				default:
					add(&a, c)
				}
			}
			ac, isAccidental := accidentalLookup[a.String()]
			if !isAccidental {
				log.Fatalf("%s is not a valid accidental\n", a.String())
			}
			if c != eof {
				x.peek = c
			}
			acc = ac
		} else {
			x.peek = c
		}
		yylval.note = &noteExpression{ note: note, octave: octave, accidental: acc }
		return NOTE
	}
	if b.String() == "stream" {
		return STREAM
	}
	return LABEL
}

// Return the next rune for the lexer.
func (x *langLex) next() rune {
	if x.peek != eof {
		r := x.peek
		x.peek = eof
		return r
	}
	if len(x.line) == 0 {
		return eof
	}
	c, size := utf8.DecodeRune(x.line)
	x.line = x.line[size:]
	if c == utf8.RuneError && size == 1 {
		log.Print("invalid utf8")
		return x.next()
	}
	return c
}

// The parser calls this method on a parse error.
func (x *langLex) Error(s string) {
	c := x.next()
	var cString string
	switch c {
		case ' ':
			cString = "space"
		case '\n':
			cString = "\\n"
		case '\t':
			cString = "\\t"
		case eof:
			cString = "EOF"
		default:
			cString = string(c)
	}
	log.Printf("parse error on line %d before character %s: %s", line, cString, s)
}
