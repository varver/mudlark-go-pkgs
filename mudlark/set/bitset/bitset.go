// Copyright 2010 -- Peter Williams, all rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The bitset package implements sets of integer numbers
package bitset

import (
	"fmt"
	"os"
)

type bitchunk uint
type bitchunkkey int64

const bitchunkSZ = (1 + ^bitchunk(0)>>32&1) * 32

// Set is a representation of integer number sets
type Set struct {
	// The number of bits in the set with a value of true
	bitcount uint64
	// A record of the bits in the Set with a value of true
	// Bit i's value is stored in bit i % 32 of bits[i / 32]
	bits map[bitchunkkey]bitchunk
}

// Location of bit representing an unsigned integer value
func ubitlocation(bit uint64) (key bitchunkkey, mask bitchunk) {
	key = bitchunkkey(bit / uint64(bitchunkSZ))
	mask = 1 << (bit % uint64(bitchunkSZ))
	return
}

// Location of bit representing a signed integer value
func sbitlocation(bit int64) (key bitchunkkey, mask bitchunk) {
	key = bitchunkkey(bit / int64(bitchunkSZ))
	if bit < 0 {
		// This is necessary because (-3 / 32) == (3 /32) etc.
		key--
		mask = 1 << uint(-bit%int64(bitchunkSZ))
	} else {
		mask = 1 << uint(bit%int64(bitchunkSZ))
	}
	return
}

// Location of bit representing arbitrary integer value
func ibitlocation(member interface{}) (key bitchunkkey, chunk bitchunk) {
	switch t := member.(type) {
	case uint:
		key, chunk = ubitlocation(uint64(member.(uint)))
	case uint8:
		key, chunk = ubitlocation(uint64(member.(uint8)))
	case uint16:
		key, chunk = ubitlocation(uint64(member.(uint16)))
	case uint32:
		key, chunk = ubitlocation(uint64(member.(uint32)))
	case uint64:
		key, chunk = ubitlocation(member.(uint64))
	case int:
		key, chunk = sbitlocation(int64(member.(int)))
	case int8:
		key, chunk = sbitlocation(int64(member.(int8)))
	case int16:
		key, chunk = sbitlocation(int64(member.(int16)))
	case int32:
		key, chunk = sbitlocation(int64(member.(int32)))
	case int64:
		key, chunk = sbitlocation(member.(int64))
	default:
		// Run time check better than no check (not as good as compile time)
		panic(os.EINVAL)
	}
	return
}

// Get the value of the member at a specific location
func imemberval(key bitchunkkey, bitn uint8) interface{} {
	if key < 0 {
		return int64(key+1)*int64(bitchunkSZ) - int64(bitn)
	}
	return uint64(key)*uint64(bitchunkSZ) + uint64(bitn)
}

// Set the specified bit to true
func (this *Set) Add(member interface{}) {
	key, mask := ibitlocation(member)
	bits := this.bits[key] | mask
	if bits != this.bits[key] {
		this.bitcount++
	}
	this.bits[key] = bits
}

// Clear the specified bit (i.e. set to false)
func (this *Set) Remove(member interface{}) {
	key, mask := ibitlocation(member)
	bits := this.bits[key] & (^mask)
	if bits != this.bits[key] {
		this.bitcount--
	}
	this.bits[key] = bits, bits != 0
}

// Get the value for the specified bit
func (this *Set) Has(member interface{}) bool {
	key, mask := ibitlocation(member)
	return (this.bits[key] & mask) != 0
}

func Make() (this *Set) {
	this = new(Set)
	this.bits = make(map[bitchunkkey]bitchunk)
	return
}

// Cardinality returns the number of items in the set.
func (this *Set) Cardinality() uint64 {
	return this.bitcount
}

func (this *Set) Clear() {
	this.bitcount = 0
	this.bits = make(map[bitchunkkey]bitchunk) // let GC clean up after us
	return
}

func bitcount(chunk bitchunk) (count uint8) {
	for temp := chunk; temp != 0; temp >>= 1 {
		if (temp & 1) != 0 {
			count++
		}
	}
	return
}

func getbits(chunk bitchunk) (bits []uint8) {
	bits = make([]uint8, bitcount(chunk))
	for bit, index := uint8(0), 0; chunk != 0; chunk >>= 1 {
		if chunk&1 == 1 {
			bits[index] = bit
			index++
		}
		bit++
	}
	return bits
}

func (this *Set) iterate(c chan<- interface{}) {
	for key, chunk := range this.bits {
		for _, bit := range getbits(chunk) {
			c <- imemberval(key, bit)
		}
	}
	close(c)
}

func (this *Set) Iter() <-chan interface{} {
	c := make(chan interface{})
	go this.iterate(c)
	return c
}

func (this *Set) String() string {
	str := "{"
	addcomma := false
	for member := range this.Iter() {
		if addcomma {
			str += fmt.Sprintf(", %v", member)
		} else {
			str += fmt.Sprintf("%v", member)
			addcomma = true
		}
	}
	str += "}"
	return str
}

// Are sets a and b equal
func Equal(a, b *Set) bool {
	if a.bitcount != b.bitcount || len(a.bits) != len(b.bits) {
		return false
	} else {
		for akey, achunk := range a.bits {
			if achunk != b.bits[akey] {
				return false
			}
		}
	}
	return true
}

// Is set a a subset of set b
func Subset(a, b *Set) bool {
	if a.bitcount > b.bitcount || len(a.bits) > len(b.bits) {
		return false
	} else {
		for akey, achunk := range a.bits {
			if (achunk & b.bits[akey]) != achunk {
				return false
			}
		}
	}
	return true
}

// Is set a a proper Subset of set b
func ProperSubset(a, b *Set) bool {
	if a.bitcount >= b.bitcount {
		return false
	}
	return Subset(a, b)
}

// Is set a a superset of set b
func Superset(a, b *Set) bool {
	return Subset(b, a)
}

// Is set a a proper superset of set b
func ProperSuperset(a, b *Set) bool {
	return ProperSubset(b, a)
}

// Are a and b disjoint sets
func Disjoint(a, b *Set) bool {
	var smallest, other *Set

	if len(a.bits) < len(b.bits) {
		smallest = a
		other = b
	} else {
		smallest = b
		other = a
	}
	for key, schunk := range smallest.bits {
		if schunk & other.bits[key] != 0 {
			return false
		}
	}
	return true
}

// Do the sets a and b intersect
func Intersect(a, b *Set) bool {
	var smallest, other *Set

	if len(a.bits) < len(b.bits) {
		smallest = a
		other = b
	} else {
		smallest = b
		other = a
	}
	for key, schunk := range smallest.bits {
		if schunk & other.bits[key] != 0 {
			return true
		}
	}
	return false
}

func Intersection(a, b *Set) (bset *Set) {
	var smallest, other *Set

	if len(a.bits) < len(b.bits) {
		smallest = a
		other = b
	} else {
		smallest = b
		other = a
	}
	bset = Make()
	for key, schunk := range smallest.bits {
		chunk := schunk & other.bits[key]
		if chunk != 0 {
			bset.bits[key] = chunk
			bset.bitcount += uint64(bitcount(chunk))
		}
	}
	return
}

func (this *Set) Copy() (bset *Set) {
	bset = Make()
	for akey, achunk := range this.bits {
		bset.bits[akey] = achunk
	}
	bset.bitcount = this.bitcount
	return
}

func Union(a, b *Set) (bset *Set) {
	bset = a.Copy()
	for bkey, bchunk := range b.bits {
		bset.bits[bkey] |= bchunk
	}
	bset.bitcount = 0
	for _, chunk := range bset.bits {
		bset.bitcount += uint64(bitcount(chunk))
	}
	return
}

func Difference(a, b *Set) (bset *Set) {
	bset = Make()
	for akey, achunk := range a.bits {
		var chunk bitchunk = achunk & (^b.bits[akey])
		if chunk != 0 {
			bset.bits[akey] = chunk
			bset.bitcount += uint64(bitcount(chunk))
		}
	}
	return
}

func SymmetricDifference(a, b *Set) (bset *Set) {
	bset = Union(Difference(a, b), Difference(b, a))
	return
}

