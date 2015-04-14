package main

/*
#include <stdarg.h>

int sum(int count, ...) {
  va_list ap;
  int i;
  int sum = 0;
  va_start(ap, count);
  for (i = 0; i < count; i++) {
    sum += va_arg(ap, int);
  }
  va_end(ap);
  return sum;
}

// START VARIADIC OMIT
// ...
int sum_if(int a, int b) {
	return sum(2, a, b);
}
*/
import "C"
import "fmt"

func sum(nums ...int) int {
	return int(C.sum_if(C.int(nums[0]), C.int(nums[1])))
}

func main() {
	// int sum(int count, ...)
	fmt.Println("Sum: ", sum(1, 2))
}

// END VARIADIC OMIT
