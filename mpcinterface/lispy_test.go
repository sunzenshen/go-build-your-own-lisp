package mpcinterface

import "testing"

func TestValidIntegerMath(t *testing.T) {
	Number, Operator, Expr, Lispy := InitLispy()
	defer CleanLispy(Number, Operator, Expr, Lispy)

	cases := []struct {
		input string
		want  int64
	}{
		{"+ 1 1", 2},
	}
	for _, c := range cases {
		got, err := ReadEval(c.input, Lispy)
		if err != nil {
			t.Errorf("ReadEval could not parse the following input: \"%s\"", c.input)
		} else if got != c.want {
			t.Errorf("ReadEval input: \"%s\" returned: %d, actually expected: %d", c.input, got, c.want)
		}
	}
}
