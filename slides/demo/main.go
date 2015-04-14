package main

import "github.com/sunzenshen/go-build-your-own-lisp/lispy"

func main() {
	l := lispy.InitLispy()
	defer lispy.CleanLispy(l)

	// workaround for present tool not being able to retrieve standard library
	l.ReadEval("def {nil} {}", false)
	l.ReadEval("def {true} 1", false)
	l.ReadEval("def {false} 0", false)
	l.ReadEval("(def {fun} (\\ {f b} {def (head f) (\\ (tail f) b)}))", false)
	l.ReadEval(`(fun {fst l} { eval (head l) })`, false)
	l.ReadEval(`(fun {snd l} { eval (head (tail l)) })`, false)
	l.ReadEval(`(fun {unpack f xs} {eval (join (list f) xs)})`, false)
	l.ReadEval(`(def {otherwise} true)`, false)

	l.ReadEval(`
; Select statement
(fun {select & cs} {
  if (== cs nil)
    {error "No selection found!"}
    {if (fst (fst cs))
      {snd (fst cs)}
      {unpack select (tail cs)}}
})
	`, false)

	input := `
	// START DEMO OMIT
; Let’s start off with some examples:
(print "Hello, Gophers")
(if (== (+ 2 2) 5)
  {print "five year plan in four years"}
  {print "Arithmetic Nonsense"}
)

; Fibonacci
(fun {fib n} {
  select
    {(== n 0) 0}
    {(== n 1) 1}
    {otherwise (+ (fib (- n 1))
                  (fib (- n 2)))}
})

(def {result} (fib 8))
(print result)
	// END DEMO OMIT
	`
	l.ReadEval(input, false)
}
