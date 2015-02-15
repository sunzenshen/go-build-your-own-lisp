package mpcinterface

// #cgo LDFLAGS: -ledit -lm
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

func getChild(node MpcAst, index int) MpcAst {
	return C.get_child(node, C.int(index))
}

func getContents(node MpcAst) string {
	return C.GoString(node.contents)
}

func getOperator(node MpcAst) string {
	return getContents(getChild(node, 1))
}

// GetOutput accesses the output field of an input MpcResult
func GetOutput(result *C.mpc_result_t) MpcAst {
	return C.get_output(result)
}

func getTag(node MpcAst) string {
	return C.GoString(node.tag)
}

// MpcAstDelete cleans up an AST contained in a MpcResult
func MpcAstDelete(result *C.mpc_result_t) {
	C.mpc_ast_delete(C.get_output(result))
}

// MpcaLang uses a language definition to generate parsers
func MpcaLang(language string,
	parser1 MpcParser,
	parser2 MpcParser,
	parser3 MpcParser,
	parser4 MpcParser) {
	cLanguage := C.CString(language)
	defer C.free(unsafe.Pointer(cLanguage))
	C.mpca_lang_if(C.MPCA_LANG_DEFAULT, cLanguage, parser1, parser2, parser3, parser4)
}

// MpcCleanup calls mpc's cleanup function indirectly, using a wrapper for the variadic args
func MpcCleanup(parser1 MpcParser, parser2 MpcParser, parser3 MpcParser, parser4 MpcParser) {
	C.mpc_cleanup_if(C.int(4), parser1, parser2, parser3, parser4)
}

func mpcNew(name string) MpcParser {
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

// PrintAst parses an input string for its AST representation
func PrintAst(input string, mpcParser MpcParser) {
	r, err := MpcParse(input, mpcParser)
	if err != nil {
		C.mpc_err_print(C.get_error(&r))
		C.mpc_err_delete(C.get_error(&r))
	} else {
		C.mpc_ast_print(C.get_output(&r))
		MpcAstDelete(&r)
	}
}
