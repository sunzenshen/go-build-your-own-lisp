package lispy

import "fmt"

type lenv struct {
	syms []string
	vals []*lval
}

func lenvNew() *lenv {
	e := new(lenv)
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

func (e *lenv) lenvGet(k *lval) *lval {
	for i, sym := range e.syms {
		if sym == k.sym {
			return lvalCopy(e.vals[i])
		}
	}
	return lvalErr("Unbound Symbol!")
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

func (e *lenv) lenvAddBuiltin(name string, function lbuiltin) {
	k := lvalSym(name)
	v := lvalFun(function)
	e.lenvPut(k, v)
}

func (e *lenv) lenvAddBuiltins() {
	// List Functions
	e.lenvAddBuiltin("def", builtinDef)
	e.lenvAddBuiltin("list", builtinList)
	e.lenvAddBuiltin("head", builtinHead)
	e.lenvAddBuiltin("tail", builtinTail)
	e.lenvAddBuiltin("eval", builtinEval)
	e.lenvAddBuiltin("join", builtinJoin)
	// Mathematical Functions
	e.lenvAddBuiltin("+", builtinAdd)
	e.lenvAddBuiltin("-", builtinSub)
	e.lenvAddBuiltin("*", builtinMul)
	e.lenvAddBuiltin("/", builtinDiv)
}
