%{
package main

%}

%union {
	stream *astStream
	str string
	instructions []instruction
	inst instruction
	expr expression
	expressions []expression
	integer int
}

%token '(' ')' '{' '}' '[' ']' '.'
%token STREAM
%token <str> LABEL
%token <integer> NUM NOTE

%type <stream> strm
%type <instructions> inst_list
%type <inst> inst
%type <expr> expr
%type <expressions> expr_list

%%

strm:
	STREAM LABEL '{' inst_list '}'
	{
		$$ = &astStream{label: $2, instructions: $4}
	}

inst_list:
	inst
	{
		$$ = []instruction{$1}
	}
	| inst_list inst
	{
		$$ = append($1, $2)
	}

inst:
	LABEL '.' LABEL '(' expr_list ')'
	{
		$$ = &methodCall{obj: &object{label: $1}, method: $3, arguments: $5}
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
%%
