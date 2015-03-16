package mpc

// #include "mpc_interface.h"
import "C"
import (
	"errors"
	"unsafe"
)

// MpcAst is a pointer to a mpc-generated AST
type MpcAst *C.mpc_ast_t

// MpcParser is a pointer to a parser created with mpc
type MpcParser *C.struct_mpc_parser_t

// MpcResult is a union that returns either an output or error
type MpcResult C.mpc_result_t

// GetNumChildren accesses the children_num of MpcAst
func GetNumChildren(node MpcAst) int {
	return int(C.get_children_num(node))
}

// GetChild accesses the child at a specific index in MpcAst
func GetChild(node MpcAst, index int) MpcAst {
	return C.get_child(node, C.int(index))
}

// GetContents accesses the contents of an MpcAst
func GetContents(node MpcAst) string {
	return C.GoString(node.contents)
}

// GetOperator accesses a MpcAst's child node representing an operator
func GetOperator(node MpcAst) string {
	return GetContents(GetChild(node, 1))
}

// GetOutput accesses the output field of an input MpcResult
func GetOutput(result *C.mpc_result_t) MpcAst {
	return C.get_output(result)
}

// GetTag accesses the tag field of an MpcAst
func GetTag(node MpcAst) string {
	return C.GoString(node.tag)
}

// MpcAstDelete cleans up an AST contained in a MpcResult
func MpcAstDelete(result *C.mpc_result_t) {
	C.mpc_ast_delete(C.get_output(result))
}

// MpcaLang uses a language definition to generate parsers
func MpcaLang(language string, parsers ...MpcParser) {
	cLanguage := C.CString(language)
	defer C.free(unsafe.Pointer(cLanguage))
	C.mpca_lang_if(C.MPCA_LANG_DEFAULT, cLanguage,
		parsers[0], parsers[1], parsers[2], parsers[3], parsers[4], parsers[5])
}

// MpcCleanup calls mpc's cleanup function indirectly, using a wrapper for the variadic args
func MpcCleanup(parsers ...MpcParser) {
	C.mpc_cleanup_if(C.int(len(parsers)),
		parsers[0], parsers[1], parsers[2], parsers[3], parsers[4], parsers[5])
}

// MpcNew returns a pointer to an initiated mpc parser
func MpcNew(name string) MpcParser {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return C.mpc_new(cName)
}

// MpcParse takes an input string and generates an MpcResult
func MpcParse(input string, parser MpcParser) (C.mpc_result_t, error) {
	var r C.mpc_result_t
	cInput := C.CString(input)
	defer C.free(unsafe.Pointer(cInput))
	stdin := C.CString("<stdin>")
	defer C.free(unsafe.Pointer(stdin))
	var err error
	if C.mpc_parse(stdin, cInput, parser, &r) == C.int(0) {
		err = errors.New("mpc: failed to parse input string")
	}
	return r, err
}

// MpcErrDelete cleans up the error info in an MpcResult
func MpcErrDelete(result *C.mpc_result_t) {
	C.mpc_err_delete(C.get_error(result))
}

// MpcErrPrint prints the reason for failed AST parsing
func MpcErrPrint(result *C.mpc_result_t) {
	C.mpc_err_print(C.get_error(result))
}

// PrintAst parses an input string for its AST representation
func PrintAst(input string, mpcParser MpcParser) {
	r, err := MpcParse(input, mpcParser)
	if err != nil {
		MpcErrPrint(&r)
		MpcErrDelete(&r)
	} else {
		C.mpc_ast_print(C.get_output(&r))
		MpcAstDelete(&r)
	}
}
