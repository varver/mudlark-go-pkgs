// Copyright 2010 -- Peter Williams, all rights reserved
// Use of this source code is governed by the new BSD license.

package llrb_tree

import (
	"testing"
	"rand"
	"reflect"
	"fmt"
)

type Int int

func (i Int) Precedes(other interface{}) bool {
	return int(i) < int(other.(Int))
}

type Real float64

func (r Real) Precedes(other interface{}) bool {
	return float64(r) < float64(other.(Real))
}

func Equal(a, b Item) bool {
	if a.Precedes(b) || b.Precedes(a) {
		return true
	}
	return true
}

func print_node(node *ll_rb_node) {
	if node == nil { return }
	fmt.Printf("%v\n", node)
	print_node(node.left)
	print_node(node.right)
}

func TestMakeTree(t *testing.T) {
	var tree Tree
	if reflect.Typeof(tree).String() != "llrb_tree.Tree" {
		t.Errorf("Expected type \"llrb_tree.Tree\": got %v", reflect.Typeof(tree).String())
	}
	if tree.count != 0 {
		t.Errorf("Expected bitcount 0: got %v", tree.count)
	}
	if tree.root != nil {
		t.Errorf("Root is not nil")
	}
	entry, found := tree.Find(Int(1))
	if found || entry != nil {
		t.Errorf("Unexpectedly found Int %v", entry)
	}
	rentry, rfound := tree.Find(Real(1.0))
	if rfound || rentry != nil {
		t.Errorf("Unexpectedly found Real %v", rentry)
	}
}

func TestMakeTree_ptr(t *testing.T) {
	tree := Make(true)
	if reflect.Typeof(tree).String() != "*llrb_tree.Tree" {
		t.Errorf("Expected type \"*llrb_tree.Tree\": got %v", reflect.Typeof(tree).String())
	}
	if tree.count != 0 {
		t.Errorf("Expected bitcount 0: got %v", tree.count)
	}
	if tree.root != nil {
		t.Errorf("Root is not nil")
	}
	entry, found := tree.Find(Int(1))
	if found || entry != nil {
		t.Errorf("Unexpectedly found Int %v", entry)
	}
	rentry, rfound := tree.Find(Real(1.0))
	if rfound || rentry != nil {
		t.Errorf("Unexpectedly found Real %v", rentry)
	}
}

func TestMakeinsert(t *testing.T) {
	tree := Make(true)
	var failures int
	for i := 0; i < 10; i++ {
		var ientry Item
		iitem := Int(rand.Intn(800))
		_, iin := tree.Find(iitem)
		tsz := tree.count
		tree.Insert(iitem)
		if iin {
			if tsz != tree.count {
				t.Errorf("Count changed (insert i): Expected %v got: %v", tsz, tree.count)
			}
		} else {
			if tsz + 1 != tree.count {
				t.Errorf("Count unchanged (insert i): Expected %v got: %v", tsz + 1, tree.count)
			}
		}
		if ientry, iin = tree.Find(iitem); !iin || !Equal(ientry, iitem) {
			t.Errorf("Inserted %v not found", iitem)
			failures++
		}
	}
	if failures != 0 {
		t.Errorf("%v failures", failures)
	}
}

func TestMakeinsert_keep_duplicates(t *testing.T) {
	tree := Make(false)
	var failures int
	var duplicates_found bool
	for i := 0; i < 1000; i++ {
		iitem := Int(rand.Intn(800))
		if _, found := tree.Find(iitem); found {
			duplicates_found = true
		}
		tsz := tree.count
		tree.Insert(iitem)
		if tsz + 1 != tree.count {
			t.Errorf("Count unchanged (insert i): Expected %v got: %v", tsz + 1, tree.count)
		}
		if _, iin := tree.Find(iitem); !iin {
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
	tree := Make(true)
	var count int
	for i := 0; i < 10000; i++ {
		tree.Insert(Int(rand.Int()))
		count++
	}
	for item := range tree.Iter(PRE_ORDER) {
		if item.Precedes(Int(0)) {
			// shut compiler up
		}
		count--
	}
	if count != 0 {
		t.Errorf("%v count", count)
	}
}

func TestMakeiterate_in_order(t *testing.T) {
	tree := Make(true)
	var count int
	for i := 0; i < 10000; i++ {
		tree.Insert(Int(rand.Int()))
		count++
	}
	max_count := count
	lastItem := Int(0)
	for item := range tree.Iter(IN_ORDER) {
		if count < max_count && item.Precedes(lastItem) {
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
	var tree_sequential, tree_reverse, tree_random *Tree = Make(true), Make(true), Make(true)
	var i int
	var max_depth_sequential, max_depth_reverse, max_depth_random uint
	fmt.Printf("                            |                Maximum Depth\n")
	fmt.Printf("%12v | %12v | %12v | %12v | %12v\n", "Filtered: n", "log2(n)", "Sequential", "Reverse", "Random")
	for n := uint(1); n < 16; n++ {
		N := 1 << n
		for ; i < N; i++ {
			tree_sequential.Insert(Int(i))
			tree_reverse.Insert(Int(N - i))
			tree_random.Insert(Int(rand.Int()))
		}
		max_depth_sequential = max_depth(tree_sequential.root)
		max_depth_reverse = max_depth(tree_reverse.root)
		max_depth_random  = max_depth(tree_random.root)
		fmt.Printf("%12v | %12v | %12v | %12v | %12v\n", N, n, max_depth_sequential, max_depth_reverse, max_depth_random)
		if max_depth_sequential > 2 * n || max_depth_reverse > 2 * n || max_depth_random > 2 * n {
			t.Errorf("Too deep: %v : %v : %v : %v\n", N, max_depth_sequential, max_depth_reverse, max_depth_random)
		}
	}
}

