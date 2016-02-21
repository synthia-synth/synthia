%{
package main
import (
	"log"
	"strconv"
	"bytes"
	"unicode/utf8"
)
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
}

%token '(' ')' '{' '}' '[' ']' '.' '=' '\n'
%token STREAM
%token <str> LABEL
%token <integer> NUM NOTE
%token <expr> TIMING

%type <stream> strm
%type <instructions> inst_list
%type <inst> inst
%type <expr> expr
%type <expressions> expr_list
%type <instumentDef> definition

%%

strm:
	STREAM LABEL '{' inst_list '}'
	{
		$$ = &astStream{label: $2, instructions: $4}
	}

inst_list:
	inst '\n'
	{
		$$ = []instruction{$1}
	}
	| inst_list inst '\n'
	{
		$$ = append($1, $2)
	}

inst:
	LABEL '.' LABEL '(' expr_list ')'
	{
		$$ = &methodCall{obj: &object{label: $1}, method: $3, arguments: $5}
	}
	| LABEL '(' expr_list ')'
	{
		$$ = &functionCall{label: $1, arguments: $3}
	}
	| LABEL '=' definition
	{
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
	| expr_list expr
	{
		$$ = append($1, $2)
	}

expr:
	NOTE '[' NUM ']'
	{
		$$ = &noteExpression{note: NoteName($1), octave: $3}
	}
	| NUM
	{
		$$ = intExp($1)
	}
	| TIMING {
		$$ = $1
	}
%%
// The parser expects the lexer to return 0 on EOF.  Give it a name
// for clarity.
const eof = 0

// The parser uses the type <prefix>Lex as a lexer.  It must provide
// the methods Lex(*<prefix>SymType) int and Error(string).
type exprLex struct {
	line []byte
	peek rune
}

// The parser calls this method to get each new token.  This
// implementation returns operators and NUM.
func (x *exprLex) Lex(yylval *langSymType) int {
	for {
		c := x.next()
		switch c {
		case eof:
			return eof
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return x.num(c, yylval)
		case '+', '-', '*', '/', '(', ')', '\n', '[', ']', '.', '{', '}':
			return int(c)

		// Recognize Unicode multiplication and division
		// symbols, returning what the parser expects.
		case '×':
			return '*'
		case '÷':
			return '/'

		case ' ', '\t', '\r':
		default:
			return x.label(c, yylval)
		}
	}
}

// Lex a number.
func (x *exprLex) num(c rune, yylval *langSymType) int {
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

func (x *exprLex) label(c rune, yylval *langSymType) int {
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
		case ' ', '.', '\t', '+', '-', '*', '/', '(', ')', '\n', '[', ']', '{', '}', eof:
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
		if c == '.' {
			modifier = Dotted
		}
		yylval.expr = &timingExpression{ timing: timing, modifier: modifier}
		return TIMING
	}
	note, isNote := noteLookup[yylval.str]
	if isNote {
		
	}
	return LABEL
}

// Return the next rune for the lexer.
func (x *exprLex) next() rune {
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
func (x *exprLex) Error(s string) {
	log.Printf("parse error: %s", s)
}
