// Copyright 2026 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccgo // import "modernc.org/ccgo/v4/lib"

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestOveralignedTagDeclaration checks that a tag declaration of a type Go
// cannot align, which declares no object, transpiles without
// -ignore-unsupported-alignment, and that an object of the type is still
// reported. mingw's excpt.h declares struct _CONTEXT, which winnt.h defines
// with DECLSPEC_ALIGN(16), so every program including <windows.h> for
// windows/amd64 has such a declaration.
func TestOveralignedTagDeclaration(t *testing.T) {
	const def = "struct s { long long a, b; } __attribute__((aligned(16)));\n"
	for _, tc := range []struct {
		src string
		err string
	}{
		{def + "struct s;\n", ""},
		{def + "struct s v;\n", "unsupported alignment 16"},
		{def + "int f(void) { struct s v; return sizeof v; }\n", "unsupported alignment 16"},
	} {
		dir := t.TempDir()
		cFile := filepath.Join(dir, "test.c")
		if err := os.WriteFile(cFile, []byte(tc.src), 0644); err != nil {
			t.Fatal(err)
		}

		var stdout, stderr bytes.Buffer
		err := NewTask(goos, goarch, []string{"ccgo", "-c", "-o", filepath.Join(dir, "test.o.go"), cFile}, &stdout, &stderr, nil).Main()
		switch {
		case tc.err == "" && err != nil:
			t.Errorf("%q: unexpected error: %v", tc.src, err)
		case tc.err != "" && (err == nil || !strings.Contains(err.Error(), tc.err)):
			t.Errorf("%q: got error %v, want one containing %q", tc.src, err, tc.err)
		}
	}
}
