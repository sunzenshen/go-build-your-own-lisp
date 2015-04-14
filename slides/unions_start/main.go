package main

/*
#include <stdarg.h>
// START UNION OMIT
union foo {
	char   c;
	int    i;
	double d;
};
// END UNION OMIT
*/
import "C"
import "fmt"

func main() {
	// START UNION MAIN OMIT
var f *C.union_foo = new(C.union_foo)
// Insert answers here:

fmt.Println(f)
	// END UNION MAIN OMIT
}
