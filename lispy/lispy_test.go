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

func TestListFunctions(t *testing.T) {
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

func TestVariableArgsAndCurrying(t *testing.T) {
	l := InitLispy()
	defer CleanLispy(l)

	cases := []struct {
		input string
		want  string
	}{
		{"\\ {args body} {def (head args) (\\ (tail args) body)}",
			"(\\ {args body} {def (head args) (\\ (tail args) body)})"},
		{"fun {add-together x y} {+ x y}", "()"},
		{"add-together 1 2", "3"},
		{"curry + {5 6 7}", "18"},
		{"uncurry head 5 6 7", "{5}"},
		{"def {add-uncurried} +", "()"},
		{"def {add-curried} (curry +)", "()"},
		{"add-curried {5 6 7}", "18"},
		{"add-uncurried 5 6 7", "18"},
	}

	for _, c := range cases {
		got := l.ReadEval(c.input, false)
		if got.lvalString() != c.want {
			t.Errorf("ReadEval input: \"%s\" returned: \"%s\", actually expected: \"%s\"", c.input, got.lvalString(), c.want)
		}
	}
}

func TestConditionals(t *testing.T) {
	l := InitLispy()
	defer CleanLispy(l)

	truth := "1"
	falsity := "0"

	cases := []struct {
		input string
		want  string
	}{
		{"> 10 5", truth},
		{"<= 88 5", falsity},
		{"== 5 6", falsity},
		{"== 5 {}", falsity},
		{"== 1 1", truth},
		{"!= {} 56", truth},
		{"== {1 2 3 { 5 6}} {1  2 3  {5 6}}", truth},
		{"def {x y} 100 200", "()"},
		{"if (== x y) {+ x y} {- x y}", "-100"},
		// Standard Library
		{"== nil {}", truth},
		{"== true 1", truth},
		{"== false 0", truth},
		{"!= true false", truth},
	}

	for _, c := range cases {
		got := l.ReadEval(c.input, false)
		if got.lvalString() != c.want {
			t.Errorf("ReadEval input: \"%s\" returned: \"%s\", actually expected: \"%s\"", c.input, got.lvalString(), c.want)
		}
	}
}

func TestRecursiveFunctions(t *testing.T) {
	l := InitLispy()
	defer CleanLispy(l)

	cases := []struct {
		input string
		want  string
	}{
		// Define recursive len function
		{`(fun {len l} {
			if (== l {})
				{0}
				{+ 1 (len (tail l))}
		})`, "()"},
		// Define recursive reverse function
		{`(fun {reverse l} {
			if (== l {})
				{{}}
				{join (reverse (tail l)) (head l)}
		})`, "()"},
		// Test out the recursive functions
		{"len {}", "0"},
		{"len {1 2 3}", "3"},
		{"reverse {}", "{}"},
		{"reverse {1 2 3}", "{3 2 1}"},
	}

	for _, c := range cases {
		got := l.ReadEval(c.input, false)
		if got.lvalString() != c.want {
			t.Errorf("ReadEval input: \"%s\" returned: \"%s\", actually expected: \"%s\"", c.input, got.lvalString(), c.want)
		}
	}
}

func TestStrings(t *testing.T) {
	l := InitLispy()
	defer CleanLispy(l)

	cases := []struct {
		input string
		want  string
	}{
		// Strings
		{"\"hello\"", "\"hello\""},
		{"\"hello\\n\"", "\"hello\\n\""},
		{"\"hello\\\"\"", "\"hello\\\"\""},
		{"head {\"hello\" \"world\"}", "{\"hello\"}"},
		{"eval (head {\"hello\" \"world\"})", "\"hello\""},
		// Comments
		{"; Ignore this comment", "()"},
		{"; + 1 2", "()"},
	}

	for _, c := range cases {
		got := l.ReadEval(c.input, false)
		if got.lvalString() != c.want {
			t.Errorf("ReadEval input: \"%s\" returned: \"%s\", actually expected: \"%s\"", c.input, got.lvalString(), c.want)
		}
	}
}

func TestStandardLibrary(t *testing.T) {
	l := InitLispy()
	defer CleanLispy(l)

	cases := []struct {
		input string
		want  string
	}{
		// fst snd trd
		{"fst {\"first\"}", "\"first\""},
		{"fst {nil}", "{}"},
		{"fst {1}", "1"},
		{"fst {{1}}", "{1}"},
		{"snd {1 2}", "2"},
		{"trd {1 2 3}", "3"},
		{"fst {-1 2 -3 4}", "-1"},
		{"snd {1 -2 -3 -4}", "-2"},
		{"trd {-1 2 -3 4}", "-3"},
		{"trd {-1 2 \"third\" 4}", "\"third\""},
		// length of list
		{"len nil", "0"},
		{"len {}", "0"},
		{"len {42}", "1"},
		{"len {1 2}", "2"},
		{"len {-1 -2 -3}", "3"},
		{"len {10 20 30 40}", "4"},
		{"len {10 true \"third\" nil}", "4"},
		// nth element of list
		{"nth 0 {1}", "1"},
		{"nth 1 {1 2}", "2"},
		{"nth 2 {1 2 3}", "3"},
		{"nth 0 {-0 1 -2}", "0"},
		{"nth 1 {0 -1 2}", "-1"},
		{"nth 2 {-0 1 -2}", "-2"},
		{"nth 2 {-0 1 nil}", "{}"},
		{"nth 2 {-0 1 \"nth\"}", "\"nth\""},
		// last element of list
		{"last {1}", "1"},
		{"last {-1 -2}", "-2"},
		{"last {1 2 \"three\"}", "\"three\""},
		// take N elements of list
		{"take 0 nil", "{}"},
		{"take 0 {nil}", "{}"},
		{"take 1 {nil}", "{nil}"},
		{"take 2 {nil \"two\"}", "{nil \"two\"}"},
		{"take 3 {1 2 3}", "{1 2 3}"},
		{"take 2 {1 2 3}", "{1 2}"},
		{"take 1 {1 2 3}", "{1}"},
		{"take 0 {1 2 3}", "{}"},
	}

	for _, c := range cases {
		got := l.ReadEval(c.input, false)
		if got.lvalString() != c.want {
			t.Errorf("ReadEval input: \"%s\" returned: \"%s\", actually expected: \"%s\"", c.input, got.lvalString(), c.want)
		}
	}
}
