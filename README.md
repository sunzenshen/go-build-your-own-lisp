# (Go) Build Your Own Lisp
[![Build Status](https://api.travis-ci.org/sunzenshen/go-build-your-own-lisp.png?branch=master)](https://travis-ci.org/sunzenshen/go-build-your-own-lisp)

This repository tracks my progress working through the book [Build Your Own Lisp](http://buildyourownlisp.com) by [Daniel Holden](https://github.com/orangeduck) using the [Go](http://golang.org) programming language.

#### ... wait, wasn’t that book about coding a Lisp in C?

After signing up for both [LambaConf](http://www.degoesconsulting.com/lambdaconf-2015/) and [GopherCon](http://gophercon.com/) for 2015, I thought it would be fun to implement a toy functional programming language in Go. Translating the C examples into equivalent Go was also an hands-on way to experiment with my understanding of both languages.

# Dependencies

This repository makes extensive use of the [mpc](https://github.com/orangeduck/mpc) library required by the book. [Cgo](http://golang.org/cmd/cgo/) is used to integrate this C code with the rest of the Go-based project.

[Here](http://sunzenshen.github.io/tutorials/2015/05/09/cgotchas-intro.html) is a post that explains some of the design decisions regarding the Cgo usage in this project.

An isolated branch also includes some experimentation with integrating [editline](https://github.com/troglobit/editline) into the command prompt code.
