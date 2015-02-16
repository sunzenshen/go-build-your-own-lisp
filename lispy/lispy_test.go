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
		got, err := l.ReadEval(c.input, false)
		if err != nil {
			t.Errorf("ReadEval could not parse the following input: \"%s\"", c.input)
		} else if got != c.want {
			t.Errorf("ReadEval input: \"%s\" returned: %d, actually expected: %d", c.input, got, c.want)
		}
	}
}
