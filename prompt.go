package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Version and Exit Information
	fmt.Println("Lispy Version 0.0.0.0.1")
	fmt.Println("Press Ctrl+c to Exit\n")
	// For reading lines of user input
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Prompt
		fmt.Print("lispy> ")
		// Read a line of user input
		scanner.Scan()
		input := scanner.Text()
		// Echo input back to user
		fmt.Printf("No you're a %s\n", input)
	}
}
