package mpcinterface

// #cgo LDFLAGS: -ledit -lm
// #include "mpc_interface.h"
import "C"
import "unsafe"

// MpcAst is a pointer to a mpc-generated AST
type MpcAst *C.mpc_ast_t

func getChild(node MpcAst, index int) MpcAst {
	return C.get_child(node, C.int(index))
}

func getContents(node MpcAst) string {
	return C.GoString(node.contents)
}

func getOperator(node MpcAst) string {
	return getContents(getChild(node, 1))
}

func getTag(node MpcAst) string {
	return C.GoString(node.tag)
}

func mpcNew(name string) *C.struct_mpc_parser_t {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return C.mpc_new(cName)
}

// ParseInput prints the AST of an input string and parser
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
