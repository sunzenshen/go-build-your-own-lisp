package mpcinterface

// #cgo LDFLAGS: -ledit -lm
// #include "mpc_interface.h"
import "C"
import "unsafe"

func getChild(node *C.mpc_ast_t, index int) *C.mpc_ast_t {
	return C.get_child(node, C.int(index))
}

func getContents(node *C.mpc_ast_t) string {
	return C.GoString(node.contents)
}

func getOperator(node *C.mpc_ast_t) string {
	return getContents(getChild(node, 1))
}

func getTag(node *C.mpc_ast_t) string {
	return C.GoString(node.tag)
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
