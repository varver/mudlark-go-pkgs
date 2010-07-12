// Copyright 2010 -- Peter Williams, all rights reserved
// Use of this source code is governed by the new BSD license.

package heteroset;

import (
	"testing";
	"rand";
	"reflect";
	"fmt";
)

type Int int;

func (i Int) Precedes(other interface{}) bool {
	return int(i) < int(other.(Int));
};

type Real float64;

func (r Real) Precedes(other interface{}) bool {
	return float64(r) < float64(other.(Real));
};

func print_node(node *ll_rb_node) {
	if node == nil { return; };
	fmt.Printf("%v\n", node);
	print_node(node.left);
	print_node(node.right);
};

func max_depth(node *ll_rb_node) uint {
	if node == nil { return 0; };
	ld := max_depth(node.left);
	rd := max_depth(node.right);
	if ld > rd {
		return ld + 1;
	};
	return rd + 1;
};

func TestMakeSet(t *testing.T) {
	set := New();
	if reflect.Typeof(set).String() != "*heteroset.Set" {
		t.Errorf("Expected type \"*heteroset.Set\": got %v", reflect.Typeof(set).String());
	};
	if set.Cardinality() != 0 {
		t.Errorf("Expected bitcount 0: got %v", set.Cardinality());
	};
	if set.root != nil {
		t.Errorf("Root is not nil");
	};
	has := set.Has(Int(1));
	if has {
		t.Errorf("Unexpectedly has Int");
	};
	if max_depth(set.root) != 0 {
		t.Errorf("Expected 0 max depth got: %v", max_depth(set.root));
	};
	has = set.Has(Real(1.0));
	if has {
		t.Errorf("Unexpectedly has Real");
	};
	if max_depth(set.root) != 0 {
		t.Errorf("Expected 0 max depth got: %v", max_depth(set.root));
	};
};

func TestMakeSetWithArgs(t *testing.T) {
	set := New(Int(1), Int(2), Int(2), Real(3), Int(4), Real(4));
	if reflect.Typeof(set).String() != "*heteroset.Set" {
		t.Errorf("Expected type \"*heteroset.Set\": got %v", reflect.Typeof(set).String());
	};
	if set.Cardinality() != 5 {
		t.Errorf("Expected count 5: got %v", set.Cardinality());
	};
	if set.root == nil {
		t.Errorf("Root is nil");
	};
	has := set.Has(Int(1));
	if !has {
		t.Errorf("Denies having Int(1)");
	};
	if max_depth(set.root) == 0 {
		t.Errorf("Expected 0 max depth got: %v", max_depth(set.root));
	};
	has = set.Has(Real(1.0));
	if has {
		t.Errorf("Unexpectedly has Real(1.0)");
	};
};

func TestMakeinsert(t *testing.T) {
	set := New();
	var failures int;
	for i := 0; i < 1000; i++ {
		iitem := Int(rand.Intn(800));
		iin := set.Has(iitem);
		tsz := set.Cardinality();
		set.Add(iitem);
		if iin {
			if tsz != set.Cardinality() {
				t.Errorf("Count changed (insert i): Expected %v got: %v", tsz, set.Cardinality());
			};
		} else {
			if tsz + 1 != set.Cardinality() {
				t.Errorf("Count uchanged (insert i): Expected %v got: %v", tsz + 1, set.Cardinality());
			};
		};
		if iin = set.Has(iitem); !iin {
			t.Errorf("Inserted %v not has", iitem);
			failures++;
		};
		ritem := Real(rand.Float64());
		rin := set.Has(ritem);
		tsz = set.Cardinality();
		set.Add(ritem);
		if rin {
			if tsz != set.Cardinality() {
				t.Errorf("Count changed (insert i): Expected %v got: %v", tsz, set.Cardinality());
			};
		} else {
			if tsz + 1 != set.Cardinality() {
				t.Errorf("Count uchanged (insert i): Expected %v got: %v", tsz + 1, set.Cardinality());
			};
		};
		if rin = set.Has(ritem); !rin {
			t.Errorf("Inserted %v not has", ritem);
			failures++;
		};
	};
	if failures != 0 {
		t.Errorf("%v failures", failures);
	};
};

func TestMakeiterate(t *testing.T) {
	set := New();
	var count int;
	for i := 0; i < 10000; i++ {
		set.Add(Int(rand.Int()));
		count++;
		set.Add(Real(rand.Float64()));
		count++;
	};
	for item := range set.Iter() {
		if cmp_type(item, Int(0)) == 0 {
			// shut compiler up
		};
		count--;
	};
	if count != 0 {
		t.Errorf("%v count", count);
	};
};

func TestMakeiterateAsync(t *testing.T) {
	set := New();
	var count int;
	for i := 0; i < 10000; i++ {
		set.Add(Int(rand.Int()));
		count++;
		set.Add(Real(rand.Float64()));
		count++;
	};
	for item := range set.IterAsync() {
		if cmp_type(item, Int(0)) == 0 {
			// shut compiler up
		};
		count--;
	};
	if count != 0 {
		t.Errorf("%v count", count);
	};
};

// test that depth of set doesn't exceed 2 * log2(cardinality) using:
//		random (best case) input
//		sequential (worst case) input
func TestMakedepth_properties(t *testing.T) {
	set_sequential, set_reverse, set_random := New(), New(), New();
	var i int;
	var max_depth_sequential, max_depth_reverse, max_depth_random uint;
	for n := uint(1); n < 16; n++ {
		N := 1 << n;
		for ; i < N; i++ {
			set_sequential.Add(Int(i));
			set_reverse.Add(Int(N - i));
			set_random.Add(Int(rand.Int()));
		};
		max_depth_sequential = max_depth(set_sequential.root);
		max_depth_reverse = max_depth(set_reverse.root);
		max_depth_random  = max_depth(set_random.root);
		if max_depth_sequential > 2 * n || max_depth_reverse > 2 * n || max_depth_random > 2 * n {
			t.Errorf("%v : %v : %v : %v\n", n, i, max_depth_sequential, max_depth_random);
		};
	};
};

func make_Int_set_serial(begin, end Int) (set *Set) {
	set = New();
	for i := begin; i <= end; i++ {
		set.Add(i);
	};
	return;
};

func TestDisjointIntersect(t *testing.T) {
	setA := make_Int_set_serial(-100, 0);
	setB := make_Int_set_serial(1, 100);
	setC := make_Int_set_serial(-50, 50);
	if Intersect(setA, setB) {
		t.Errorf("setA and setB should be disjoint");
	};
	if !Intersect(setA, setC) {
		t.Errorf("setA and setC should intersect");
	};
	if Intersect(setA, setB) && Disjoint(setA, setB) {
		t.Errorf("Intersect(A, B) and Disjoint(A, B) should be mutually exclusive");
	};
	if Intersect(setA, setC) && Disjoint(setA, setC) {
		t.Errorf("Intersect(A, B) and Disjoint(A, B) should be mutually exclusive");
	};
	if Intersect(setA, setB) != Intersect(setB, setA) {
		t.Errorf("Intersect(A, B) and Intersect(B, A) should be equal");
	};
	if Intersect(setA, setC) != Intersect(setC, setA) {
		t.Errorf("Intersect(A, C) and Intersect(C, A) should be equal");
	};
	if Disjoint(setA, setB) != Disjoint(setB, setA) {
		t.Errorf("Disjoint(A, B) and Disjoint(B, A) should be equal");
	};
	if Disjoint(setA, setC) != Disjoint(setC, setA) {
		t.Errorf("Disjoint(A, C) and Disjoint(C, A) should be equal");
	};
};

func TestUnion(t *testing.T) {
	setA := make_Int_set_serial(-100, 0);
	setB := make_Int_set_serial(1, 100);
	setC := make_Int_set_serial(-50, 50);
	setAuB := Union(setA, setB);
	setAuC := Union(setA, setC);
	if !Intersect(setA, setAuB) || !Intersect(setB, setAuB) {
		t.Errorf("setAuB should intersect with both setA and SetB");
	};
	if setAuB.Cardinality() != setA.Cardinality() + setB.Cardinality() {
		t.Errorf("Cardinality of a union of disjoint sets should be the sum of their cardinalities");
	};
	for item := range setAuB.Iter() {
		if !setA.Has(item) && !setB.Has(item) {
			t.Errorf("Items in setAuB should be in either setA or setB");
		};
	};
	if !Intersect(setA, setAuC) || !Intersect(setC, setAuC) {
		t.Errorf("setAuC should intersect with both setA and SetC");
	};
	var incommon uint;
	for item := range setAuC.Iter() {
		if setA.Has(item) {
			if setC.Has(item) {
				incommon++;
			};
		} else if !setC.Has(item) {
			t.Errorf("Items in setAuC should be in either setA or setC");
		};
	};
	if setAuC.Cardinality() != setA.Cardinality() + setC.Cardinality() - incommon {
		t.Errorf("Cardinality of a union of intesecting sets should be the sum of their cardinalities minus the size of their intersection");
	};
};

func TestIntersection(t *testing.T) {
	setA := make_Int_set_serial(-100, 0);
	setB := make_Int_set_serial(1, 100);
	setC := make_Int_set_serial(-50, 50);
	setAiB := Intersection(setA, setB);
	setAiC := Intersection(setA, setC);
	if (Intersect(setA, setAiB) && !Intersect(setB, setAiB)) || (!Intersect(setA, setAiB) && Intersect(setB, setAiB)) {
		t.Errorf("setAiB should intersect with both setA and SetB or neither");
	};
	if setAiB.Cardinality() != 0 {
		t.Errorf("Cardinality of an intersection of disjoint sets should be 0");
	};
	for item := range setAiB.Iter() {
		if !setA.Has(item) || !setB.Has(item) {
			t.Errorf("Items in setAiB should be in both setA and setB");
		};
	};
	if !Intersect(setA, setAiC) || !Intersect(setC, setAiC) {
		t.Errorf("setAiC should intersect with both setA and SetC");
	};
	for item := range setAiC.Iter() {
		if setA.Has(item) {
			if !setC.Has(item) {
				t.Errorf("Items in setAiC should be in both setA and setC");
			};
		} else if setC.Has(item) {
			t.Errorf("Items in setAiC should be in both setA and setC");
		};
	};
	if setAiC.Cardinality() != 51 {
		t.Errorf("Cardinality of an intersection of intesecting sets should be the size of their intersection");
	};
};

