//line lang.y:2
package main

import __yyfmt__ "fmt"

//line lang.y:2
import (
	"bytes"
	"log"
	"strconv"
	"unicode/utf8"
)

var stm = false
var line int = 1

//line lang.y:14
type langSymType struct {
	yys          int
	stream       *astStream
	str          string
	instructions []instruction
	inst         instruction
	expr         expression
	expressions  []expression
	integer      int
	instumentDef instrument
	headers      []header
	headerField  header
	note         *noteExpression
	notes        []*noteExpression
}

const STREAM = 57346
const LABEL = 57347
const NUM = 57348
const TIMING = 57349
const NOTE = 57350

var langToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"'('",
	"')'",
	"'{'",
	"'}'",
	"'['",
	"']'",
	"'.'",
	"'='",
	"';'",
	"STREAM",
	"LABEL",
	"NUM",
	"TIMING",
	"NOTE",
	"','",
}
var langStatenames = [...]string{}

const langEofCode = 1
const langErrCode = 2
const langMaxDepth = 200

//line lang.y:150

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
L:
	for {
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
L:
	for {
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
		yylval.expr = &timingExpression{timing: timing, modifier: modifier}
		return TIMING
	}
	note, isNote := noteLookup[yylval.str]
	if isNote {
		c = x.next()
		if c != '[' {
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
		if c != ']' {
			log.Fatalf("Invalid Note %s with no octave. Should be %s[n] where n is a single-digit integer\n", yylval.str, yylval.str)
		}
		acc := AccidentalNatural
		c = x.next()
		if c == '.' {
			var a bytes.Buffer
		LA:
			for {
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
		yylval.expr = &noteExpression{note: note, octave: octave, accidental: acc}
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

//line yacctab:1
var langExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const langNprod = 24
const langPrivate = 57344

var langTokenNames []string
var langStates []string

const langLast = 54

var langAct = [...]int{

	13, 21, 15, 30, 14, 20, 19, 50, 48, 33,
	22, 24, 36, 6, 5, 23, 16, 17, 19, 31,
	25, 25, 34, 27, 25, 44, 35, 31, 11, 40,
	32, 43, 38, 37, 12, 39, 41, 42, 9, 8,
	28, 45, 47, 46, 3, 10, 26, 7, 49, 2,
	1, 18, 29, 4,
}
var langPact = [...]int{

	0, -1000, 0, 27, -1000, 34, 14, 22, -1000, 1,
	-4, 9, -1000, 6, -1000, -1000, -1000, -1000, -1000, -1000,
	-11, -1000, 30, 13, -1000, 1, 4, -1000, 12, 5,
	20, 25, -1000, -1000, -11, -1000, -1000, 19, -1000, 11,
	1, -4, -1000, -1000, 38, 3, -1000, 1, -1000, 2,
	-1000,
}
var langPgo = [...]int{

	0, 53, 52, 3, 4, 51, 0, 1, 50, 49,
	44, 2, 46,
}
var langR1 = [...]int{

	0, 8, 9, 9, 10, 10, 10, 1, 2, 2,
	3, 3, 3, 7, 6, 6, 4, 4, 4, 4,
	5, 12, 12, 11,
}
var langR2 = [...]int{

	0, 1, 2, 3, 1, 4, 3, 5, 2, 3,
	6, 4, 3, 3, 1, 3, 1, 1, 1, 1,
	3, 1, 3, 1,
}
var langChk = [...]int{

	-1000, -8, -9, -10, -1, 14, 13, -10, 12, 4,
	11, 14, 12, -6, -4, -11, 15, 16, -5, 17,
	4, -7, 14, 6, 5, 18, -12, -11, 10, -2,
	-3, 14, -4, 5, 18, 14, 7, -3, 12, 10,
	4, 11, -11, 12, 14, -6, -7, 4, 5, -6,
	5,
}
var langDef = [...]int{

	0, -2, 1, 0, 4, 0, 0, 0, 2, 0,
	0, 0, 3, 0, 14, 16, 17, 18, 19, 23,
	0, 6, 0, 0, 5, 0, 0, 21, 0, 0,
	0, 0, 15, 20, 0, 13, 7, 0, 8, 0,
	0, 0, 22, 9, 0, 0, 12, 0, 11, 0,
	10,
}
var langTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	4, 5, 3, 3, 18, 3, 10, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 12,
	3, 11, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 8, 3, 9, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 6, 3, 7,
}
var langTok2 = [...]int{

	2, 3, 13, 14, 15, 16, 17,
}
var langTok3 = [...]int{
	0,
}

var langErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	langDebug        = 0
	langErrorVerbose = false
)

type langLexer interface {
	Lex(lval *langSymType) int
	Error(s string)
}

type langParser interface {
	Parse(langLexer) int
	Lookahead() int
}

type langParserImpl struct {
	lookahead func() int
}

func (p *langParserImpl) Lookahead() int {
	return p.lookahead()
}

func langNewParser() langParser {
	p := &langParserImpl{
		lookahead: func() int { return -1 },
	}
	return p
}

const langFlag = -1000

func langTokname(c int) string {
	if c >= 1 && c-1 < len(langToknames) {
		if langToknames[c-1] != "" {
			return langToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func langStatname(s int) string {
	if s >= 0 && s < len(langStatenames) {
		if langStatenames[s] != "" {
			return langStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func langErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !langErrorVerbose {
		return "syntax error"
	}

	for _, e := range langErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + langTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := langPact[state]
	for tok := TOKSTART; tok-1 < len(langToknames); tok++ {
		if n := base + tok; n >= 0 && n < langLast && langChk[langAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if langDef[state] == -2 {
		i := 0
		for langExca[i] != -1 || langExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; langExca[i] >= 0; i += 2 {
			tok := langExca[i]
			if tok < TOKSTART || langExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if langExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += langTokname(tok)
	}
	return res
}

func langlex1(lex langLexer, lval *langSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = langTok1[0]
		goto out
	}
	if char < len(langTok1) {
		token = langTok1[char]
		goto out
	}
	if char >= langPrivate {
		if char < langPrivate+len(langTok2) {
			token = langTok2[char-langPrivate]
			goto out
		}
	}
	for i := 0; i < len(langTok3); i += 2 {
		token = langTok3[i+0]
		if token == char {
			token = langTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = langTok2[1] /* unknown char */
	}
	if langDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", langTokname(token), uint(char))
	}
	return char, token
}

func langParse(langlex langLexer) int {
	return langNewParser().Parse(langlex)
}

func (langrcvr *langParserImpl) Parse(langlex langLexer) int {
	var langn int
	var langlval langSymType
	var langVAL langSymType
	var langDollar []langSymType
	_ = langDollar // silence set and not used
	langS := make([]langSymType, langMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	langstate := 0
	langchar := -1
	langtoken := -1 // langchar translated into internal numbering
	langrcvr.lookahead = func() int { return langchar }
	defer func() {
		// Make sure we report no lookahead when not parsing.
		langstate = -1
		langchar = -1
		langtoken = -1
	}()
	langp := -1
	goto langstack

ret0:
	return 0

ret1:
	return 1

langstack:
	/* put a state and value onto the stack */
	if langDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", langTokname(langtoken), langStatname(langstate))
	}

	langp++
	if langp >= len(langS) {
		nyys := make([]langSymType, len(langS)*2)
		copy(nyys, langS)
		langS = nyys
	}
	langS[langp] = langVAL
	langS[langp].yys = langstate

langnewstate:
	langn = langPact[langstate]
	if langn <= langFlag {
		goto langdefault /* simple state */
	}
	if langchar < 0 {
		langchar, langtoken = langlex1(langlex, &langlval)
	}
	langn += langtoken
	if langn < 0 || langn >= langLast {
		goto langdefault
	}
	langn = langAct[langn]
	if langChk[langn] == langtoken { /* valid shift */
		langchar = -1
		langtoken = -1
		langVAL = langlval
		langstate = langn
		if Errflag > 0 {
			Errflag--
		}
		goto langstack
	}

langdefault:
	/* default state action */
	langn = langDef[langstate]
	if langn == -2 {
		if langchar < 0 {
			langchar, langtoken = langlex1(langlex, &langlval)
		}

		/* look through exception table */
		xi := 0
		for {
			if langExca[xi+0] == -1 && langExca[xi+1] == langstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			langn = langExca[xi+0]
			if langn < 0 || langn == langtoken {
				break
			}
		}
		langn = langExca[xi+1]
		if langn < 0 {
			goto ret0
		}
	}
	if langn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			langlex.Error(langErrorMessage(langstate, langtoken))
			Nerrs++
			if langDebug >= 1 {
				__yyfmt__.Printf("%s", langStatname(langstate))
				__yyfmt__.Printf(" saw %s\n", langTokname(langtoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for langp >= 0 {
				langn = langPact[langS[langp].yys] + langErrCode
				if langn >= 0 && langn < langLast {
					langstate = langAct[langn] /* simulate a shift of "error" */
					if langChk[langstate] == langErrCode {
						goto langstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if langDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", langS[langp].yys)
				}
				langp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if langDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", langTokname(langtoken))
			}
			if langtoken == langEofCode {
				goto ret1
			}
			langchar = -1
			langtoken = -1
			goto langnewstate /* try again in the same state */
		}
	}

	/* reduction by production langn */
	if langDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", langn, langStatname(langstate))
	}

	langnt := langn
	langpt := langp
	_ = langpt // guard against "declared and not used"

	langp -= langR2[langn]
	// langp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if langp+1 >= len(langS) {
		nyys := make([]langSymType, len(langS)*2)
		copy(nyys, langS)
		langS = nyys
	}
	langVAL = langS[langp+1]

	/* consult goto table to find next state */
	langn = langR1[langn]
	langg := langPgo[langn]
	langj := langg + langS[langp].yys + 1

	if langj >= langLast {
		langstate = langAct[langg]
	} else {
		langstate = langAct[langj]
		if langChk[langstate] != -langn {
			langstate = langAct[langg]
		}
	}
	// dummy call; replaced with literal code
	switch langnt {

	case 1:
		langDollar = langS[langpt-1 : langpt+1]
		//line lang.y:50
		{
			langVAL.headers = langDollar[1].headers
		}
	case 2:
		langDollar = langS[langpt-2 : langpt+1]
		//line lang.y:55
		{
			langVAL.headers = []header{langDollar[1].headerField}
		}
	case 3:
		langDollar = langS[langpt-3 : langpt+1]
		//line lang.y:59
		{
			langVAL.headers = append(langDollar[1].headers, langDollar[2].headerField)
		}
	case 4:
		langDollar = langS[langpt-1 : langpt+1]
		//line lang.y:64
		{
			stm = true
			langVAL.headerField = langDollar[1].stream
		}
	case 5:
		langDollar = langS[langpt-4 : langpt+1]
		//line lang.y:69
		{
			stm = true
			langVAL.headerField = &functionCall{label: langDollar[1].str, arguments: langDollar[3].expressions}
		}
	case 6:
		langDollar = langS[langpt-3 : langpt+1]
		//line lang.y:74
		{
			stm = true
			langVAL.headerField = &instrumentInstance{label: langDollar[1].str, inst: langDollar[3].instumentDef}
		}
	case 7:
		langDollar = langS[langpt-5 : langpt+1]
		//line lang.y:81
		{
			langVAL.stream = &astStream{label: langDollar[2].str, instructions: langDollar[4].instructions}
		}
	case 8:
		langDollar = langS[langpt-2 : langpt+1]
		//line lang.y:87
		{
			langVAL.instructions = []instruction{langDollar[1].inst}
		}
	case 9:
		langDollar = langS[langpt-3 : langpt+1]
		//line lang.y:91
		{
			langVAL.instructions = append(langDollar[1].instructions, langDollar[2].inst)
		}
	case 10:
		langDollar = langS[langpt-6 : langpt+1]
		//line lang.y:97
		{
			stm = true
			langVAL.inst = &methodCall{obj: &object{label: langDollar[1].str}, method: langDollar[3].str, arguments: langDollar[5].expressions}
		}
	case 11:
		langDollar = langS[langpt-4 : langpt+1]
		//line lang.y:102
		{
			stm = true
			langVAL.inst = &functionCall{label: langDollar[1].str, arguments: langDollar[3].expressions}
		}
	case 12:
		langDollar = langS[langpt-3 : langpt+1]
		//line lang.y:107
		{
			stm = true
			langVAL.inst = &instrumentInstance{label: langDollar[1].str, inst: langDollar[3].instumentDef}
		}
	case 13:
		langDollar = langS[langpt-3 : langpt+1]
		//line lang.y:114
		{
			i, err := instrumentLookup(langDollar[1].str, langDollar[3].str)
			if err != nil {
				log.Fatalf("Instrument lookup with (%s,%s) failed: %s\n", langDollar[1].str, langDollar[3].str, err)
			}
			langVAL.instumentDef = i
		}
	case 14:
		langDollar = langS[langpt-1 : langpt+1]
		//line lang.y:124
		{
			langVAL.expressions = []expression{langDollar[1].expr}
		}
	case 15:
		langDollar = langS[langpt-3 : langpt+1]
		//line lang.y:128
		{
			langVAL.expressions = append(langDollar[1].expressions, langDollar[3].expr)
		}
	case 16:
		langDollar = langS[langpt-1 : langpt+1]
		//line lang.y:133
		{
			langVAL.expr = langDollar[1].note
		}
	case 17:
		langDollar = langS[langpt-1 : langpt+1]
		//line lang.y:134
		{
			langVAL.expr = intExp(langDollar[1].integer)
		}
	case 18:
		langDollar = langS[langpt-1 : langpt+1]
		//line lang.y:135
		{
			langVAL.expr = langDollar[1].expr
		}
	case 19:
		langDollar = langS[langpt-1 : langpt+1]
		//line lang.y:136
		{
			langVAL.expr = langDollar[1].expr
		}
	case 20:
		langDollar = langS[langpt-3 : langpt+1]
		//line lang.y:139
		{
			langVAL.expr = &chordExpression{notes: langDollar[2].notes}
		}
	case 21:
		langDollar = langS[langpt-1 : langpt+1]
		//line lang.y:142
		{
			langVAL.notes = []*noteExpression{langDollar[1].note}
		}
	case 22:
		langDollar = langS[langpt-3 : langpt+1]
		//line lang.y:143
		{
			langVAL.notes = append(langDollar[1].notes, langDollar[3].note)
		}
	case 23:
		langDollar = langS[langpt-1 : langpt+1]
		//line lang.y:147
		{
			langVAL.note = langDollar[1].note
		}
	}
	goto langstack /* stack new state and value */
}
