// Copyright 2026 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccgo // import "modernc.org/ccgo/v4/lib"

import (
	"testing"
)

// TestNestedStatementExpressions verifies that the value of a statement
// expression nested in another one reaches every place it is used, not only
// the last statement of the enclosing one, and that a value an expression
// statement inside a statement expression discards does not become a Go
// statement of its own.
//
// gcc's <stdatomic.h> atomic_load and atomic_store are statement expressions,
// so storing a value computed from two loads nests them in an initializer.
// sqlite-multiwriter has
//
//	atomic_store(&sh->log_pos, MW_LOG_POS(atomic_load(&sh->sl_seg), atomic_load(&sh->sl_end)));
//
// and the Go code for it had empty operands where the loaded values belong.
func TestNestedStatementExpressions(t *testing.T) {
	const src = `
#include <stdio.h>

// The shapes of atomic_load and atomic_store in gcc's <stdatomic.h>.
#define LOAD(p) ({ __auto_type __atomic_load_ptr = (p); __typeof__((void)0, *__atomic_load_ptr) __atomic_load_tmp; __atomic_load(__atomic_load_ptr, &__atomic_load_tmp, __ATOMIC_SEQ_CST); __atomic_load_tmp; })
#define STORE(p, v) ({ __auto_type __atomic_store_ptr = (p); __typeof__((void)0, *__atomic_store_ptr) __atomic_store_tmp = (v); __atomic_store(__atomic_store_ptr, &__atomic_store_tmp, __ATOMIC_SEQ_CST); })

struct s { unsigned long long pos; unsigned seg; unsigned long long end; };

static int g(int x) { return x + 100; }

int main(void) {
	struct s s = { 0, 3, 5 };
	int a = 3, b = 4, arr[] = { 7, 8, 9 };

	STORE(&s.pos, ((unsigned long long)LOAD(&s.seg) << 40) | LOAD(&s.end));
	printf("%x %x\n", (unsigned)(s.pos >> 32), (unsigned)s.pos);

	int v1 = ({ int t = ({ a; }) * 10 + ({ b; }); t; });                        // operands of an initializer
	int v2 = ({ int t = a; ({ t + b; }); });                                     // the last statement
	int v3 = ({ ({ a + 1; }); ({ b++; }); b; });                                 // discarded
	long v4 = ({ (long)({ a; }) + 1; });                                         // under a cast in the last statement
	int v5 = ({ g(({ a + b; })); });                                             // an argument
	int v6 = ({ ({ ({ 3; }) + 1; }) * 2; });                                     // three levels
	int v7 = ({ int r = 0; if (({ a; }) == 3) r = 1; while (({ r; }) < 5) r++; r; }); // conditions
	int v8 = ({ ({ (void)a; }); 8; });                                           // void
	int v9 = ({ a; a + b; ({ a; }) + ({ b; }); 9; });                            // discarded values
	unsigned v10 = ({ ({ &s; })->seg + 1; });                                    // a pointer
	int v11 = ({ arr[({ 1; })] + ({ arr; })[2]; });                              // an index, an array
	unsigned v12 = ({ ({ s; }).seg; });                                          // a struct
	printf("%d %d %d %d %d %d %d %d %d %u %d %u %d\n", v1, v2, v3, (int)v4, v5, v6, v7, v8, v9, v10, v11, v12, b);
	return 0;
}
`
	const exp = `300 5
34 7 5 4 108 8 5 8 9 4 17 3 5`
	if g := testHostCCvsCcgo(t, src); g != exp {
		t.Fatalf("got\n%s\nexpected\n%s", g, exp)
	}
}
