// Copyright 2010 -- Peter Williams, all rights reserved
// Use of this source code is governed by the new BSD license.

// The heteroset package implements heterogenous sets.
//
// (Uses 2-3 left Leaning Red Black Binary Trees as the low level component.
// Based on the Java implementation described by Robert Sedgewick
// in his paper entitled "Left-leaning Red-Black Trees"
// available at: <www.cs.princeton.edu/~rs/talks/LLRB/LLRB.pdf>.
// The principal difference (other than the conversion to Go) is that the items
// being inserted combine the roles of both key and value and the items
// being inserted do not have to be of the same type.)
package heteroset;

import "reflect";

// The type of potential set items must implement this interface and must
// satisfy the following formal requirements (where a, b and c are all
// instances of the same type):
//	 a.Precedes(b) implies !b.Precedes(a)
//	 a.Precedes(b) && b.Precedes(c) implies a.Precedes(c)
//	 !a.Precedes(b) && !b.Precedes(a) implies a == b
// This method will only be used when reflect.Typeof() the calling object
// matches reflect.Typeof() of other.
type Item interface {
	Precedes(other interface{}) bool;
};

// LLRB tree node
type ll_rb_node struct {
	item Item;
	left, right *ll_rb_node;
	red bool;
};

func new_ll_rb_node(item Item) *ll_rb_node {
	node := new(ll_rb_node);
	node.item = item;
	node.red = true;
	return node;
};

func min(a, b int) int { if a < b { return a }; return b; };

func cmp_string(a, b string) int {
	for i, lim := 0, min(len(a), len(b)); i < lim; i++ {
		if a[i] < b[i] {
			return -1;
		} else if a[i] > b[i] {
			return 1;
		};
	};
	return len(a) - len(b);
};

func cmp_type(a, b interface{}) int {
	ta := reflect.Typeof(a);
	tb := reflect.Typeof(b);
	if ta == tb {
		return 0;
	};
	if cp := cmp_string(ta.PkgPath(), tb.PkgPath()); cp != 0 {
		return cp;
	};
	return cmp_string(ta.Name(), tb.Name());
};

func (this *ll_rb_node) compare_item(item Item) int {
	if ct := cmp_type(this.item, item); ct != 0 {
		return ct;
	};
	if this.item.Precedes(item) {
		return -1;
	} else if item.Precedes(this.item) {
		return 1;
	};
	return 0;
};

func is_red(node *ll_rb_node) bool { return node != nil && node.red; };

func flip_colours(node *ll_rb_node) {
	node.red = !node.red;
	node.left.red = !node.left.red;
	node.right.red = !node.right.red;
};

func rotate_left(node *ll_rb_node) *ll_rb_node {
	tmp := node.right;
	node.right = tmp.left;
	tmp.left = node;
	tmp.red = node.red;
	node.red = true;
	return tmp;
};

func rotate_right(node *ll_rb_node) *ll_rb_node {
	tmp := node.left;
	node.left = tmp.right;
	tmp.right = node;
	tmp.red = node.red;
	node.red = true;
	return tmp;
};

func fix_up(node *ll_rb_node) *ll_rb_node {
	if is_red(node.right) && !is_red(node.left) {
		node = rotate_left(node);
	};
	if is_red(node.left) && is_red(node.left.left) {
		node = rotate_right(node);
	};
	if is_red(node.left) && is_red(node.right) {
		flip_colours(node);
	};
	return node;
};

func insert(node *ll_rb_node, item Item) (*ll_rb_node, bool) {
	if node == nil {
		return new_ll_rb_node(item), true;
	};
	inserted := false;
	switch cmp := node.compare_item(item); {
	case cmp > 0:
		node.left, inserted = insert(node.left, item);
	case cmp < 0:
		node.right, inserted = insert(node.right, item);
	default:
		// overwrite the existing equivalent item so that Sets are useful
		// with (key, value) items
		node.item = item;
	};
	return fix_up(node), inserted;
};

func move_red_left(node *ll_rb_node) *ll_rb_node {
	flip_colours(node);
	if (is_red(node.right.left)) {
		node.right = rotate_right(node.right);
		node = rotate_left(node);
		flip_colours(node);
	};
	return node;
};

func move_red_right(node *ll_rb_node) *ll_rb_node {
	flip_colours(node);
	if (is_red(node.left.left)) {
		node = rotate_right(node);
		flip_colours(node);
	};
	return node;
};

func delete_left_most(node *ll_rb_node) *ll_rb_node {
	if node.left == nil {
		return nil;
	};
	if !is_red(node.left) && !is_red(node.left.left) {
		node = move_red_left(node);
	};
	node.left = delete_left_most(node.left);
	return fix_up(node);
};

func delete(node *ll_rb_node, item Item) (*ll_rb_node, bool) {
	var deleted bool;
	if node.compare_item(item) > 0 {
		if !is_red(node.left) && !is_red(node.left.left) {
			node = move_red_left(node);
		};
		node.left, deleted = delete(node.left, item);
	} else {
		if is_red(node.left) {
			node = rotate_right(node);
		};
		if node.compare_item(item) == 0 && node.right == nil {
			return nil, true;
		};
		if !is_red(node.right) && !is_red(node.right.left) {
			node = move_red_right(node);
		};
		if node.compare_item(item) == 0 {
			left_most := node.right;
			for left_most.left != nil {
				left_most = left_most.left;
			};
			node.item = left_most.item;
			node.right = delete_left_most(node.right);
			deleted = true;
		} else {
			node.right, deleted = delete(node.right, item);
		};
	};
	return fix_up(node), deleted;
};

// Iteration using recursion is safe because the depth of the tree should never
// be greater than 2Log2(N) where N is the number of nodes in the tree and
// (in general) will be approximately Log2(N).

func iterate_inorder(node *ll_rb_node, c chan<- Item) {
	if node == nil {
		return;
	}
	iterate_inorder(node.left, c);
	c <- node.item;
	iterate_inorder(node.right, c);
};

func iterate(node *ll_rb_node, c chan<- Item) {
	iterate_inorder(node, c);
	close(c);
};

func copy(node *ll_rb_node) *ll_rb_node {
	if node == nil { return nil; };
	clone := new(ll_rb_node);
	clone.item = node.item;
	clone.red = node.red;
	clone.left = copy(node.left);
	clone.right = copy(node.right);
	return clone;
};

// Set is a set of hetrogeneous objects whos types implement the Item
// interface. Instances of Set must be created using New()
// before use.  E.g.:
//	var s Set = heteroset.New(item1, ....)
type Set struct {
	root *ll_rb_node;
	count uint;
};

// Make a Set. The optional Item parameters will be used to initialize the set's
// contents.
func New(items ...Item) (set *Set) {
	set = new(Set);
	for _, item := range items {
		set.Add(item);
	};
	return;
};

// Len returns the number of items in the set.
func (this *Set) Cardinality() uint {
	return this.count;
};

// Make a copy of this set.
func (this *Set) Copy() (set *Set) {
	set = new(Set);
	set.root = copy(this.root);
	set.count = this.count;
	return;
};

// Find an instance equal to item in the set.
// This function is useful in the case where the item has a (key, value)
// structure and only the key is used for implementing Precedes() for using
// a Set as a look up table.
func (this *Set) Find(item Item) (instance Item, found bool) {
	if this.count == 0 {
		return;
	};
	for node := this.root; node != nil && !found; {
		switch cmp := node.compare_item(item); {
		case cmp > 0:
			node = node.left;
		case cmp < 0:
			node = node.right;
		default:
			found = true;
			instance = node.item;
		};
	};
	return;
};

// Is there an instance equal to item in the set.
func (this *Set) Has(item Item) (has bool) {
	_, has = this.Find(item);
	return;
};

// Add an item to the set.
// If an Item equal to item is already present in the set it is overwritten.
// This makes sets useful in the case where the items have a (key, value)
// structure and only the key is used for implementing Precedes() for use as a
// look up table.
func (this *Set) Add(item Item) {
	var inserted bool;
	this.root, inserted = insert(this.root, item);
	if inserted {
		this.count++;
	};
	this.root.red = false;
};

// Remove item from the set.
func (this *Set) Remove(item Item) {
	var deleted bool;
	this.root, deleted = delete(this.root, item);
	if deleted {
		this.count--;
	};
	this.root.red = false;
};

// Iterate over the set members in arbitrary type order and in order within type.
func (this *Set) Iter() <-chan Item {
	c := make(chan Item);
	go iterate(this.root, c);
	return c;
};

// Iterate asynchronously over the set members in arbitrary type order and in
// order within type. This method uses more memory than Iter() and is only
// recommended for use when circumstances preclude the use of Iter().
func (this *Set) IterAsync() <-chan Item {
	c := make(chan Item, this.count);
	iterate(this.root, c);
	return c;
};

func in_size_order(setA, setB *Set) (smallest, other *Set) {
	if setA.Cardinality() < setB.Cardinality() {
		smallest, other = setA, setB;
	} else {
		smallest, other = setB, setA;
	};
	return;
};

// Disjoint returns true if setA and setB have no members in common
func Disjoint(setA, setB *Set) bool {
	smallest, other := in_size_order(setA, setB);
	for item := range smallest.Iter() {
		if other.Has(item) {
			return false;
		};
	};
	return true;
};

// Intersect returns true if setA and setB have at least one member common
func Intersect(setA, setB *Set) bool {
	smallest, other := in_size_order(setA, setB);
	for item := range smallest.Iter() {
		if other.Has(item) {
			return true;
		};
	};
	return false;
};

// Subset returns true if every member of setA is also member of setB
//	Intersection(setA, setB) == setA
func Subset(setA, setB *Set) bool {
	if setA.Cardinality() > setB.Cardinality() { return false; };
	for item := range setA.Iter() {
		if !setB.Has(item) {
			return false;
		};
	};
	return true;
};

// ProperSubset returns true if every member of setA is also member of setB
// and they are not equal
//	Intersection(setA, setB) == setA && setA != setB
func ProperSubset(setA, setB *Set) bool {
	if setA.Cardinality() >= setB.Cardinality() { return false; };
	return Subset(setA, setB);
};

// Superset returns true if every member of setB is also member of setA
//	Intersection(setA, setB) == setB
func Superset(setA, setB *Set) bool {
	return Subset(setB, setA);
};

// Propersupeset returns true if every member of setB is also member of setA
// and they are not equal
//	Intersection(setA, setB) == setA && setA != setB
func ProperSuperset(setA, setB *Set) bool {
	return ProperSubset(setB, setA);
};

// Equal returns true if setA and setB contain exactly the same members
//	Intersection(setA, setB) == setA == setB
func Equal(setA, setB *Set) bool {
	if setA.Cardinality() != setB.Cardinality() { return false; };
	return Subset(setA, setB);
};

// Precedes() implements Item.Precedes() method for sets so that sets of sets are
// possible
func (this *Set) Precedes(other interface{}) bool {
	thisiter := this.Iter();
	otheriter := other.(*Set).Iter();
	for {
		thisitem, this_ok := <- thisiter;
		otheritem, other_ok := <- otheriter;
		if !this_ok {
			return !other_ok;
		} else if !other_ok {
			return false;
		};
		ct := cmp_type(thisitem, otheritem);
		if ct == 0 {
			if thisitem.Precedes(otheritem) {
				return true;
			} else if otheritem.Precedes(thisitem) {
				return false;
			};
		} else {
			return ct < 0;
		};
	};
	return false;
};

// Union returns a set that is the union of setA and setB
//	for any Item i:
//		(setA.Has(i) || setB.Has(i)) == Union(setA, setB).Has(i)
func Union(setA, setB *Set) (set *Set) {
	smallest, other := in_size_order(setA, setB);
	set = other.Copy();
	for item := range smallest.Iter() {
		set.Add(item);
	};
	return;
};

// Intersection returns a set that is the intersection of setA and setB
//	for any Item i:
//		(setA.Has(i) && setB.Has(i)) == Intersection(setA, setB).Has(i)
func Intersection(setA, setB *Set) (set *Set) {
	smallest, other := in_size_order(setA, setB);
	set = New();
	for item := range smallest.Iter() {
		if other.Has(item) {
			set.Add(item);
		};
	};
	return;
};

// Difference returns a set that contains the items in setA minus any items in setB
//	for any Item i:
//		(setA.Has(i) && !setB.Has(i)) == Difference(setA, setB).Has(i)
func Difference(setA, setB *Set) (set *Set) {
	set = New();
	for item := range setA.Iter() {
		if !setB.Has(item) {
			set.Add(item);
		};
	};
	return;
};

// SymmetricDifference returns a set that contains the items in setA minus or setB
// but not both
//	for any Item i:
//		((setA.Has(i) && !setB.Has(i)) || (!setA.Has(i) && setB.Has(i))) == SymmetricDifference(setA, setB).Has(i)
func SymmetricDifference(setA, setB *Set) (set *Set) {
	set = New();
	for item := range setA.Iter() {
		if !setB.Has(item) {
			set.Add(item);
		};
	};
	for item := range setB.Iter() {
		if !setA.Has(item) {
			set.Add(item);
		};
	};
	return;
};

