// Copyright 2010 -- Peter Williams, all rights reserved
// Use of this source code is governed by the new BSD license.

package sort_test

import (
	"testing"
	"rand"
	"mudlark/sort"
//	"fmt"
//	"reflect"
)

type Int int

func (i Int) Precedes(other interface{}) bool {
	return int(i) < int(other.(Int))
}

type IntArray []Int

func (this IntArray) iterate(c chan<- sort.Item) {
	for _, i := range this {
		c <- i
	}
	close(c)
}

func (this IntArray) iterator() <-chan sort.Item {
	c := make(chan sort.Item)
	go this.iterate(c)
	return c
}

type Real float64

func (r Real) Precedes(other interface{}) bool {
	return float64(r) < float64(other.(Real))
}

func TestMakeSortSlice(t *testing.T) {
	const sz = 1000
	ints := make([]sort.Item, sz)
	for i := 0; i < sz; i++ {
		ints[i] = Int(rand.Intn(8 * sz / 10))
	}
	count := 0
	var lasti sort.Item
	for _, i := range sort.SortSlice(ints) {
		if count != 0 && i.Precedes(lasti) {
			t.Errorf("Unexpected order: %v : %v", i, lasti)
		}
		count++
		lasti = i
	}
	if count != sz {
		t.Errorf("Expected count %v: got %v", sz, count)
	}
}

func TestMakeSortChan(t *testing.T) {
	const sz = 1000
	ints := make(IntArray, sz)
	for i := 0; i < sz; i++ {
		ints[i] = Int(rand.Intn(8 * sz / 10))
	}
	count := 0
	var lasti sort.Item
	for i := range sort.SortChan(ints.iterator()) {
		if count != 0 && i.Precedes(lasti) {
			t.Errorf("Unexpected order: %v : %v", i, lasti)
		}
		count++
		lasti = i
	}
	if count != sz {
		t.Errorf("Expected count %v: got %v", sz, count)
	}
}

func BenchmarkSortSlice(b *testing.B) {
	const sz = 1000
	b.SetBytes(sz)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ints := make([]sort.Item, sz)
		for i := 0; i < sz; i++ {
			ints[i] = Int(rand.Intn(8 * sz / 10))
		}
		b.StartTimer()
		ints = sort.SortSlice(ints)
	}
}

