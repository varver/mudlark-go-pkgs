// Copyright 2010 -- Peter Williams, all rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The heteroset package implements heterogeneou sets
package heteroset

import "reflect"
import "fmt"

// Implement 2-3 left Leaning Red Black Trees for for internal representation.
// It is based on the Java implementation described by Robert Sedgewick
// in his paper entitled "left-leaning Red-Black Trees"
// available at: <www.cs.princeton.edu/~rs/talks/LLRB/LLRB.pdf>.
// The principal difference (other than the conversion to Go) is that the items
// being inserted combine the roles of both key and value

// Prospective set items must implement this interface and must satisfy the
// following formal requirements:
//	 a.Compare(b) < 0 implies b.Compare(a) > 0
//	 a.Compare(b) > 0 implies b.Compare(a) < 0
//	 a.Compare(b) == 0 implies b.Compare(a) == 0
// This method will only be called if the type of the calling object matches
// that of other.
type Item interface {
	Compare(other Item) int
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
	return this.item.Compare(item)
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
	case cmp < 0:
		node.left, inserted = insert(node.left, item)
	case cmp > 0:
		node.right, inserted = insert(node.right, item)
	default:
	}
	return fix_up(node), inserted
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
	if node.compare_item(item) < 0 {
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

func iterate_preorder(node *ll_rb_node, c chan<- Item) {
	if node == nil {
		return
	}
	c <- node.item
	iterate_preorder(node.left, c)
	iterate_preorder(node.right, c)
}

func iterate(node *ll_rb_node, c chan<- Item) {
	iterate_preorder(node, c)
	close(c)
}

func print_node(node *ll_rb_node) {
	if node == nil { return }
	fmt.Printf("%v\n", node)
	print_node(node.left)
	print_node(node.right)
}

type ll_rb_tree struct {
	root *ll_rb_node
	count uint64
}

func (this ll_rb_tree) find(item Item) (found bool, iterations uint) {
	if this.count == 0 {
		return
	}
	for node := this.root; node != nil && !found; {
		iterations++
		switch cmp := node.compare_item(item); {
		case cmp < 0:
			node = node.left
		case cmp > 0:
			node = node.right
		default:
			found = true
		}
	}
	return
}

func (this *ll_rb_tree) insert(item Item) {
	var inserted bool
	this.root, inserted = insert(this.root, item)
	if inserted {
		this.count++
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

func (this ll_rb_tree) iterator() <-chan Item {
	c := make(chan Item)
	go iterate(this.root, c)
	return c
}

