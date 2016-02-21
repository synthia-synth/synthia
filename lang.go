//line lang.y:2
package main

import __yyfmt__ "fmt"

//line lang.y:2
//line lang.y:6
type langSymType struct {
	yys          int
	stream       *astStream
	str          string
	instructions []instruction
	inst         instruction
	expr         expression
	expressions  []expression
	integer      int
}

const STREAM = 57346
const LABEL = 57347
const NUM = 57348
const NOTE = 57349

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
	"STREAM",
	"LABEL",
	"NUM",
	"NOTE",
}
var langStatenames = [...]string{}

const langEofCode = 1
const langErrCode = 2
const langInitialStackSize = 16

//line lang.y:66

//line yacctab:1
var langExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const langNprod = 8
const langPrivate = 57344

var langTokenNames []string
var langStates []string

const langLast = 22

var langAct = [...]int{

	16, 15, 14, 19, 11, 8, 7, 10, 3, 15,
	7, 2, 20, 18, 4, 6, 17, 12, 13, 5,
	1, 9,
}
var langPact = [...]int{

	0, -1000, -4, 8, -6, -2, -1000, -3, -1000, -1000,
	-8, 13, -13, -5, -1000, 5, -1000, -1000, -10, 3,
	-1000,
}
var langPgo = [...]int{

	0, 20, 19, 15, 2, 18,
}
var langR1 = [...]int{

	0, 1, 2, 2, 3, 5, 5, 4,
}
var langR2 = [...]int{

	0, 5, 1, 2, 6, 1, 2, 4,
}
var langChk = [...]int{

	-1000, -1, 11, 12, 6, -2, -3, 12, 7, -3,
	10, 12, 4, -5, -4, 14, 5, -4, 8, 13,
	9,
}
var langDef = [...]int{

	0, -2, 0, 0, 0, 0, 2, 0, 1, 3,
	0, 0, 0, 0, 5, 0, 4, 6, 0, 0,
	7,
}
var langTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	4, 5, 3, 3, 3, 3, 10, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 8, 3, 9, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 6, 3, 7,
}
var langTok2 = [...]int{

	2, 3, 11, 12, 13, 14,
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
	lval  langSymType
	stack [langInitialStackSize]langSymType
	char  int
}

func (p *langParserImpl) Lookahead() int {
	return p.char
}

func langNewParser() langParser {
	return &langParserImpl{}
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
	var langVAL langSymType
	var langDollar []langSymType
	_ = langDollar // silence set and not used
	langS := langrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	langstate := 0
	langrcvr.char = -1
	langtoken := -1 // langrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		langstate = -1
		langrcvr.char = -1
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
	if langrcvr.char < 0 {
		langrcvr.char, langtoken = langlex1(langlex, &langrcvr.lval)
	}
	langn += langtoken
	if langn < 0 || langn >= langLast {
		goto langdefault
	}
	langn = langAct[langn]
	if langChk[langn] == langtoken { /* valid shift */
		langrcvr.char = -1
		langtoken = -1
		langVAL = langrcvr.lval
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
		if langrcvr.char < 0 {
			langrcvr.char, langtoken = langlex1(langlex, &langrcvr.lval)
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
			langrcvr.char = -1
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
		langDollar = langS[langpt-5 : langpt+1]
		//line lang.y:31
		{
			langVAL.stream = &astStream{label: langDollar[2].str, instructions: langDollar[4].instructions}
		}
	case 2:
		langDollar = langS[langpt-1 : langpt+1]
		//line lang.y:37
		{
			langVAL.instructions = []instruction{langDollar[1].inst}
		}
	case 3:
		langDollar = langS[langpt-2 : langpt+1]
		//line lang.y:41
		{
			langVAL.instructions = append(langDollar[1].instructions, langDollar[2].inst)
		}
	case 4:
		langDollar = langS[langpt-6 : langpt+1]
		//line lang.y:47
		{
			langVAL.inst = &methodCall{obj: &object{label: langDollar[1].str}, method: langDollar[3].str, arguments: langDollar[5].expressions}
		}
	case 5:
		langDollar = langS[langpt-1 : langpt+1]
		//line lang.y:53
		{
			langVAL.expressions = []expression{langDollar[1].expr}
		}
	case 6:
		langDollar = langS[langpt-2 : langpt+1]
		//line lang.y:57
		{
			langVAL.expressions = append(langDollar[1].expressions, langDollar[2].expr)
		}
	case 7:
		langDollar = langS[langpt-4 : langpt+1]
		//line lang.y:63
		{
			langVAL.expr = &noteExpression{note: NoteName(langDollar[1].integer), octave: langDollar[3].integer}
		}
	}
	goto langstack /* stack new state and value */
}
