package lispy

import "github.com/sunzenshen/go-build-your-own-lisp/mpc"

type lenv struct {
	parser mpc.ParserPtr
	par    *lenv
	syms   map[string]*lval
}

func lenvNew() *lenv {
	e := new(lenv)
	e.parser = nil
	e.par = nil
	e.syms = make(map[string]*lval)
	return e
}

func (e *lenv) count() int {
	return len(e.syms)
}

func lenvCopy(e *lenv) *lenv {
	n := new(lenv)
	e.parser = nil
	n.par = e.par
	n.syms = make(map[string]*lval)
	for k, v := range e.syms {
		n.syms[k] = v
	}
	return n
}

func (e *lenv) lenvGet(k *lval) *lval {
	previous := e.syms[k.sym]
	if previous != nil {
		return lvalCopy(previous)
	}

	// If no symbol is found yet, check in the parent
	if e.par != nil {
		return e.par.lenvGet(k)
	}
	return lvalErr("Unbound Symbol: '%s'", k.sym)
}

func (e *lenv) lenvPut(k, v *lval) {
	// If existing entry is found, overwrite it
	previous := e.syms[k.sym]
	if previous != nil {
		e.syms[k.sym] = lvalCopy(v)
	} else {
		// If no existing entry is found, add new entry
		e.syms[k.sym] = v
	}
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
	// String Functions
	e.lenvAddBuiltin("load", builtinLoad)
	e.lenvAddBuiltin("error", builtinError)
	e.lenvAddBuiltin("print", builtinPrint)
	// Mathematical Functions
	e.lenvAddBuiltin("+", builtinAdd)
	e.lenvAddBuiltin("-", builtinSub)
	e.lenvAddBuiltin("*", builtinMul)
	e.lenvAddBuiltin("/", builtinDiv)
	e.lenvAddBuiltin("%", builtinMod)
	e.lenvAddBuiltin("^", builtinPow)
}
