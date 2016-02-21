%{

%}

%union {
	stream *astStream
	str string
	instructions []instruction
	inst instruction
}

%token '(' ')' '{' '}'
%token STREAM
%token <str> LABEL

%type <stream> strm
%type <instructions> inst_list

%%

strm:
	STREAM LABEL '{' inst_list '}'
	{
		$$ = {label: $2, instructions: inst_list}
	}

inst_list:
	inst
	{
		$$ = []instructions{$1}
	}
	| inst_list inst
	{
		$$ = append($1, $2)
	}

inst:
	

%%
