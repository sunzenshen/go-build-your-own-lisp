package lispy

import "fmt"
import "github.com/sunzenshen/go-build-your-own-lisp/mpc"

type lbuiltin func(*lenv, *lval) *lval

func builtinOp(e *lenv, a *lval, op string) *lval {
	// Ensure all arguments are numbers
	for _, cell := range a.cells {
		if cell.ltype != lvalNumType {
			return lvalErr("Cannot operate on non-number: %s", cell.ltypeName())
		}
	}
	// Pop the first element
	x := a.lvalPop(0)
	// Handle unary negation
	if op == "-" && a.cellCount() == 0 {
		x.num = -x.num
	}
	// Process remaining elements
	for a.cellCount() > 0 {
		// Pop the next element
		y := a.lvalPop(0)
		// Perform symbol's operation
		if op == "+" {
			x.num += y.num
		} else if op == "-" {
			x.num -= y.num
		} else if op == "*" {
			x.num *= y.num
		} else if op == "/" {
			if y.num == 0 {
				x = lvalErr("Division By Zero!")
			} else {
				x.num /= y.num
			}
		}
	}
	return x
}

func builtinDef(e *lenv, a *lval) *lval {
	return builtinVar(e, a, "def")
}

func builtinPut(e *lenv, a *lval) *lval {
	return builtinVar(e, a, "=")
}

func builtinVar(e *lenv, a *lval, function string) *lval {
	if a.cells[0].ltype != lvalQexprType {
		return lvalErr("Function %s passed incorrect type: %s", function, a.cells[0].ltypeName())
	}
	// First argument is symbol list
	syms := a.cells[0]
	// Ensure elements of first list are symbols
	for _, cell := range syms.cells {
		if cell.ltype != lvalSymType {
			return lvalErr("Function %s cannot define non-symbol: %s", function, cell.ltypeName())
		}
	}
	// Check for the correct number of symbols and values
	if syms.cellCount() != a.cellCount()-1 {
		return lvalErr("Function %s cannot define incorrect number of values to symbols", function)
	}
	// Assign copies of values to symbols
	for i, cell := range syms.cells {
		// 'def' to define globally
		if function == "def" {
			e.lenvDef(cell, a.cells[i+1])
		}
		// 'put' to define locally
		if function == "=" {
			e.lenvPut(cell, a.cells[i+1])
		}
	}
	return lvalSexpr()
}

func builtinHead(e *lenv, a *lval) *lval {
	// Check for error conditions
	if a.cellCount() != 1 {
		return lvalErr("Function 'head' passed too many arguments: %s", a.lvalString())
	}
	if a.cells[0].ltype != lvalQexprType {
		return lvalErr("Function 'head' passed incorrect types: %s", a.lvalString())
	}
	if a.cells[0].cellCount() == 0 {
		return lvalErr("Function 'head' passed {}!")
	}
	// Otherwise, get the head
	v := a.lvalTake(0)
	for v.cellCount() > 1 {
		v.lvalPop(1)
	}
	return v
}

func builtinTail(e *lenv, a *lval) *lval {
	// Check for error conditions
	if a.cellCount() != 1 {
		return lvalErr("Function 'tail' passed too many arguments: %s", a.lvalString())
	}
	if a.cells[0].ltype != lvalQexprType {
		return lvalErr("Function 'tail' passed incorrect types: %s", a.lvalString())
	}
	if a.cells[0].cellCount() == 0 {
		return lvalErr("Function 'tail' passed {}!")
	}
	// Otherwise, get the tail
	v := a.lvalTake(0)
	v.lvalPop(0)
	return v
}

func builtinList(e *lenv, a *lval) *lval {
	a.ltype = lvalQexprType
	return a
}

func builtinEval(e *lenv, a *lval) *lval {
	if a.cellCount() != 1 {
		return lvalErr("Function 'eval' passed too many arguments: %s", a.lvalString())
	}
	if a.cells[0].ltype != lvalQexprType {
		return lvalErr("Function 'eval' passed incorrect type: %s", a.lvalString())
	}
	x := a.lvalTake(0)
	x.ltype = lvalSexprType
	return x.lvalEval(e)
}

func builtinJoin(e *lenv, a *lval) *lval {
	for _, cell := range a.cells {
		if cell.ltype != lvalQexprType {
			return lvalErr("Function 'join' passed incorrect type: %s", a.lvalString())
		}
	}
	x := a.lvalPop(0)
	for a.cellCount() > 0 {
		x = lvalJoin(x, a.lvalPop(0))
	}
	return x
}

func builtinLambda(e *lenv, a *lval) *lval {
	if a.cellCount() != 2 {
		return lvalErr("Lambda has %d arguments, not 2 as expected", a.cellCount())
	} else if a.cells[0].ltype != lvalQexprType {
		return lvalErr("Lambda cell[0] has unexpected type %d", a.cells[0].ltype)
	} else if a.cells[1].ltype != lvalQexprType {
		return lvalErr("Lambda cell[1] has unexpected type %d", a.cells[1].ltype)
	}
	// Check that the first Q-expression contains only Symbols
	for _, cell := range a.cells[0].cells {
		if cell.ltype != lvalSymType {
			return lvalErr("Cannot define non-symbol. Got type %s instead", cell.ltype)
		}
	}
	// Pop first 2 arguments and pass them to lvalLambda
	formals := a.lvalPop(0)
	body := a.lvalPop(0)
	return lvalLambda(formals, body)
}

func builtinGreaterThan(e *lenv, a *lval) *lval {
	return builtinOrd(e, a, ">")
}

func builtinLessThan(e *lenv, a *lval) *lval {
	return builtinOrd(e, a, "<")
}

func builtinGreaterEqual(e *lenv, a *lval) *lval {
	return builtinOrd(e, a, ">=")
}

func builtinLessEqual(e *lenv, a *lval) *lval {
	return builtinOrd(e, a, "<=")
}

func builtinOrd(e *lenv, a *lval, op string) *lval {
	if a.cellCount() != 2 {
		return lvalErr("%s passed in with %d cells not 2", op, a.cellCount())
	}
	if a.cells[0].ltype != lvalNumType {
		return lvalErr("%s cell0 is not a number, but type %s", op, a.cells[0].ltypeName())
	}
	if a.cells[1].ltype != lvalNumType {
		return lvalErr("%s cell1 is not a number, but type %s", op, a.cells[1].ltypeName())
	}
	var cmp bool
	if op == ">" {
		cmp = a.cells[0].num > a.cells[1].num
	} else if op == "<" {
		cmp = a.cells[0].num < a.cells[1].num
	} else if op == ">=" {
		cmp = a.cells[0].num >= a.cells[1].num
	} else if op == "<=" {
		cmp = a.cells[0].num <= a.cells[1].num
	}
	// 0 = false, and everything else is true
	if cmp {
		return lvalNum(1)
	}
	return lvalNum(0)
}

func builtinCmp(e *lenv, a *lval, op string) *lval {
	if a.cellCount() != 2 {
		return lvalErr("%s passed in with %d args, not 2", op, a.cellCount())
	}
	var cmp bool
	if op == "==" {
		cmp = lvalEq(a.cells[0], a.cells[1])
	} else if op == "!=" {
		cmp = !lvalEq(a.cells[0], a.cells[1])
	}
	if cmp {
		return lvalNum(1)
	}
	return lvalNum(0)
}

func builtinEqual(e *lenv, a *lval) *lval {
	return builtinCmp(e, a, "==")
}

func builtinNotEqual(e *lenv, a *lval) *lval {
	return builtinCmp(e, a, "!=")
}

func builtinIf(e *lenv, a *lval) *lval {
	if a.cellCount() != 3 {
		return lvalErr("if passed in %d args, not 3", a.cellCount())
	}
	if a.cells[0].ltype != lvalNumType {
		return lvalErr("if cell0 is not a number")
	}
	if a.cells[1].ltype != lvalQexprType {
		return lvalErr("if cell1 is not a Q-exp")
	}
	if a.cells[2].ltype != lvalQexprType {
		return lvalErr("if cell2 is not a Q-exp")
	}
	// Mark both expressions as evaluable
	a.cells[1].ltype = lvalSexprType
	a.cells[2].ltype = lvalSexprType
	// Determine branch direction
	var x *lval
	if a.cells[0].num != 0 {
		// If condition is true, evaluate the first expression
		x = a.lvalPop(1).lvalEval(e)
	} else {
		// Otherwise, evaluate the second expression
		x = a.lvalPop(2).lvalEval(e)
	}
	return x
}

func builtinLoad(e *lenv, a *lval) *lval {
	if a.cellCount() != 1 {
		return lvalErr("load passed in a with %d cells, expected 1", a.cellCount())
	}
	if a.cells[0].ltype != lvalStrType {
		return lvalErr("load did not get a string for input")
	}
	if e.parser == nil {
		return lvalErr("load env is missing parser")
	}
	// Parse string as a file name
	var ret *lval
	r, err := mpc.MpcParseContents(a.cells[0].str, e.parser)
	if err != nil {
		// Get parse error in string format
		errMsg := mpc.GetErrorStr(r)
		ret = lvalErr("Could not load library %s", errMsg)
	} else if r != nil {
		// Read contents
		expr := lvalRead(mpc.GetOutput(r))
		mpc.DeleteAstPtr(r)
		// Evaluate each expression
		for expr.cellCount() > 0 {
			x := expr.lvalPop(0).lvalEval(e)
			// If evaluation leads to an error, print it
			if x.ltype == lvalErrType {
				x.lvalPrintLn()
			}
		}
		// Return an empty list
		ret = lvalSexpr()
	}

	return ret
}

func builtinPrint(e *lenv, a *lval) *lval {
	// Print each argument followed by a space
	for _, cell := range a.cells {
		cell.lvalPrint()
		fmt.Println(" ")
	}
	fmt.Println("")
	return lvalSexpr()
}

func builtinError(e *lenv, a *lval) *lval {
	if a.cellCount() != 1 {
		return lvalErr("error: %d parameters were passed in, expected only 1", a.cellCount())
	}
	if a.cells[0].ltype != lvalStrType {
		return lvalErr("error: argument is not a string")
	}
	// Construct error from the first argument
	return lvalErr(a.cells[0].str)
}

func builtinAdd(e *lenv, a *lval) *lval {
	return builtinOp(e, a, "+")
}

func builtinSub(e *lenv, a *lval) *lval {
	return builtinOp(e, a, "-")
}

func builtinMul(e *lenv, a *lval) *lval {
	return builtinOp(e, a, "*")
}

func builtinDiv(e *lenv, a *lval) *lval {
	return builtinOp(e, a, "/")
}
