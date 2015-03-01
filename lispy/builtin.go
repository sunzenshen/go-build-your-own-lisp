package lispy

type lbuiltin func(*lenv, *lval) *lval

func builtinOp(e *lenv, a *lval, op string) *lval {
	// Ensure all arguments are numbers
	for _, cell := range a.cells {
		if cell.ltype != lvalNumType {
			return lvalErr("Cannot operate on non-number!")
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

func builtinHead(e *lenv, a *lval) *lval {
	// Check for error conditions
	if a.cellCount() != 1 {
		return lvalErr("Function 'head' passed too many arguments!")
	}
	if a.cells[0].ltype != lvalQexprType {
		return lvalErr("Function 'head' passed incorrect types!")
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
		return lvalErr("Function 'tail' passed too many arguments!")
	}
	if a.cells[0].ltype != lvalQexprType {
		return lvalErr("Function 'tail' passed incorrect types!")
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
		return lvalErr("Function 'eval' passed too many arguments!")
	}
	if a.cells[0].ltype != lvalQexprType {
		return lvalErr("Function 'eval' passed incorrect type!")
	}
	x := a.lvalTake(0)
	x.ltype = lvalSexprType
	return x.lvalEval(e)
}

func builtinJoin(e *lenv, a *lval) *lval {
	for _, cell := range a.cells {
		if cell.ltype != lvalQexprType {
			return lvalErr("Function 'join' passed incorrect type.")
		}
	}
	x := a.lvalPop(0)
	for a.cellCount() > 0 {
		x = lvalJoin(x, a.lvalPop(0))
	}
	return x
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
