#ifndef mpc_interface_h
#define mpc_interface_h

#include "mpc.h"

inline int get_children_num(mpc_ast_t* node)
{
  return node->children_num;
}

inline mpc_ast_t* get_child(mpc_ast_t* node, int index)
{
  return node->children[index]; // index into double pointer
}

inline mpc_err_t* get_error(mpc_result_t* result)
{
  return (result == NULL) ? NULL : result->error;
}

inline mpc_ast_t* get_output(mpc_result_t* result)
{
  return (result == NULL) ? NULL : result->output;
}

inline void mpc_cleanup_if
(
  int n,
  mpc_parser_t* parser1, // variadic args
  mpc_parser_t* parser2,
  mpc_parser_t* parser3,
  mpc_parser_t* parser4,
  mpc_parser_t* parser5,
  mpc_parser_t* parser6
)
{
  mpc_cleanup(n, parser1, parser2, parser3, parser4, parser5, parser6);
}

inline mpc_err_t* mpca_lang_if
(
  int flags,
  const char *language,
  mpc_parser_t* parser1, // variadic args
  mpc_parser_t* parser2,
  mpc_parser_t* parser3,
  mpc_parser_t* parser4,
  mpc_parser_t* parser5,
  mpc_parser_t* parser6
)
{
  return mpca_lang(flags, language, parser1, parser2, parser3, parser4, parser5, parser6);
}

#endif

