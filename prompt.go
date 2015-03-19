package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sunzenshen/go-build-your-own-lisp/lispy"
)

func main() {
	// Version and Exit Information
	fmt.Println("Lispy Version 0.0.0.0.4")
	fmt.Print("Press Ctrl+c to Exit\n\n")
	// For reading lines of user input
	scanner := bufio.NewScanner(os.Stdin)

	l := lispy.InitLispy()
	defer lispy.CleanLispy(l)

	// Supplied with a list of files
	if len(os.Args) > 1 {
		fmt.Println("Files passed into Lispy interpreter")
		l.LoadFiles(os.Args[1:])
	}

	for {
		// Prompt
		fmt.Print("lispy> ")
		// Read a line of user input
		scanner.Scan()
		input := scanner.Text()
		// Echo input back to user
		l.ReadEvalPrint(input)
	}
}
