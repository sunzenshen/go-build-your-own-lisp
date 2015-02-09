package golispy

// #cgo LDFLAGS: -ledit -lm
// #include "mpc_interface.h"
import "C"
import "unsafe"

func CleanLispy(number *C.struct_mpc_parser_t,
	operator *C.struct_mpc_parser_t,
	expr *C.struct_mpc_parser_t,
	lispy *C.struct_mpc_parser_t) {
	defer C.mpc_cleanup_if(C.int(4), number, operator, expr, lispy)
}

func InitLispy() (
	*C.struct_mpc_parser_t,
	*C.struct_mpc_parser_t,
	*C.struct_mpc_parser_t,
	*C.struct_mpc_parser_t) {
	Number := mpcNew("number")
	Operator := mpcNew("operator")
	Expr := mpcNew("expr")
	Lispy := mpcNew("lispy")
	language := "" +
		"number : /-?[0-9]+/                               ; " +
		"operator : '+' | '-' | '*' | '/'                  ; " +
		"expr     : <number> | '(' <operator> <expr>+ ')'  ; " +
		"lispy    : /^/ <operator> <expr>+ /$/             ; "
	cLanguage := C.CString(language)
	defer C.free(unsafe.Pointer(cLanguage))
	C.mpca_lang_if(C.MPCA_LANG_DEFAULT,
		cLanguage,
		Number, Operator, Expr, Lispy)
	return Number, Operator, Expr, Lispy
}

func mpcNew(name string) *C.struct_mpc_parser_t {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return C.mpc_new(cName)
}

func ParseInput(input string, mpcParser *C.struct_mpc_parser_t) {
	var r C.mpc_result_t
	cInput := C.CString(input)
	defer C.free(unsafe.Pointer(cInput))
	stdin := C.CString("<stdin>")
	defer C.free(unsafe.Pointer(stdin))
	if C.mpc_parse(stdin, cInput, mpcParser, &r) != C.int(0) {
		C.mpc_ast_print(C.get_output(&r))
		C.mpc_ast_delete(C.get_output(&r))
	} else {
		C.mpc_err_print(C.get_error(&r))
		C.mpc_err_delete(C.get_error(&r))
	}
}
