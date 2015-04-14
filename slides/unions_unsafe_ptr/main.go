package main

/*
union foo {
	char   c;
	int    i;
	double d;
};
*/
import "C"
import "fmt"
import "unsafe"

func main() {
	var f *C.union_foo = new(C.union_foo)
	p := (*C.int)(unsafe.Pointer(f))
	*p = 32767
	fmt.Println(f)
	fmt.Println(*p)
}
