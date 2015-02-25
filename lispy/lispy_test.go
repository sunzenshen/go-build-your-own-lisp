package lispy

import "testing"

func TestLvalRead(t *testing.T) {
	l := InitLispy()
	defer CleanLispy(l)

	cases := []struct {
		input, want string
	}{
		{"+ 2 2", "(+ 2 2)"},
		{"+ 2 (* 7 6) (* 2 5)", "(+ 2 (* 7 6) (* 2 5))"},
		{"*     55     101     (+ 0 0 0)", "(* 55 101 (+ 0 0 0))"},
	}

	for _, c := range cases {
		got := l.Read(c.input, false).lvalString()
		if got != c.want {
			t.Errorf("Read input: \"%s\" returned %s, actually expected %s", c.input, got, c.want)
		}
	}
}

func TestValidIntegerMath(t *testing.T) {
	l := InitLispy()
	defer CleanLispy(l)

	cases := []struct {
		input string
		want  int64
	}{
		{"+ 1 1", 2},
		{"+ 2 -3", -1},
		{"- 3 2", 1},
		{"- 100", -100},
		{"- 0", 0},
		{"* -2 -3", 6},
		{"* 2 3", 6},
		{"* 2 -3", -6},
		{"/ 9 3", 3},
		{"/ -9 3", -3},
		{"/ -9 -3", 3},
		{"/ 7 3", 2},
		{"+ 5 6", 11},
		{"- (* 10 10) (+ 1 1 1)", 97},
		{"+ 1 (* 7 5) 3", 39},
	}

	for _, c := range cases {
		got := l.ReadEval(c.input, false)
		if got.ltype != lvalNumType {
			t.Errorf("ReadEval input: \"%s\" returned ltype %s, err %s actually expected lvalNumType",
				c.input, got.ltypeString(), got.err)
		}
		if got.num != c.want {
			t.Errorf("ReadEval input: \"%s\" returned: %d, actually expected: %d", c.input, got.num, c.want)
		}
	}
}

func TestError(t *testing.T) {
	l := InitLispy()
	defer CleanLispy(l)

	cases := []struct {
		input, want string
	}{
		{"/ 10 0", "Division By Zero!"},
		{"(/ ())", "Cannot operate on non-number!"},
		{"(1 2 3)", "S-expression does not start with symbol!"},
		{"+ - +", "Cannot operate on non-number!"},
		{"The quick brown fox jumps over the very lazy dog.", "Failed to parse input!"},
	}

	for _, c := range cases {
		got := l.ReadEval(c.input, false)
		if got.ltype != lvalErrType {
			t.Errorf("ReadEval input: \"%s\" returned ltype %s, actually expected lvalErrType", c, got.ltypeString())
		}
		if got.err != c.want {
			t.Errorf("ReadEval input: \"%s\" returned err \"%s\", actually expected \"%s\"", c, got.err, c.want)
		}
	}
}
