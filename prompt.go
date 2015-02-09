package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sunzenshen/golispy/parser"
)

func main() {
	// Version and Exit Information
	fmt.Println("Lispy Version 0.0.0.0.2")
	fmt.Print("Press Ctrl+c to Exit\n\n")
	// For reading lines of user input
	scanner := bufio.NewScanner(os.Stdin)

	Number, Operator, Expr, Lispy := golispy.InitLispy()
	defer golispy.CleanLispy(Number, Operator, Expr, Lispy)

	for {
		// Prompt
		fmt.Print("lispy> ")
		// Read a line of user input
		scanner.Scan()
		input := scanner.Text()
		// Echo input back to user
		golispy.ParseInput(input, Lispy)
	}
}
