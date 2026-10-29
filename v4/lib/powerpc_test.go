// Copyright 2026 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccgo // import "modernc.org/ccgo/v4/lib"

import (
	"math"
	"testing"
)

func TestPowerPCPointerWidth(t *testing.T) {
	for _, tc := range []struct {
		arch       string
		maxAlign   int
		maxUintptr uint64
	}{
		{"ppc", 4, math.MaxUint32},
		{"ppc64", 8, math.MaxUint64},
		{"386", 4, math.MaxUint32},
		{"arm", 4, math.MaxUint32},
		{"amd64", 8, math.MaxUint64},
	} {
		t.Run(tc.arch, func(t *testing.T) {
			task := &Task{goarch: tc.arch}
			if got := newCtx(task, nil).maxAlign; got != tc.maxAlign {
				t.Errorf("maxAlign = %d, want %d", got, tc.maxAlign)
			}
			linker, err := newLinker(task, nil)
			if err != nil {
				t.Fatal(err)
			}
			if linker.maxUintptr != tc.maxUintptr {
				t.Errorf("maxUintptr = %d, want %d", linker.maxUintptr, tc.maxUintptr)
			}
		})
	}
}
