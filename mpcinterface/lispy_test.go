package mpcinterface

import "testing"

func TestValidIntegerMath(t *testing.T) {
	lispy := InitLispy()
	defer CleanLispy(lispy)

	cases := []struct {
		input string
		want  int64
	}{
		{"+ 1 1", 2},
		{"+ 2 -3", -1},
		{"* -2 -3", 6},
		{"* 2 3", 6},
		{"* 2 -3", -6},
		{"/ 9 3", 3},
		{"/ -9 3", -3},
		{"/ -9 -3", 3},
		{"/ 7 3", 2},
		{"+ 5 6", 11},
		{"- (* 10 10) (+ 1 1 1)", 97},
	}
	for _, c := range cases {
		got, err := ReadEval(c.input, lispy.Lispy)
		if err != nil {
			t.Errorf("ReadEval could not parse the following input: \"%s\"", c.input)
		} else if got != c.want {
			t.Errorf("ReadEval input: \"%s\" returned: %d, actually expected: %d", c.input, got, c.want)
		}
	}
}
