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
		{"{1 2 3 4}", "({1 2 3 4})"},
		{"{1 2 (+ 5 6) 4}", "({1 2 (+ 5 6) 4})"},
		{"{{2 3 4} {1}}", "({{2 3 4} {1}})"},
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
				c.input, got.ltypeName(), got.err)
		}
		if got.num != c.want {
			t.Errorf("ReadEval input: \"%s\" returned: %d, actually expected: %d", c.input, got.num, c.want)
		}
	}
}

func TestStringOutput(t *testing.T) {
	l := InitLispy()
	defer CleanLispy(l)

	cases := []struct {
		input string
		want  string
	}{
		{"list 1 2 3 4", "{1 2 3 4}"},
		{"{head (list 1 2 3 4)}", "{head (list 1 2 3 4)}"},
		{"eval {head (list 1 2 3 4)}", "{1}"},
		{"tail {tail tail tail}", "{tail tail}"},
		{"eval (tail {tail tail {5 6 7}})", "{6 7}"},
		{"eval (head {(+ 1 2) (+ 10 20)})", "3"},
		{"eval (head {5 10 11 15})", "5"},
		{"+", "<builtin>"},
		{"eval (head {+ - = - * /})", "<builtin>"},
		{"(eval (head {+ - = - * /})) 10 20", "30"},
		{"hello", "Error: Unbound Symbol: 'hello'"},
		{"+ 1 {5 6 7}", "Error: Cannot operate on non-number: Q-Expression"},
		{"head {1 2 3} {4 5 6}",
			"Error: Function 'head' passed too many arguments: ({1 2 3} {4 5 6})"},
	}

	for _, c := range cases {
		got := l.ReadEval(c.input, false)
		if got.lvalString() != c.want {
			t.Errorf("ReadEval input: \"%s\" returned: \"%s\", actually expected: \"%s\"", c.input, got.lvalString(), c.want)
		}
	}
}

func TestVariableAssignments(t *testing.T) {
	l := InitLispy()
	defer CleanLispy(l)

	cases := []struct {
		input string
		want  string
	}{
		{"def {x} 100", "()"},
		{"def {y} 200", "()"},
		{"x", "100"},
		{"y", "200"},
		{"+ x y", "300"},
		{"def {a b} 5 6", "()"},
		{"+ a b", "11"},
		{"def {arglist} {a b x y}", "()"},
		{"arglist", "{a b x y}"},
		{"def arglist 1 2 3 4", "()"},
		{"list a b x y", "{1 2 3 4}"},
	}

	for _, c := range cases {
		got := l.ReadEval(c.input, false)
		if got.lvalString() != c.want {
			t.Errorf("ReadEval input: \"%s\" returned: \"%s\", actually expected: \"%s\"", c.input, got.lvalString(), c.want)
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
		{"(/ ())", "Cannot operate on non-number: S-Expression"},
		{"(1 2 3)", "S-expression does not start with symbol! got: Number"},
		{"+ - +", "Cannot operate on non-number: Function"},
		{"The quick brown fox jumps over the very lazy dog.",
			"Failed to parse input: 'The quick brown fox jumps over the very lazy dog.'"},
	}

	for _, c := range cases {
		got := l.ReadEval(c.input, false)
		if got.ltype != lvalErrType {
			t.Errorf("ReadEval input: \"%s\" returned ltype %s, actually expected lvalErrType", c, got.ltypeName())
		}
		if got.err != c.want {
			t.Errorf("ReadEval input: \"%s\" returned err \"%s\", actually expected \"%s\"", c, got.err, c.want)
		}
	}
}

func TestFunctionDefinitions(t *testing.T) {
	l := InitLispy()
	defer CleanLispy(l)

	cases := []struct {
		input string
		want  string
	}{
		{"(\\ {x y} {+ x y})", "(\\ {x y} {+ x y})"},
		{"(\\ {x y} {+ x y}) 10 20", "30"},
		{"def {add-together} (\\ {x y} {+ x y})", "()"},
		{"add-together", "(\\ {x y} {+ x y})"},
		{"add-together 10 20", "30"},
		{"add-together", "(\\ {x y} {+ x y})"}, // Check for accidental modification
		{"def {add-mul} (\\ {x y} {+ x (* x y)})", "()"},
		{"add-mul", "(\\ {x y} {+ x (* x y)})"},
		{"add-mul 10 20", "210"},
		{"add-mul 10", "(\\ {y} {+ x (* x y)})"},
		{"def {add-mul-ten} (add-mul 10)", "()"},
		{"add-mul-ten", "(\\ {y} {+ x (* x y)})"},
		{"add-mul 10 50", "510"},
		{"add-mul-ten 50", "510"},
		{"add-mul", "(\\ {x y} {+ x (* x y)})"},   // Check for accidental modification
		{"add-mul-ten", "(\\ {y} {+ x (* x y)})"}, // Check for accidental modification
	}

	for _, c := range cases {
		got := l.ReadEval(c.input, false)
		if got.lvalString() != c.want {
			t.Errorf("ReadEval input: \"%s\" returned: \"%s\", actually expected: \"%s\"", c.input, got.lvalString(), c.want)
		}
	}
}
