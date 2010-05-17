// Copyright 2010 -- Peter Williams, all rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Implement 2-3 left Leaning Red Black Binary Trees.
//
// Based on the Java implementation described by Robert Sedgewick
// in his paper entitled "Left-leaning Red-Black Trees"
// available at: <www.cs.princeton.edu/~rs/talks/LLRB/LLRB.pdf>.
// The principal difference (other than the conversion to Go) is that the items
// being inserted combine the roles of both key and value.
package llrb_tree

import "reflect"
import "fmt"


// Items to be inserted in a tree must implement this interface and must
// satisfy the following formal requirements (where a, b and c are all
// instances of the same type):
//	 a.Less(b) implies !b.Less(a)
//	 a.Less(b) && b.Less(c) implies a.Less(c)
//	 !a.Less(b) && !b.Less(a) implies a == b
// This method will only be used when reflect.Typeof() the calling object
// matches reflect.Typeof() of other.
type Item interface {
	Less(other interface{}) bool
}

// LLRB tree node
type ll_rb_node struct {
	item Item
	left, right *ll_rb_node
	red bool
}

func new_ll_rb_node(item Item) *ll_rb_node {
	node := new(ll_rb_node)
	node.item = item
	node.red = true
	return node
}

func min(a, b int) int { if a < b { return a }; return b }

func cmp_string(a, b string) int {
	for i, lim := 0, min(len(a), len(b)); i < lim; i++ {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}
	return len(a) - len(b)
}

func cmp_type(a, b interface{}) int {
	ta := reflect.Typeof(a)
	tb := reflect.Typeof(b)
	if ta == tb {
		return 0
	}
	if cp := cmp_string(ta.PkgPath(), tb.PkgPath()); cp != 0 {
		return cp
	}
	return cmp_string(ta.Name(), tb.Name())
}

func (this *ll_rb_node) compare_item(item Item) int {
	if ct := cmp_type(this.item, item); ct != 0 {
		return ct
	}
	if this.item.Less(item) {
		return -1
	} else if item.Less(this.item) {
		return 1
	}
	return 0
}

func is_red(node *ll_rb_node) bool { return node != nil && node.red }

func flip_colours(node *ll_rb_node) {
	node.red = !node.red
	node.left.red = !node.left.red
	node.right.red = !node.right.red
}

func rotate_left(node *ll_rb_node) *ll_rb_node {
	tmp := node.right
	node.right = tmp.left
	tmp.left = node
	tmp.red = node.red
	node.red = true
	return tmp
}

func rotate_right(node *ll_rb_node) *ll_rb_node {
	tmp := node.left
	node.left = tmp.right
	tmp.right = node
	tmp.red = node.red
	node.red = true
	return tmp
}

func fix_up(node *ll_rb_node) *ll_rb_node {
	if is_red(node.right) && !is_red(node.left) {
		node = rotate_left(node)
	}
	if is_red(node.left) && is_red(node.left.left) {
		node = rotate_right(node)
	}
	if is_red(node.left) && is_red(node.right) {
		flip_colours(node)
	}
	return node
}

func insert(node *ll_rb_node, item Item) (*ll_rb_node, bool) {
	if node == nil {
		return new_ll_rb_node(item), true
	}
	inserted := false
	switch cmp := node.compare_item(item); {
	case cmp > 0:
		node.left, inserted = insert(node.left, item)
	case cmp < 0:
		node.right, inserted = insert(node.right, item)
	default:
	}
	return fix_up(node), inserted
}

func insert_keep_duplicates(node *ll_rb_node, item Item) (*ll_rb_node) {
	if node == nil {
		return new_ll_rb_node(item)
	}
	if node.compare_item(item) > 0 {
		node.left = insert_keep_duplicates(node.left, item)
	} else {
		node.right = insert_keep_duplicates(node.right, item)
	}
	return fix_up(node)
}

func move_red_left(node *ll_rb_node) *ll_rb_node {
	flip_colours(node)
	if (is_red(node.right.left)) {
		node.right = rotate_right(node.right)
		node = rotate_left(node)
		flip_colours(node)
	}
	return node
}

func move_red_right(node *ll_rb_node) *ll_rb_node {
	flip_colours(node)
	if (is_red(node.left.left)) {
		node = rotate_right(node)
		flip_colours(node)
	}
	return node
}

func delete_left_most(node *ll_rb_node) *ll_rb_node {
	if node.left == nil {
		return nil
	}
	if !is_red(node.left) && !is_red(node.left.left) {
		node = move_red_left(node)
	}
	node.left = delete_left_most(node.left)
	return fix_up(node)
}

func delete(node *ll_rb_node, item Item) (*ll_rb_node, bool) {
	var deleted bool
	if node.compare_item(item) > 0 {
		if !is_red(node.left) && !is_red(node.left.left) {
			node = move_red_left(node)
		}
		node.left, deleted = delete(node.left, item)
	} else {
		if is_red(node.left) {
			node = rotate_right(node)
		}
		if node.compare_item(item) == 0 && node.right == nil {
			return nil, true
		}
		if !is_red(node.right) && !is_red(node.right.left) {
			node = move_red_right(node)
		}
		if node.compare_item(item) == 0 {
			left_most := node.right
			for left_most.left != nil {
				left_most = left_most.left
			}
			node.item = left_most.item
			node.right = delete_left_most(node.right)
			deleted = true
		} else {
			node.right, deleted = delete(node.right, item)
		}
	}
	return fix_up(node), deleted
}

// Iteration using recursion is safe because the depth of the tree should never
// be greater than 2Log2(N) where N is the number of nodes in the tree and
// (in general) will be approximately Log2(N).

func iterate_preorder(node *ll_rb_node, c chan<- Item) {
	if node == nil {
		return
	}
	c <- node.item
	iterate_preorder(node.left, c)
	iterate_preorder(node.right, c)
}

func iterate_inorder(node *ll_rb_node, c chan<- Item) {
	if node == nil {
		return
	}
	iterate_inorder(node.left, c)
	c <- node.item
	iterate_inorder(node.right, c)
}

func iterate_postorder(node *ll_rb_node, c chan<- Item) {
	if node == nil {
		return
	}
	iterate_postorder(node.left, c)
	iterate_postorder(node.right, c)
	c <- node.item
}

func iterate_reverseorder(node *ll_rb_node, c chan<- Item) {
	if node == nil {
		return
	}
	iterate_reverseorder(node.right, c)
	c <- node.item
	iterate_reverseorder(node.left, c)
}

// Specify output order for iteration.
const (
	PRE_ORDER = iota
	IN_ORDER
	POST_ORDER
	REVERSE_ORDER
)

func iterate(node *ll_rb_node, c chan<- Item, order int) {
	switch order {
	case PRE_ORDER:
		iterate_preorder(node, c)
	case IN_ORDER:
		iterate_inorder(node, c)
	case POST_ORDER:
		iterate_postorder(node, c)
	case REVERSE_ORDER:
		iterate_reverseorder(node, c)
	}
	close(c)
}

func print_node(node *ll_rb_node) {
	if node == nil { return }
	fmt.Printf("%v\n", node)
	print_node(node.left)
	print_node(node.right)
}

func copy(node *ll_rb_node) *ll_rb_node {
	if node == nil { return nil }
	clone := new(ll_rb_node)
	clone.item = node.item
	clone.red = node.red
	clone.left = copy(node.left)
	clone.right = copy(node.right)
	return clone
}

type ll_rb_tree struct {
	root *ll_rb_node
	count uint
	keep_duplicates bool
}

func (this ll_rb_tree) find(item Item) (found bool, iterations uint) {
	if this.count == 0 {
		return
	}
	for node := this.root; node != nil && !found; {
		iterations++
		switch cmp := node.compare_item(item); {
		case cmp > 0:
			node = node.left
		case cmp < 0:
			node = node.right
		default:
			found = true
		}
	}
	return
}

func (this *ll_rb_tree) insert(item Item) {
	if this.keep_duplicates {
		this.root = insert_keep_duplicates(this.root, item)
		this.count++
	} else {
		var inserted bool
		this.root, inserted = insert(this.root, item)
		if inserted {
			this.count++
		}
	}
	this.root.red = false
}

func (this *ll_rb_tree) delete(item Item) {
	var deleted bool
	this.root, deleted = delete(this.root, item)
	if deleted {
		this.count--
	}
	this.root.red = false
}

func (this ll_rb_tree) iterator(order int) <-chan Item {
	c := make(chan Item)
	go iterate(this.root, c, order)
	return c
}

// Tree is a Left-leaning Red/Black binary tree of objects that satisfy the
// Item interface.  Instances of Tree must be initialized using Make()
// before use.  E.g.:
//	var t Tree = llrb_tree.Make(true)
type Tree struct {
	Data *ll_rb_tree
}

// Make a Tree. The parameter "filtered" determines whether duplicate items
// will be filtered out (or kept) during insertion.
func Make(filtered bool) (tree Tree) {
	tree.Data = new(ll_rb_tree)
	tree.Data.keep_duplicates = !filtered
	return
}

// Make a copy of this tree.
func (this Tree) Copy() (tree Tree) {
	tree = Make(!this.Data.keep_duplicates)
	tree.Data.root = copy(this.Data.root)
	tree.Data.count = this.Data.count
	return
}

// Len returns the number of items in the tree.
func (this Tree) Len() uint {
	return this.Data.count
}

// Insert item in the tree.
func (this *Tree) Insert(item Item) {
	this.Data.insert(item)
}

// Delete item from the tree. If item has duplicates in the tree only one will
// be deleted.
func (this *Tree) Delete(item Item) {
	this.Data.delete(item)
}

// Is there an instance equal to item in the tree.
func (this *Tree) Has(item Item) (found bool) {
	found, _ = this.Data.find(item)
	return
}

// Iterate over the tree in the order specified:
//	order == IN_ORDER: in order as defined by Item.Less()
//	order == REVERSE_ORDER: in reverse order as defined by Item.Less()
//	order == PRE_ORDER: in binary tree pre order
//	order == POST_ORDER: in binary tree post order
func (this Tree) Iter(order int) <-chan Item {
	return this.Data.iterator(order)
}


