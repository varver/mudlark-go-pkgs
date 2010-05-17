// Copyright 2010 -- Peter Williams, all rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sort

import (
	"testing"
	"rand"
//	"fmt"
//	"reflect"
)

type Int int

func (i Int) Less(other interface{}) bool {
	return int(i) < int(other.(Int))
}

type Real float64

func (r Real) Less(other interface{}) bool {
	return float64(r) < float64(other.(Real))
}

func TestMakeSortSlice(t *testing.T) {
	const sz = 1000
	ints := make([]Item, sz)
	for i := 0; i < sz; i++ {
		ints[i] = Int(rand.Intn(8 * sz / 10))
	}
	count := 0
	var lasti Item
	for _, i := range SortSlice(ints) {
		if count != 0 && i.Less(lasti) {
			t.Errorf("Unexpected order: %v : %v", i, lasti)
		}
		count++
		lasti = i
	}
	if count != sz {
		t.Errorf("Expected count %v: got %v", sz, count)
	}
}

