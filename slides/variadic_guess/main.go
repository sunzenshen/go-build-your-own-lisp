package main

/*
#include <stdarg.h>

int sum(int count, ...) {
  va_list ap; int i; int sum = 0;
  va_start(ap, count);
  for (i = 0; i < count; i++) { sum += va_arg(ap, int); }
  va_end(ap);
  return sum;
}
*/
import "C"
import "fmt"

func main() {
	// START VARIADIC GUESS OMIT
	// Insert your answer here:
	sum := 0
	fmt.Println("Sum: ", sum)
	// END VARIADIC GUESS OMIT
}
