// Copyright 2010 -- Peter Williams, all rights reserved
// Use of this source code is governed by the new BSD license.

// Implement sort operations.
package sort

import "mudlark/tree/llrb_tree"

// Items to be sorted must implement this interface and must satisfy the
// following formal requirements (where a, b and c are all instances of the
// same type):
//	 a.Less(b) implies !b.Less(a)
//	 a.Less(b) && b.Less(c) implies a.Less(c)
//	 !a.Less(b) && !b.Less(a) implies a == b
// This method will only be used when reflect.Typeof() the calling object
// matches reflect.Typeof() of other.
type Item interface {
	Less(other interface{}) bool
}

func slice_to_tree(slice []Item, filtered bool) (tree llrb_tree.Tree) {
	tree = llrb_tree.Make(filtered)
	for _, item := range slice {
		tree.Insert(item)
	}
	return
}

func tree_to_slice(tree llrb_tree.Tree, order int) (slice []Item) {
	slice = make([]Item, tree.Len())
	var i int
	for item := range tree.Iter(order) {
		slice[i] = item
		i++
	}
	return
}

// SortSlice() returns a copy of a slice in order as defined by Item.Less().
func SortSlice(slice []Item) (sorted []Item) {
	tree := slice_to_tree(slice, false)
	return tree_to_slice(tree, llrb_tree.IN_ORDER)
}

// SortFilteredSlice() returns a copy of a slice in order as defined by
// Item.Less() filtering out duplicate items.
func SortFilteredSlice(slice []Item) (sorted []Item) {
	tree := slice_to_tree(slice, true)
	return tree_to_slice(tree, llrb_tree.IN_ORDER)
}

// ReverseSortSlice() returns a copy of a slice in reverse order as defined
// by Item.Less().
func ReverseSortSlice(slice []Item) (sorted []Item) {
	tree := slice_to_tree(slice, false)
	return tree_to_slice(tree, llrb_tree.REVERSE_ORDER)
}

// ReverseSortFilteredSlice() returns a copy of a slice in reverse order as
// defined by Item.Less() filtering out duplicate items.
func ReverseSortFilteredSlice(slice []Item) (sorted []Item) {
	tree := slice_to_tree(slice, true)
	return tree_to_slice(tree, llrb_tree.REVERSE_ORDER)
}

// Now do the same thing for channels (for use with iterators)

func chan_to_tree(channel <-chan Item, filtered bool) (tree llrb_tree.Tree) {
	tree = llrb_tree.Make(filtered)
	for item := range channel {
		tree.Insert(item)
	}
	return
}

func tree_to_chan(tree llrb_tree.Tree, order int) (channel chan Item) {
	channel = make(chan Item, tree.Len())
	for item := range tree.Iter(order) {
		channel <- item
	}
	close(channel)
	return channel
}

// SortChan() returns a new <-chan which will emit contents of channel
// in order as defined by Item.Less().
func SortChan(channel <-chan Item) (<-chan Item) {
	tree := chan_to_tree(channel, false)
	return tree_to_chan(tree, llrb_tree.IN_ORDER)
}

// SortFilteredChan() returns a copy of a chan in order as defined by
// Item.Less() filtering out duplicate items.
func SortFilteredChan(channel <-chan Item) (<-chan Item) {
	tree := chan_to_tree(channel, true)
	return tree_to_chan(tree, llrb_tree.IN_ORDER)
}

// ReverseSortChan() returns a copy of a chan in reverse order as defined
// by Item.Less().
func ReverseSortChan(channel <-chan Item) (<-chan Item) {
	tree := chan_to_tree(channel, false)
	return tree_to_chan(tree, llrb_tree.REVERSE_ORDER)
}

// ReverseSortFilteredChan() returns a copy of a chan in reverse order as
// defined by Item.Less() filtering out duplicate items.
func ReverseSortFilteredChan(channel <-chan Item) (<-chan Item) {
	tree := chan_to_tree(channel, true)
	return tree_to_chan(tree, llrb_tree.REVERSE_ORDER)
}

