#ifndef mpc_interface_h
#define mpc_interface_h

#include "mpc.h"
/* CHEAP HACK ALERT!
 * I know that including the c file is a terrible idea, but I've been spending
 * an embarrassing amount of time trying to figure out why cgo won't link the
 * "undefined" symbols from that file. I'll fix this once my subconscious
 * discovers my bone-headed omission, or when I break down and ask for help.
 */
#include "mpc.c"

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
  mpc_parser_t* parser4
)
{
  mpc_cleanup(n, parser1, parser2, parser3, parser4);
}

inline mpc_err_t* mpca_lang_if
(
  int flags,
  const char *language,
  mpc_parser_t* parser1, // variadic args
  mpc_parser_t* parser2,
  mpc_parser_t* parser3,
  mpc_parser_t* parser4
)
{
  return mpca_lang(flags, language, parser1, parser2, parser3, parser4);
}

#endif

