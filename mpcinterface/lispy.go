package mpcinterface

// #cgo LDFLAGS: -ledit -lm
// #include "mpc_interface.h"
import "C"
import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// CleanLispy is used after parsers initiated by InitLispy are not longer to be used
func CleanLispy(number MpcParser, operator MpcParser, expr MpcParser, lispy MpcParser) {
	MpcCleanup(number, operator, expr, lispy)
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
	return 0
}

// InitLispy returns the parsers for the Lispy language definition
func InitLispy() (MpcParser, MpcParser, MpcParser, MpcParser) {
	Number := mpcNew("number")
	Operator := mpcNew("operator")
	Expr := mpcNew("expr")
	Lispy := mpcNew("lispy")
	language := "" +
		"number : /-?[0-9]+/                               ; " +
		"operator : '+' | '-' | '*' | '/'                  ; " +
		"expr     : <number> | '(' <operator> <expr>+ ')'  ; " +
		"lispy    : /^/ <operator> <expr>+ /$/             ; "
	MpcaLang(language, Number, Operator, Expr, Lispy)
	return Number, Operator, Expr, Lispy
}

// ReadEval takes a string, tries to interpret it in Lispy
func ReadEval(input string, mpcParser MpcParser) (int64, error) {
	r, err := MpcParse(input, mpcParser)
	if err != nil {
		return 0, errors.New("mpc: ReadEval call to MpcParse failed")
	}
	defer MpcAstDelete(&r)
	return Eval(GetOutput(&r)), nil
}

// ReadEvalPrint takes a string, tries to interpret it in Lispy, or returns an parsing error
func ReadEvalPrint(input string, mpcParser MpcParser) {
	r, err := MpcParse(input, mpcParser)
	if err != nil {
		MpcErrPrint(&r)
		C.mpc_err_delete(C.get_error(&r))
	} else {
		fmt.Println(Eval(C.get_output(&r)))
		MpcAstDelete(&r)
	}
}
