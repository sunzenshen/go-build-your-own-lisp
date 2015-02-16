package mpcinterface

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Lispy is a collection of the Lispy parser definitions
type Lispy struct {
	numberParser, operatorParser, exprParser, lispyParser MpcParser
}

// CleanLispy is used after parsers initiated by InitLispy are not longer to be used
func CleanLispy(l Lispy) {
	MpcCleanup(l.numberParser, l.operatorParser, l.exprParser, l.lispyParser)
}

// Eval translates an AST into the final result of the represented instructions
func Eval(tree MpcAst) int64 {
	if strings.Contains(getTag(tree), "number") {
		num, _ := strconv.ParseInt(getContents(tree), 10, 0)
		return num
	}
	op := getOperator(tree)
	x := Eval(getChild(tree, 2))
	i := 3
	for strings.Contains(getTag(getChild(tree, i)), "expr") {
		x = evalOp(x, op, Eval(getChild(tree, i)))
		i++
	}
	return x
}

func evalOp(x int64, op string, y int64) int64 {
	if strings.Contains(op, "+") {
		return x + y
	}
	if strings.Contains(op, "-") {
		return x - y
	}
	if strings.Contains(op, "*") {
		return x * y
	}
	if strings.Contains(op, "/") {
		return x / y
	}
	if strings.Contains(op, "%") {
		return x % y
	}
	return 0
}

// InitLispy returns the parsers for the Lispy language definition
func InitLispy() Lispy {
	number := mpcNew("number")
	operator := mpcNew("operator")
	expr := mpcNew("expr")
	lispy := mpcNew("lispy")
	language := "" +
		"number : /-?[0-9]+/                               ; " +
		"operator : '+' | '-' | '*' | '/' | '%'            ; " +
		"expr     : <number> | '(' <operator> <expr>+ ')'  ; " +
		"lispy    : /^/ <operator> <expr>+ /$/             ; "
	MpcaLang(language, number, operator, expr, lispy)
	parserSet := Lispy{}
	parserSet.numberParser = number
	parserSet.operatorParser = operator
	parserSet.exprParser = expr
	parserSet.lispyParser = lispy
	return parserSet
}

// ReadEval takes a string, tries to interpret it in Lispy
func (l *Lispy) ReadEval(input string, printErrors bool) (int64, error) {
	r, err := MpcParse(input, l.lispyParser)
	if err != nil {
		if printErrors {
			MpcErrPrint(&r)
		}
		MpcErrDelete(&r)
		return 0, errors.New("mpc: ReadEval call to MpcParse failed")
	}
	defer MpcAstDelete(&r)
	return Eval(GetOutput(&r)), nil
}

// ReadEvalPrint takes a string, tries to interpret it in Lispy, or returns an parsing error
func (l *Lispy) ReadEvalPrint(input string) {
	evalResult, err := l.ReadEval(input, true)
	if err == nil {
		fmt.Println(evalResult)
	}
}
