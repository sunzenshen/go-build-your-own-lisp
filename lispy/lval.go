package lispy

import "fmt"

// ltype values for lval
const (
	lvalNumType = iota
	lvalErrType
)

// err values for lval
const (
	lerrDivZero = iota
	lerrBadOp
	lerrBadNum
	lerrParseFail
)

type lval struct {
	ltype int8
	num   int64
	err   int8
}

// lvalNum creates an lval number
func lvalNum(x int64) lval {
	var v lval
	v.ltype = lvalNumType
	v.num = x
	return v
}

// lvalErr creates an lval error
func lvalErr(x int8) lval {
	var v lval
	v.ltype = lvalErrType
	v.err = x
	return v
}

func lvalPrint(v lval) {
	switch v.ltype {
	case lvalNumType:
		fmt.Print(v.num)
	case lvalErrType:
		if v.err == lerrDivZero {
			fmt.Print("Error: Division By Zero!")
		} else if v.err == lerrBadOp {
			fmt.Print("Error: Invalid Operator!")
		} else if v.err == lerrBadNum {
			fmt.Print("Error: Invalid Number!")
		} else if v.err == lerrParseFail {
			fmt.Print("Error: Failed to parse input!")
		}
	}
}

func lvalPrintLn(v lval) {
	lvalPrint(v)
	fmt.Print("\n")
}
