// Copyright 2010 -- Peter Williams, all rights reserved
// Use of this source code is governed by the new BSD license.

// Implement 2-3 left Leaning Red Black Binary Trees.
//
// Based on the Java implementation described by Robert Sedgewick
// in his paper entitled "Left-leaning Red-Black Trees"
// available at: <www.cs.princeton.edu/~rs/talks/LLRB/LLRB.pdf>.
// The principal difference (other than the conversion to Go) is that the items
// being inserted combine the roles of both key and value.
package llrb_tree

// Items to be inserted in a tree must implement this interface and must
// satisfy the following formal requirements (where a, b and c are all
// instances of the same type):
//	 a.Precedes(b) implies !b.Precedes(a)
//	 a.Precedes(b) && b.Precedes(c) implies a.Precedes(c)
//	 !a.Precedes(b) && !b.Precedes(a) implies a == b
type Item interface {
	Precedes(other interface{}) bool
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
	if item.Precedes(node.item) {
		node.left, inserted = insert(node.left, item)
	} else if node.item.Precedes(item) {
		node.right, inserted = insert(node.right, item)
	} else {
		node.item = item
	} 
	return fix_up(node), inserted
}

func insert_keep_duplicates(node *ll_rb_node, item Item) (*ll_rb_node) {
	if node == nil {
		return new_ll_rb_node(item)
	}
	if item.Precedes(node.item) {
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
	if item.Precedes(node.item) {
		if !is_red(node.left) && !is_red(node.left.left) {
			node = move_red_left(node)
		}
		node.left, deleted = delete(node.left, item)
	} else {
		if is_red(node.left) {
			node = rotate_right(node)
		}
		if !node.item.Precedes(item) && !item.Precedes(node.item) && node.right == nil {
			return nil, true
		}
		if !is_red(node.right) && !is_red(node.right.left) {
			node = move_red_right(node)
		}
		if !node.item.Precedes(item) && !item.Precedes(node.item) {
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

func copy(node *ll_rb_node) *ll_rb_node {
	if node == nil { return nil }
	clone := new(ll_rb_node)
	clone.item = node.item
	clone.red = node.red
	clone.left = copy(node.left)
	clone.right = copy(node.right)
	return clone
}

// Tree is a Left-leaning Red/Black binary tree of objects that satisfy the
// Item interface.  Instances of Tree must be initialized using Make()
// before use.  E.g.:
//	var t Tree = llrb_tree.Make(true)
type Tree struct {
	root *ll_rb_node
	count uint
	keep_duplicates bool
}

// Find an item in the tree.  Useful for look up tables.
func (this *Tree) Find(item Item) (entry Item, found bool) {
	if this.count == 0 {
		return
	}
	for node := this.root; node != nil && !found; {
		if item.Precedes(node.item) {
			node = node.left
		} else if node.item.Precedes(item) {
			node = node.right
		} else {
			entry = node.item
			found = true
		}
	}
	return
}

// Insert item in the tree.  If the tree was initialized to filter out
// duplicates the item being inserted will overwrite any equal item already
// in the tree.  This allows the tree to be used as a look up table using
// {key, value} item types where Precedes() ony uses the key.
func (this *Tree) Insert(item Item) {
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

// Delete item from the tree. If item has duplicates in the tree only one will
// be deleted.
func (this *Tree) Delete(item Item) {
	var deleted bool
	this.root, deleted = delete(this.root, item)
	if deleted {
		this.count--
	}
	this.root.red = false
}

// Iterate over the tree in the order specified:
//	order == IN_ORDER: in order as defined by Item.Precedes()
//	order == REVERSE_ORDER: in reverse order as defined by Item.Precedes()
//	order == PRE_ORDER: in binary tree pre order
//	order == POST_ORDER: in binary tree post order
func (this *Tree) Iter(order int) <-chan Item {
	c := make(chan Item)
	go iterate(this.root, c, order)
	return c
}

// Make a Tree. The parameter "filtered" determines whether duplicate items
// will be filtered out (or kept) during insertion.
func Make(filtered bool) (tree *Tree) {
	tree = new(Tree)
	tree.keep_duplicates = !filtered
	return
}

// Make a copy of this tree.
func (this *Tree) Copy() (tree *Tree) {
	tree = Make(!this.keep_duplicates)
	tree.root = copy(this.root)
	tree.count = this.count
	return
}

// Len returns the number of items in the tree.
func (this *Tree) Len() uint {
	return this.count
}

// Is there an instance equal to item in the tree.
func (this *Tree) Has(item Item) (found bool) {
	_, found = this.Find(item)
	return
}

