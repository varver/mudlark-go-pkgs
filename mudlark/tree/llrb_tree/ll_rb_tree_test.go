// Copyright 2010 -- Peter Williams, all rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package llrb_tree

import (
	"testing"
	"rand"
	"reflect"
	"fmt"
)

type Int int

func (i Int) Less(other interface{}) bool {
	return int(i) < int(other.(Int))
}

type Real float64

func (r Real) Less(other interface{}) bool {
	return float64(r) < float64(other.(Real))
}

func print_node(node *ll_rb_node) {
	if node == nil { return }
	fmt.Printf("%v\n", node)
	print_node(node.left)
	print_node(node.right)
}

func TestMakell_rb_tree(t *testing.T) {
	var tree ll_rb_tree
	if reflect.Typeof(tree).String() != "llrb_tree.ll_rb_tree" {
		t.Errorf("Expected type \"llrb_tree.ll_rb_tree\": got %v", reflect.Typeof(tree).String())
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
	if reflect.Typeof(tree).String() != "*llrb_tree.ll_rb_tree" {
		t.Errorf("Expected type \"*llrb_tree.ll_rb_tree\": got %v", reflect.Typeof(tree).String())
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
	for i := 0; i < 10; i++ {
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
				t.Errorf("Count unchanged (insert i): Expected %v got: %v", tsz + 1, tree.count)
			}
		}
		if iin, _ = tree.find(iitem); !iin {
			t.Errorf("Inserted %v not found", iitem)
			failures++
		}
	}
	if failures != 0 {
		t.Errorf("%v failures", failures)
	}
}

func TestMakeinsert_keep_duplicates(t *testing.T) {
	var tree ll_rb_tree
	var failures int
	var duplicates_found bool
	tree.keep_duplicates = true
	for i := 0; i < 1000; i++ {
		iitem := Int(rand.Intn(800))
		if found, _ := tree.find(iitem); found {
			duplicates_found = true
		}
		tsz := tree.count
		tree.insert(iitem)
		if tsz + 1 != tree.count {
			t.Errorf("Count unchanged (insert i): Expected %v got: %v", tsz + 1, tree.count)
		}
		if iin, _ := tree.find(iitem); !iin {
			t.Errorf("Inserted %v not found", iitem)
			failures++
		}
	}
	if failures != 0 {
		t.Errorf("%v failures", failures)
	}
	if !duplicates_found {
		t.Errorf("Test invalid: no duplicates inserted")
	}
}

func TestMakeiterate(t *testing.T) {
	var tree ll_rb_tree
	var count int
	for i := 0; i < 10000; i++ {
		tree.insert(Int(rand.Int()))
		count++
	}
	for item := range tree.iterator(PRE_ORDER) {
		if item.Less(Int(0)) {
			// shut compiler up
		}
		count--
	}
	if count != 0 {
		t.Errorf("%v count", count)
	}
}

func TestMakeiterate_in_order(t *testing.T) {
	var tree ll_rb_tree
	var count int
	for i := 0; i < 10000; i++ {
		tree.insert(Int(rand.Int()))
		count++
	}
	max_count := count
	lastItem := Int(0)
	for item := range tree.iterator(IN_ORDER) {
		if count < max_count && item.Less(lastItem) {
			t.Errorf("%v !< %v", item, lastItem)
		}
		count--
	}
	if count != 0 {
		t.Errorf("%v count", count)
	}
}

func max_depth(node *ll_rb_node) uint {
	if node == nil { return 0 }
	ld := max_depth(node.left)
	rd := max_depth(node.right)
	if ld > rd {
		return ld + 1
	}
	return rd + 1
}

// test that depth of tree doesn't exceed 2 * log2(cardinality) using:
//		random (best case) input
//		sequential (worst case) input
func TestMakedepth_properties(t *testing.T) {
	var tree_sequential, tree_reverse, tree_random ll_rb_tree
	var i int
	var max_depth_sequential, max_depth_reverse, max_depth_random uint
	for n := uint(1); n < 16; n++ {
		N := 1 << n
		for ; i < N; i++ {
			tree_sequential.insert(Int(i))
			tree_reverse.insert(Int(N - i))
			tree_random.insert(Int(rand.Int()))
		}
		max_depth_sequential = max_depth(tree_sequential.root)
		max_depth_reverse = max_depth(tree_reverse.root)
		max_depth_random  = max_depth(tree_random.root)
		if max_depth_sequential > 2 * n || max_depth_reverse > 2 * n || max_depth_random > 2 * n {
			t.Errorf("%v : %v : %v : %v\n", n, i, max_depth_sequential, max_depth_random)
		}
	}
}

