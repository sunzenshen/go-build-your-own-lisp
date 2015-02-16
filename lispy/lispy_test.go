package lispy

import "testing"

func TestValidIntegerMath(t *testing.T) {
	l := InitLispy()
	defer CleanLispy(l)

	cases := []struct {
		input string
		want  int64
	}{
		{"+ 1 1", 2},
		{"+ 2 -3", -1},
		{"add 2 -3", -1},
		{"* -2 -3", 6},
		{"* 2 3", 6},
		{"* 2 -3", -6},
		{"mul 2 -3", -6},
		{"/ 9 3", 3},
		{"/ -9 3", -3},
		{"/ -9 -3", 3},
		{"/ 7 3", 2},
		{"div -9 -3", 3},
		{"% 7 3", 1},
		{"% 6 3", 0},
		{"% -7 3", -1},
		{"% -7 -3", -1},
		{"mod -7 -3", -1},
		{"^ -2 0", 1},
		{"^ -2 1", -2},
		{"^ -2 2", 4},
		{"^ 9 -1", 1},
		{"pow -2 2", 4},
		{"+ 5 6", 11},
		{"- (* 10 10) (+ 1 1 1)", 97},
	}
	for _, c := range cases {
		got := l.ReadEval(c.input, false)
		if got.ltype != lvalNumType {
			t.Errorf("ReadEval input: \"%s\" returned ltype %s, err %s actually expected lvalNumType",
				c.input, ltypeString(got), lerrString(got))
		}
		if got.num != c.want {
			t.Errorf("ReadEval input: \"%s\" returned: %d, actually expected: %d", c.input, got, c.want)
		}
	}
}

func TestDivisionByZero(t *testing.T) {
	l := InitLispy()
	defer CleanLispy(l)

	cases := []string{
		"/ 42 0",
		"div 42 0",
	}

	for _, c := range cases {
		got := l.ReadEval(c, false)
		if got.ltype != lvalErrType {
			t.Errorf("ReadEval input: \"%s\" returned ltype %s, actually expected lvalErrType", c, ltypeString(got))
		}
		if got.err != lerrDivZero {
			t.Errorf("ReadEval input: \"%s\" returned err %s, actually expected lerrDivZero", c, lerrString(got))
		}
	}
}

func TestFailedParse(t *testing.T) {
	l := InitLispy()
	defer CleanLispy(l)
	c := "The quick brown fox jumps over the very lazy dog."
	got := l.ReadEval(c, false)
	if got.ltype != lvalErrType || got.err != lerrParseFail {
		t.Errorf("ReadEval input: \"%s\" returned ltype %s, num %d, err %s, actually expected lerrParseFail",
			c, ltypeString(got), got.num, lerrString(got))
	}
}
