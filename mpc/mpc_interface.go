package mpc

// #cgo LDFLAGS: -lm
// #include "mpc_interface.h"
import "C"
import (
	"errors"
	"unsafe"
)

// AstPtr is a pointer to a mpc-generated AST
type AstPtr *C.mpc_ast_t

// ErrorPtr is a pointer to an mpc error
type ErrorPtr *C.mpc_err_t

// ParserPtr is a pointer to a parser created with mpc
type ParserPtr *C.struct_mpc_parser_t

// MpcResult is a union that returns either an output or error
type MpcResult C.mpc_result_t

// GetNumChildren accesses the children_num of AstPtr
func GetNumChildren(node AstPtr) int {
	return int(C.get_children_num(node))
}

// GetChild accesses the child at a specific index in AstPtr
func GetChild(node AstPtr, index int) AstPtr {
	return C.get_child(node, C.int(index))
}

// GetContents accesses the contents of an AstPtr
func GetContents(node AstPtr) string {
	return C.GoString(node.contents)
}

// GetOperator accesses a AstPtr's child node representing an operator
func GetOperator(node AstPtr) string {
	return GetContents(GetChild(node, 1))
}

// GetOutput accesses the output field of an input MpcResult
func GetOutput(result *C.mpc_result_t) AstPtr {
	cast := (*AstPtr)(unsafe.Pointer(result))
	return *cast
}

// GetError accesses the error field of an input MpcResult
func GetError(result *C.mpc_result_t) ErrorPtr {
	cast := (*ErrorPtr)(unsafe.Pointer(result))
	return *cast
}

// GetTag accesses the tag field of an AstPtr
func GetTag(node AstPtr) string {
	return C.GoString(node.tag)
}

// DeleteAstPtr cleans up an AST contained in a MpcResult
func DeleteAstPtr(result *C.mpc_result_t) {
	C.mpc_ast_delete(GetOutput(result))
}

// MpcaLang uses a language definition to generate parsers
func MpcaLang(language string, parsers ...ParserPtr) {
	cLanguage := C.CString(language)
	defer C.free(unsafe.Pointer(cLanguage))
	C.mpca_lang_if(C.MPCA_LANG_DEFAULT, cLanguage,
		parsers[0], parsers[1], parsers[2], parsers[3], parsers[4], parsers[5], parsers[6], parsers[7])
}

// MpcCleanup calls mpc's cleanup function indirectly, using a wrapper for the variadic args
func MpcCleanup(parsers ...ParserPtr) {
	C.mpc_cleanup_if(C.int(len(parsers)),
		parsers[0], parsers[1], parsers[2], parsers[3], parsers[4], parsers[5], parsers[6], parsers[7])
}

// MpcNew returns a pointer to an initiated mpc parser
func MpcNew(name string) ParserPtr {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return C.mpc_new(cName)
}

// MpcParse takes an input string and generates an MpcResult
func MpcParse(input string, parser ParserPtr) (C.mpc_result_t, error) {
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
	C.mpc_err_delete(GetError(result))
}

// MpcErrPrint prints the reason for failed AST parsing
func MpcErrPrint(result *C.mpc_result_t) {
	C.mpc_err_print(GetError(result))
}

// PrintAst parses an input string for its AST representation
func PrintAst(input string, mpcParser ParserPtr) {
	r, err := MpcParse(input, mpcParser)
	if err != nil {
		MpcErrPrint(&r)
		MpcErrDelete(&r)
	} else {
		C.mpc_ast_print(GetOutput(&r))
		DeleteAstPtr(&r)
	}
}

// MpcfEscape inserts escape characters into an input string
func MpcfEscape(input string) string {
	cInput := C.CString(input)
	// mpcf_escape calls free() on the input string. Refer to mpc.c for context.
	ret := C.mpcf_escape(unsafe.Pointer(cInput))
	return C.GoString((*C.char)(ret))
}

// MpcfUnescape converts escape characters into encoded equivalents
func MpcfUnescape(input string) string {
	cInput := C.CString(input)
	// mpcf_unescape calls free() on the input string. Refer to mpc.c for context.
	ret := C.mpcf_unescape(unsafe.Pointer(cInput))
	return C.GoString((*C.char)(ret))
}

// MpcParseContents parses the contents of a file
func MpcParseContents(input string, parser ParserPtr) (*C.mpc_result_t, error) {
	cInput := C.CString(input)
	defer C.free(unsafe.Pointer(cInput))
	result := new(C.mpc_result_t)
	var err error
	ret := C.mpc_parse_contents(cInput, parser, result)
	if ret == 0 {
		err = errors.New("mpc: failed to parse input file")
	}
	return result, err
}

// MpcErrString gets the error string from an mpc result
func MpcErrString(result *C.mpc_result_t) string {
	err := GetError(result)
	if err == nil {
		return "<Failed to load error>"
	}
	cErrMsg := C.mpc_err_string(err)
	MpcErrDelete(result)
	return C.GoString(cErrMsg)
}
