// Copyright 2010 -- Peter Williams, all rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heteroset

import (
	"testing"
	"rand"
	"reflect"
	//"fmt"
)

type Int int

func (i Int) Compare(other Item) int {
	return int(i) - int(other.(Int))
}

type Real float64

func (r Real) Compare(other Item) int {
	if float64(r) < float64(other.(Real)) {
		return -1
	} else if float64(r) > float64(other.(Real)) {
		return 1
	}
	return 0
}

func TestMakell_rb_tree(t *testing.T) {
	var tree ll_rb_tree
	if reflect.Typeof(tree).String() != "heteroset.ll_rb_tree" {
		t.Errorf("Expected type \"heteroset.ll_rb_tree\": got %v", reflect.Typeof(tree).String())
	}
	if tree.count != 0 {
		t.Errorf("Expected bitcount 0: got %v", tree.count)
	}
	if tree.root != nil {
		t.Errorf("Root is not nil")
	}
	found, iterations := tree.find(Int(1))
	if found {
		t.Errorf("Unexpectedly found Int")
	}
	if iterations != 0 {
		t.Errorf("Expected 0 iteretions got: %v", iterations)
	}
	found, iterations = tree.find(Real(1.0))
	if found {
		t.Errorf("Unexpectedly found Real")
	}
	if iterations != 0 {
		t.Errorf("Expected 0 iteretions got: %v", iterations)
	}
}

func TestMakell_rb_tree_ptr(t *testing.T) {
	tree := new(ll_rb_tree)
	if reflect.Typeof(tree).String() != "*heteroset.ll_rb_tree" {
		t.Errorf("Expected type \"*heteroset.ll_rb_tree\": got %v", reflect.Typeof(tree).String())
	}
	if tree.count != 0 {
		t.Errorf("Expected bitcount 0: got %v", tree.count)
	}
	if tree.root != nil {
		t.Errorf("Root is not nil")
	}
	found, iterations := tree.find(Int(1))
	if found {
		t.Errorf("Unexpectedly found Int")
	}
	if iterations != 0 {
		t.Errorf("Expected 0 iteretions got: %v", iterations)
	}
	found, iterations = tree.find(Real(1.0))
	if found {
		t.Errorf("Unexpectedly found Real")
	}
	if iterations != 0 {
		t.Errorf("Expected 0 iteretions got: %v", iterations)
	}
}

func TestMakeinsert(t *testing.T) {
	var tree ll_rb_tree
	var failures int
	for i := 0; i < 1000; i++ {
		iitem := Int(rand.Intn(800))
		iin, _ := tree.find(iitem)
		tsz := tree.count
		tree.insert(iitem)
		if iin {
			if tsz != tree.count {
				t.Errorf("Count changed (insert i): Expected %v got: %v", tsz, tree.count)
			}
		} else {
			if tsz + 1 != tree.count {
				t.Errorf("Count uchanged (insert i): Expected %v got: %v", tsz + 1, tree.count)
			}
		}
		if iin, _ = tree.find(iitem); !iin {
			t.Errorf("Inserted %v not found", iitem)
			failures++
		}
		ritem := Real(rand.Float64())
		rin, _:= tree.find(ritem)
		tsz = tree.count
		tree.insert(ritem)
		if rin {
			if tsz != tree.count {
				t.Errorf("Count changed (insert i): Expected %v got: %v", tsz, tree.count)
			}
		} else {
			if tsz + 1 != tree.count {
				t.Errorf("Count uchanged (insert i): Expected %v got: %v", tsz + 1, tree.count)
			}
		}
		if rin, _ = tree.find(ritem); !rin {
			t.Errorf("Inserted %v not found", ritem)
			failures++
		}
	}
	if failures != 0 {
		t.Errorf("%v failures", failures)
	}
}

