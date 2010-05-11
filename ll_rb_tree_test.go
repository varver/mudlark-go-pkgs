// Copyright 2010 -- Peter Williams, all rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heteroset

import (
	"testing"
//	"rand"
	"reflect"
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
	if tree.find(Int(1)) {
		t.Errorf("Unexpectedly found Int")
	}
	if tree.find(Real(1.0)) {
		t.Errorf("Unexpectedly found Real")
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
	if tree.find(Int(1)) {
		t.Errorf("Unexpectedly found Int")
	}
	if tree.find(Real(1.0)) {
		t.Errorf("Unexpectedly found Real")
	}
}

