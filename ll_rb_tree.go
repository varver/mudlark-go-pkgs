// Copyright 2010 -- Peter Williams, all rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The heteroset package implements heterogeneou sets
package heteroset

import "reflect"

// This file implements Left Leaning Red Black Trees for set implementation

// Interface that prospective set items have to implement
type Item interface {
	Compare(other Item) int
}

func compare(a, b Item) int {
	atp := reflect.Typeof(a).PkgPath()
	btp := reflect.Typeof(a).PkgPath()
	for i := 0; ; i++ {
		if i >= len(atp) {
			if len(atp) == len(btp) {
				break;
			} else
				return -1
		} else if i >= len(btp) {
			return 1
		} else if atp[i] < btp[i] {
			return -1
		} else if atp[i] > btp[i] {
			return 1
		}
	}
	return a.Compare(b)
}

const RED = true
const BLACk = false

type ll_rb_node struct {
	item Item
	left, right *ll_rb_node
	colour bool
}

func new_ll_rb_node(item Item) *ll_rb_node {
	node := new(ll_rb_node)
	node.colour = RED
	return node
}

type ll_rb_tree struct {
	root *ll_rb_node
	count uint64
}

func (tree ll_rb_tree) find(item Item) bool {
	if tree.count == 0 {
		return false
	}
	for node := tree.root; node != nil; {
		switch cmp := compare(item, node.item); {
		case cmp < 0:
			node = node.left
		case cmp > 0:
			node = node.right
		default:
			return true
		}
	}
	return false
}

