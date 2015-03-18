package lispy

import "fmt"

type lenv struct {
	par  *lenv
	syms []string
	vals []*lval
}

func lenvNew() *lenv {
	e := new(lenv)
	e.par = nil
	e.syms = nil
	e.vals = nil
	return e
}

func (e *lenv) count() int {
	if len(e.syms) != len(e.vals) {
		fmt.Printf("Counts for lenv don't match! syms: %d vals:%d", len(e.syms), len(e.vals))
		return -1
	}
	return len(e.syms)
}

func lenvCopy(e *lenv) *lenv {
	n := new(lenv)
	n.par = e.par
	for i := 0; i < e.count(); i++ {
		n.syms = append(n.syms, string(e.syms[i]))
		n.vals = append(n.vals, lvalCopy(e.vals[i]))
	}
	return n
}

func (e *lenv) lenvGet(k *lval) *lval {
	for i := 0; i < e.count(); i++ {
		if e.syms[i] == k.sym {
			return lvalCopy(e.vals[i])
		}
	}
	// If no symbol is found yet, check in the parent
	if e.par != nil {
		return e.par.lenvGet(k)
	}
	return lvalErr("Unbound Symbol: '%s'", k.sym)
}

func (e *lenv) lenvPut(k, v *lval) {
	// If existing entry is found, overwrite it
	for i, sym := range e.syms {
		if sym == k.sym {
			e.vals[i] = lvalCopy(v)
			return
		}
	}
	// If no existing entry is found, add new entry
	e.vals = append(e.vals, v)
	e.syms = append(e.syms, k.sym)
}

func (e *lenv) lenvDef(k *lval, v *lval) {
	// Find top parent
	for e.par != nil {
		e = e.par
	}
	// Put value in e
	e.lenvPut(k, v)
}

func (e *lenv) lenvAddBuiltin(name string, function lbuiltin) {
	k := lvalSym(name)
	v := lvalFun(function)
	e.lenvPut(k, v)
}

func (e *lenv) lenvAddBuiltins() {
	// Define Functions
	e.lenvAddBuiltin("def", builtinDef)
	e.lenvAddBuiltin("=", builtinPut)
	e.lenvAddBuiltin("\\", builtinLambda)
	// List Functions
	e.lenvAddBuiltin("list", builtinList)
	e.lenvAddBuiltin("head", builtinHead)
	e.lenvAddBuiltin("tail", builtinTail)
	e.lenvAddBuiltin("eval", builtinEval)
	e.lenvAddBuiltin("join", builtinJoin)
	// Comparison Functions
	e.lenvAddBuiltin("if", builtinIf)
	e.lenvAddBuiltin("==", builtinEqual)
	e.lenvAddBuiltin("!=", builtinNotEqual)
	e.lenvAddBuiltin(">", builtinGreaterThan)
	e.lenvAddBuiltin("<", builtinLessThan)
	e.lenvAddBuiltin(">=", builtinGreaterEqual)
	e.lenvAddBuiltin("<=", builtinLessEqual)
	// Mathematical Functions
	e.lenvAddBuiltin("+", builtinAdd)
	e.lenvAddBuiltin("-", builtinSub)
	e.lenvAddBuiltin("*", builtinMul)
	e.lenvAddBuiltin("/", builtinDiv)
}
