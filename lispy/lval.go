package lispy

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sunzenshen/go-build-your-own-lisp/mpc"
)

// ltype values for lval
const (
	lvalNumType = iota
	lvalSymType
	lvalFunType
	lvalSexprType
	lvalQexprType
	lvalErrType
)

type lval struct {
	ltype int

	// Basic
	num int64  // lvalNumType
	err string // lvalErrType
	sym string // lvalSymType

	// Function
	builtin lbuiltin // lvalFunType, nil for user defined function
	env     *lenv
	formals *lval
	body    *lval

	// Expression
	cells []*lval // lvalSexprType, lvalQexprType
}

// lvalNum creates an lval number
func lvalNum(x int64) *lval {
	v := new(lval)
	v.ltype = lvalNumType
	v.num = x
	return v
}

// lvalErr creates an lval error
func lvalErr(f string, a ...interface{}) *lval {
	v := new(lval)
	v.ltype = lvalErrType
	v.err = fmt.Sprintf(f, a...)
	return v
}

// lvalSym creates an lval symbol
func lvalSym(s string) *lval {
	v := new(lval)
	v.ltype = lvalSymType
	v.sym = string(s)
	return v
}

// lvalSexpr creates an lval S-expression
func lvalSexpr() *lval {
	v := new(lval)
	v.ltype = lvalSexprType
	v.cells = make([]*lval, 0)
	return v
}

// lvalQexpr creates an lval Q-expression
func lvalQexpr() *lval {
	v := new(lval)
	v.ltype = lvalQexprType
	v.cells = make([]*lval, 0)
	return v
}

// lvalFun creates a function lval
func lvalFun(function lbuiltin) *lval {
	v := new(lval)
	v.ltype = lvalFunType
	v.builtin = function
	return v
}

// lvalLambda creates a user defined lval function
func lvalLambda(formals *lval, body *lval) *lval {
	v := new(lval)
	v.ltype = lvalFunType
	// Builtin is nil for user defined functions
	v.builtin = nil
	// Init environment
	v.env = lenvNew()
	// Set formals and body
	v.formals = formals
	v.body = body
	return v
}

func (v *lval) cellCount() int {
	return len(v.cells)
}

func (v *lval) ltypeName() string {
	return ltypeName(v.ltype)
}

func ltypeName(i int) string {
	switch i {
	case lvalNumType:
		return "Number"
	case lvalErrType:
		return "Error"
	case lvalSymType:
		return "Symbol"
	case lvalFunType:
		return "Function"
	case lvalSexprType:
		return "S-Expression"
	case lvalQexprType:
		return "Q-Expression"
	}
	return "Unknown:" + strconv.Itoa(i)
}

func (v *lval) lvalString() string {
	switch v.ltype {
	case lvalNumType:
		return strconv.FormatInt(v.num, 10)
	case lvalErrType:
		return ("Error: " + v.err)
	case lvalSymType:
		return (v.sym)
	case lvalFunType:
		if v.builtin == nil {
			return "(\\ " + v.formals.lvalString() + " " + v.body.lvalString() + ")"
		}
		return "<builtin>"
	case lvalSexprType:
		return v.lvalExprString("(", ")")
	case lvalQexprType:
		return v.lvalExprString("{", "}")
	}
	return fmt.Sprintf("Error: lvalString() unhandled ltype %d", v.ltype)
}

func (v *lval) lvalPrint() {
	fmt.Print(v.lvalString())
}

func (v *lval) lvalPrintLn() {
	v.lvalPrint()
	fmt.Print("\n")
}

func lvalCopy(v *lval) *lval {
	x := new(lval)
	x.ltype = v.ltype
	switch v.ltype {
	case lvalFunType:
		if v.builtin == nil {
			x.builtin = nil
			x.env = lenvCopy(v.env)
			x.formals = lvalCopy(v.formals)
			x.body = lvalCopy(v.body)
		} else {
			x.builtin = v.builtin
		}
	case lvalNumType:
		x.num = v.num
	case lvalErrType:
		x.err = string(v.err)
	case lvalSymType:
		x.sym = string(v.sym)
	case lvalSexprType:
		fallthrough
	case lvalQexprType:
		for _, cell := range v.cells {
			x.cells = append(x.cells, lvalCopy(cell))
		}
	}
	return x
}

func lvalAdd(v *lval, x *lval) *lval {
	v.cells = append(v.cells, x)
	return v
}

func (v *lval) lvalExprString(openChar string, closeChar string) string {
	s := openChar
	for i, cell := range v.cells {
		s += cell.lvalString()
		if i < v.cellCount()-1 {
			s += " "
		}
	}
	s += closeChar
	return s
}

func (v *lval) lvalExprPrint(openChar string, closeChar string) {
	fmt.Print(v.lvalExprString(openChar, closeChar))
}

func lvalReadNum(tree mpc.MpcAst) *lval {
	x, err := strconv.ParseInt(mpc.GetContents(tree), 10, 0)
	if err != nil {
		return lvalErr("Invalid Number: %s", mpc.GetContents(tree))
	}
	return lvalNum(x)
}

func lvalRead(tree mpc.MpcAst) *lval {
	// If Symbol or Number, return conversion to that type
	if strings.Contains(mpc.GetTag(tree), "number") {
		return lvalReadNum(tree)
	}
	if strings.Contains(mpc.GetTag(tree), "symbol") {
		return lvalSym(mpc.GetContents(tree))
	}
	// If root (>) or S-expression, then create empty list
	var x *lval
	if mpc.GetTag(tree) == ">" {
		x = lvalSexpr()
	} else if strings.Contains(mpc.GetTag(tree), "sexpr") {
		x = lvalSexpr()
	} else if strings.Contains(mpc.GetTag(tree), "qexpr") {
		x = lvalQexpr()
	}
	// Fill the cell list with any valid expressions in the children
	for i := 0; i < mpc.GetNumChildren(tree); i++ {
		iChild := mpc.GetChild(tree, i)
		iContents := mpc.GetContents(iChild)
		if iContents == "(" ||
			iContents == ")" ||
			iContents == "{" ||
			iContents == "}" ||
			mpc.GetTag(iChild) == "regex" {
			continue
		} else {
			x = lvalAdd(x, lvalRead(iChild))
		}
		strconv.ParseInt(mpc.GetContents(tree), 10, 0)
	}
	return x
}

func (v *lval) lvalEvalSexpr(e *lenv) *lval {
	// Evaluate children
	for i, cell := range v.cells {
		v.cells[i] = cell.lvalEval(e)
	}
	// Error checking
	for i, cell := range v.cells {
		if cell.ltype == lvalErrType {
			return v.lvalTake(i)
		}
	}
	// Empty Expression
	if v.cellCount() == 0 {
		return v
	}
	// Single Expression
	if v.cellCount() == 1 {
		return v.lvalTake(0).lvalEval(e)
	}
	// Ensure first element is a symbol
	f := v.lvalPop(0)
	if f.ltype != lvalFunType {
		return lvalErr("S-expression does not start with symbol! got: %s", f.ltypeName())
	}
	// Use first element as a function to get result
	return lvalCall(e, f, v)
}

func (v *lval) lvalEval(e *lenv) *lval {
	if v.ltype == lvalSymType {
		return e.lenvGet(v)
	}
	if v.ltype == lvalSexprType {
		return v.lvalEvalSexpr(e)
	}
	return v
}

func (v *lval) lvalPop(i int) *lval {
	x := v.cells[i]
	copy(v.cells[i:], v.cells[i+1:])
	v.cells[len(v.cells)-1] = nil
	v.cells = v.cells[:len(v.cells)-1]
	return x
}

func (v *lval) lvalTake(i int) *lval {
	return v.lvalPop(i)
}

func lvalJoin(x *lval, y *lval) *lval {
	for y.cellCount() > 0 {
		x = lvalAdd(x, y.lvalPop(0))
	}
	return x
}

func lvalCall(e *lenv, f *lval, a *lval) *lval {
	// Simple Builtin case:
	if f.builtin != nil {
		return f.builtin(e, a)
	}
	// Record argument counts
	given := a.cellCount()
	total := f.formals.cellCount()
	// While arguments still remain to be processed
	for a.cellCount() > 0 {
		// If we've ran out of formal arguments to bind
		if f.formals.cellCount() == 0 {
			return lvalErr("Function passed too many arguments. Got %d, Expected %d", given, total)
		}
		// Pop the first symbol from the formal
		sym := f.formals.lvalPop(0)
		// Special case to deal with '&'
		if sym.sym == "&" {
			// Ensure '&' is followed by another symbol
			if f.formals.cellCount() != 1 {
				return lvalErr("Function format invalid. Symbol '&' was not followed by 1 symbol.)")
			}
			// Next formal should be bound to the remaining arguments
			nsym := f.formals.lvalPop(0)
			f.env.lenvPut(nsym, builtinList(e, a))
			break
		}
		// Pop the next argument from the list
		val := a.lvalPop(0)
		// Bind a copy into the function's environment
		f.env.lenvPut(sym, val)
	}
	// If '&' remains in the formal list, bind to an empty list
	if f.formals.cellCount() > 0 && f.formals.cells[0].sym == "&" {
		// Check to ensure that '&' is not passed in invalidly
		if f.formals.cellCount() != 2 {
			return lvalErr("Function forma invalid. Symbol '&' not followed by single symbol")
		}
		// Pop '&' symbol
		f.formals.lvalPop(0)
		// Pop the next symbol and create an empty list
		sym := f.formals.lvalPop(0)
		val := lvalQexpr()
		// Bind to the environment
		f.env.lenvPut(sym, val)
	}
	// If all formals have been bound, evaluate
	if f.formals.cellCount() == 0 {
		// Set environment parent to evaluation environment
		f.env.par = e
		// Evaluate and return
		return builtinEval(f.env, lvalAdd(lvalSexpr(), lvalCopy(f.body)))
	}
	// Otherwise, return partially evaluated function
	return lvalCopy(f)
}
