package main

/*
#cgo LDFLAGS: -ledit
#include <stdlib.h>
#include <editline/readline.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {
	// Version and Exit Information
	fmt.Println("Lispy Version 0.0.0.0.2")
	fmt.Println("Pressing Ctrl+c is broken in this version\n")
	// Prompt for user input
	prompt := C.CString("lispy> ")

	for {
		// Output prompt and get user input
		raw_input := C.readline(prompt)
		input := C.GoString(raw_input)
		// Add input to history
		C.add_history(raw_input)
		// Echo input back to user
		fmt.Printf("No you're a %s\n", input)
		// Free retrieved input
		C.free(unsafe.Pointer(raw_input))
	}
	C.free(unsafe.Pointer(prompt))
}
