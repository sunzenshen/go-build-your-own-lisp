package main

/*
#include <stdarg.h>
// START UNION OMIT
union foo {
	char   c;
	int    i;
	double d;
};
union bar {
	char   c;
	int    i;
};
union baz { char   c; };
// END UNION OMIT
*/
import "C"
import "fmt"
import "unsafe"

func main() {
	// START UNION MAIN OMIT
var f *C.union_foo = new(C.union_foo)
var b *C.union_bar = new(C.union_bar)
var z *C.union_baz = new(C.union_baz)
p := (*C.int)(unsafe.Pointer(f))
*p = 32767
fmt.Println(f)
fmt.Println(b)
fmt.Println(z)
	// END UNION MAIN OMIT
}
