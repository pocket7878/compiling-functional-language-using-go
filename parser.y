%{

package main

import (
    "fmt"
    "io"
    "text/scanner"
    "unicode"
    "strconv"
)

type (
	item struct {
		typ int
		val string
	}
)
%}

%token PLUS
%token TIMES
%token MINUS
%token DIVIDE
%token <number> INT
%token DEFN
%token DATA
%token CASE
%token OF
%token OCURLY
%token CCURLY
%token OPAREN
%token CPAREN
%token COMMA
%token ARROW
%token EQUAL
%token <lid> LID
%token <uid> UID

%type <params> lowercaseParams uppercaseParams
%type <definitions> program definitions
%type <branches> branches
%type <constructors> constructors
%type <ast> aAdd aMul case app appBase
%type <definition> definition defn data 
%type <branch> branch
%type <pattern> pattern
%type <constructor> constructor

%union {
	Token item
    number int
    params []string
    definition definition
    definitions []definition
    branch branch
    branches []branch
    pattern pattern
    constructor constructor
    constructors []constructor
    ast ast
    lid string
    uid string
}

%start program

%%

program
    : definitions 
    { 
        $$ = $1
        yylex.(*lexer).result = $$
    }
    ;

definitions
    : definitions definition { $$ = $1; $$ = append($$, $2); }
    | definition { $$ = make([]definition, 0); $$ = append($$, $1); }
    ;

definition
    : defn { $$ = $1; }
    | data { $$ = $1; }
    ;

defn
    : DEFN LID lowercaseParams EQUAL OCURLY aAdd CCURLY
        { $$ = definitionDefn{$2, $3, $6} }
    ;

lowercaseParams 
    : { $$ = make([]string, 0); }
    | lowercaseParams LID { $$ = $1; $$ = append($$, $2); }
    ;

uppercaseParams
    : { $$ = make([]string, 0); }
    | uppercaseParams UID { $$ = $1; $$ = append($$, $2); }
    ;

aAdd
    : aAdd PLUS aMul { $$ = astBinOp{binOpPlus, $1, $3}; }
    | aAdd MINUS aMul { $$ = astBinOp{binOpMinus, $1, $3}; }
    | aMul { $$ = $1; }
    ;

aMul
    : aAdd TIMES aMul { $$ = astBinOp{binOpTimes, $1, $3}; }
    | aAdd DIVIDE aMul { $$ = astBinOp{binOpDivide, $1, $3}; }
    | app { $$ = $1; }
    ;

app
    : app appBase { $$ = astApp{$1, $2}; }
    | appBase { $$ = $1; }
    ;

appBase
    : INT { $$ = astInt{$1}; fmt.Println("astInt: %v", $1) }
    | LID { $$ = astLID{$1}; }
    | UID { $$ = astUID{$1}; }
    | OPAREN aAdd CPAREN { $$ = $2; }
    | case { $$ = $1; }
    ;

case
    : CASE aAdd OF OCURLY branches CCURLY 
        { $$ = astCase{$2, $5}; }
    ;

branches
    : branches branch { $$ = $1; $$ = append($$, $2); }
    | branch { $$ = make([]branch, 0); $$ = append($$, $1);}
    ;

branch
    : pattern ARROW OCURLY aAdd CCURLY
        { $$ = branch{$1, $4}; }
    ;

pattern
    : LID { $$ = &patternVar{$1}; }
    | UID lowercaseParams
        { $$ = patternConstr{$1, $2}; }
    ;

data
    : DATA UID EQUAL OCURLY constructors CCURLY
        { $$ = definitionData{$2, $5}; }
    ;

constructors
    : constructors COMMA constructor { $$ = $1; $$ = append($$, $3); }
    | constructor
        { $$ = make([]constructor, 0); $$ = append($$, $1); }
    ;

constructor
    : UID uppercaseParams
        { $$ = constructor{$1, $2}; }
    ;

%%

var simpleTokenTypeTable = map[string]int{
	"+":    PLUS,
	"*":    TIMES,
	"/":    DIVIDE,
	"defn": DEFN,
	"data": DATA,
	"case": CASE,
	"of":   OF,
	"{":    OCURLY,
	"}":    CCURLY,
	"(":    OPAREN,
	")":    CPAREN,
	",":    COMMA,
	"=":    EQUAL,
}


type lexer struct {
	scanner scanner.Scanner
    result []definition
}

func newLexer(reader io.Reader) *lexer {
	var s scanner.Scanner
	s.Init(reader)
	return &lexer{
		s,
        make([]definition, 0),
	}
}

func (l *lexer) Lex(lval *yySymType) int {
    tok := l.scanner.Scan()
    if tok == scanner.EOF {
        return 0
    }

    tokenText := l.scanner.TokenText()
    if tok == scanner.Int {
        lval.number, _ = strconv.Atoi(tokenText)
        return INT
    }

    it, ok := simpleTokenTypeTable[tokenText]
    if ok {
        return it
    }

    if tokenText == "-" {
        nextToken := l.scanner.Peek()
        if nextToken == scanner.EOF {
            return MINUS
        } else if nextToken == '>' {
            l.scanner.Scan()

            return ARROW
        } else {
            return MINUS
        }
    }

    if unicode.IsLower(rune(tokenText[0])) {
        lval.lid = tokenText

        return LID
    } else if unicode.IsUpper(rune(tokenText[0])) {
        lval.uid = tokenText

        return UID
    } else {

        return 0
    }
}

func (l *lexer) Error(e string) {
    panic(e)
}