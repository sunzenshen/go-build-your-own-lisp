package lispy

import "fmt"
import "github.com/sunzenshen/go-build-your-own-lisp/mpc"

// Lispy is a collection of the Lispy parser definitions
type Lispy struct {
	env           *lenv
	numberParser  mpc.MpcParser
	symbolParser  mpc.MpcParser
	strParser     mpc.MpcParser
	commentParser mpc.MpcParser
	sexprParser   mpc.MpcParser
	qexprParser   mpc.MpcParser
	exprParser    mpc.MpcParser
	lispyParser   mpc.MpcParser
}

// CleanLispy is used after parsers initiated by InitLispy are not longer to be used
func CleanLispy(l Lispy) {
	mpc.MpcCleanup(
		l.numberParser,
		l.symbolParser,
		l.strParser,
		l.commentParser,
		l.sexprParser,
		l.qexprParser,
		l.exprParser,
		l.lispyParser)
}

// InitLispy returns the parsers for the Lispy language definition
func InitLispy() Lispy {
	number := mpc.MpcNew("number")
	symbol := mpc.MpcNew("symbol")
	str := mpc.MpcNew("string")
	comment := mpc.MpcNew("comment")
	sexpr := mpc.MpcNew("sexpr")
	qexpr := mpc.MpcNew("qexpr")
	expr := mpc.MpcNew("expr")
	lispy := mpc.MpcNew("lispy")
	language := "" +
		"number  : /-?[0-9]+/                                                     ; " +
		"symbol  : /[a-zA-Z0-9_+\\-*\\/\\\\=<>!&]+/                               ; " +
		"string  : /\"(\\\\.|[^\"])*\"/                                           ; " +
		"comment : /;[^\\r\\n]*/                                                  ; " +
		"sexpr   : '(' <expr>* ')'                                                ; " +
		"qexpr   : '{' <expr>* '}'                                                ; " +
		"expr    : <number> | <symbol> | <string> | <comment> | <sexpr> | <qexpr> ; " +
		"lispy   : /^/ <expr>* /$/                                                ; "
	mpc.MpcaLang(language, number, symbol, str, comment, sexpr, qexpr, expr, lispy)
	l := Lispy{}
	l.numberParser = number
	l.symbolParser = symbol
	l.strParser = str
	l.commentParser = comment
	l.sexprParser = sexpr
	l.qexprParser = qexpr
	l.exprParser = expr
	l.lispyParser = lispy
	// Init environment
	l.env = lenvNew()
	l.env.lenvAddBuiltins()
	l.env.parser = lispy // For loading files with builtin
	return l
}

// PrintAst prints the AST of a Lispy expression.
func (l *Lispy) PrintAst(input string) {
	mpc.PrintAst(input, l.lispyParser)
}

// Read takes a string and parses it into an lval
func (l *Lispy) Read(input string, printErrors bool) *lval {
	r, err := mpc.MpcParse(input, l.lispyParser)
	if err != nil {
		if printErrors {
			mpc.MpcErrPrint(&r)
		}
		mpc.MpcErrDelete(&r)
		return lvalErr("Failed to parse input: '%s'", input)
	}
	defer mpc.DeleteAstPtr(&r)
	return lvalRead(mpc.GetOutput(&r))
}

// Eval translates an lval into the final result of the represented instructions
func (v *lval) Eval(e *lenv) *lval {
	return v.lvalEval(e)
}

// ReadEval takes a string, tries to interpret it in Lispy
func (l *Lispy) ReadEval(input string, printErrors bool) *lval {
	return l.Read(input, printErrors).Eval(l.env)
}

// ReadEvalPrint takes a string, tries to interpret it in Lispy, and prints the result
func (l *Lispy) ReadEvalPrint(input string) {
	l.ReadEval(input, true).lvalPrintLn()
}

// LoadFiles loads a list of files into the Lispy environment
func (l *Lispy) LoadFiles(files []string) {
	for _, file := range files {
		fmt.Println("Loading: ", file)
		// Argument list with a single argument, the file name
		args := lvalAdd(lvalSexpr(), lvalStr(file))
		// Pass into builtin load to get the result
		x := builtinLoad(l.env, args)
		// If the result is an error, print it
		if x.ltype == lvalErrType {
			x.lvalPrintLn()
		}
	}
}
