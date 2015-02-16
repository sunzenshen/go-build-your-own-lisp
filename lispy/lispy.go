package lispy

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	mpc "github.com/sunzenshen/golispy/mpcinterface"
)

// Lispy is a collection of the Lispy parser definitions
type Lispy struct {
	numberParser, operatorParser, exprParser, lispyParser mpc.MpcParser
}

// CleanLispy is used after parsers initiated by InitLispy are not longer to be used
func CleanLispy(l Lispy) {
	mpc.MpcCleanup(l.numberParser, l.operatorParser, l.exprParser, l.lispyParser)
}

// Eval translates an AST into the final result of the represented instructions
func Eval(tree mpc.MpcAst) int64 {
	if strings.Contains(mpc.GetTag(tree), "number") {
		num, _ := strconv.ParseInt(mpc.GetContents(tree), 10, 0)
		return num
	}
	op := mpc.GetOperator(tree)
	x := Eval(mpc.GetChild(tree, 2))
	i := 3
	for strings.Contains(mpc.GetTag(mpc.GetChild(tree, i)), "expr") {
		x = evalOp(x, op, Eval(mpc.GetChild(tree, i)))
		i++
	}
	return x
}

func evalOp(x int64, op string, y int64) int64 {
	if strings.Contains(op, "+") || strings.Contains(op, "add") {
		return x + y
	}
	if strings.Contains(op, "-") || strings.Contains(op, "sub") {
		return x - y
	}
	if strings.Contains(op, "*") || strings.Contains(op, "mul") {
		return x * y
	}
	if strings.Contains(op, "/") || strings.Contains(op, "div") {
		return x / y
	}
	if strings.Contains(op, "%") || strings.Contains(op, "mod") {
		return x % y
	}
	if strings.Contains(op, "^") || strings.Contains(op, "pow") {
		z := big.NewInt(0)
		z.Exp(big.NewInt(x), big.NewInt(y), nil)
		return z.Int64()
	}
	return 0
}

// InitLispy returns the parsers for the Lispy language definition
func InitLispy() Lispy {
	number := mpc.MpcNew("number")
	operator := mpc.MpcNew("operator")
	expr := mpc.MpcNew("expr")
	lispy := mpc.MpcNew("lispy")
	language := "" +
		"number : /-?[0-9]+/                                                  ; " +
		"operator :   '+'   |   '-'   |   '*'   |   '/'   |   '%'   |   '^'     " +
		"         | \"add\" | \"sub\" | \"mul\" | \"div\" | \"mod\" | \"pow\" ; " +
		"expr     : <number> | '(' <operator> <expr>+ ')'                     ; " +
		"lispy    : /^/ <operator> <expr>+ /$/                                ; "
	mpc.MpcaLang(language, number, operator, expr, lispy)
	parserSet := Lispy{}
	parserSet.numberParser = number
	parserSet.operatorParser = operator
	parserSet.exprParser = expr
	parserSet.lispyParser = lispy
	return parserSet
}

// PrintAst prints the AST of a Lispy expression.
func (l *Lispy) PrintAst(input string) {
	mpc.PrintAst(input, l.lispyParser)
}

// ReadEval takes a string, tries to interpret it in Lispy
func (l *Lispy) ReadEval(input string, printErrors bool) (int64, error) {
	r, err := mpc.MpcParse(input, l.lispyParser)
	if err != nil {
		if printErrors {
			mpc.MpcErrPrint(&r)
		}
		mpc.MpcErrDelete(&r)
		return 0, errors.New("mpc: ReadEval call to MpcParse failed")
	}
	defer mpc.MpcAstDelete(&r)
	return Eval(mpc.GetOutput(&r)), nil
}

// ReadEvalPrint takes a string, tries to interpret it in Lispy, or returns an parsing error
func (l *Lispy) ReadEvalPrint(input string) {
	evalResult, err := l.ReadEval(input, true)
	if err == nil {
		fmt.Println(evalResult)
	}
}
