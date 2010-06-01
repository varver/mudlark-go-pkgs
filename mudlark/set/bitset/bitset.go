// Copyright 2010 -- Peter Williams, all rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The set package implements sets
package set

import (
	"fmt"
	"os"
)

type bitchunk uint
type bitchunkkey int64

const bitchunkSZ = (1 + ^bitchunk(0)>>32&1) * 32

// Bitset is a sparse representation of a bitset for use as a basis for
// integer sets
type bitset struct {
	// The number of bits in the set with a value of true
	bitcount uint64
	// A record of the bits in the bitset with a value of true
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
func (bset *bitset) setBit(key bitchunkkey, mask bitchunk) {
	bits := bset.bits[key] | mask
	if bits != bset.bits[key] {
		bset.bitcount++
	}
	bset.bits[key] = bits
}

// Clear the specified bit (i.e. set to false)
func (bset *bitset) clearBit(key bitchunkkey, mask bitchunk) {
	bits := bset.bits[key] & (^mask)
	if bits != bset.bits[key] {
		bset.bitcount--
	}
	bset.bits[key] = bits, bits != 0
}

// Get the value for the specified bit
func (bset *bitset) getBit(key bitchunkkey, mask bitchunk) bool {
	return (bset.bits[key] & mask) != 0
}

func makebitset() (bset *bitset) {
	bset = new(bitset)
	bset.bits = make(map[bitchunkkey]bitchunk)
	return
}

func (bset *bitset) clear() {
	bset.bitcount = 0
	bset.bits = make(map[bitchunkkey]bitchunk) // let GC clean up after us
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

func (bset *bitset) iterate(c chan<- interface{}) {
	for key, chunk := range bset.bits {
		for _, bit := range getbits(chunk) {
			c <- imemberval(key, bit)
		}
	}
	close(c)
}

func (bset *bitset) iter() <-chan interface{} {
	c := make(chan interface{})
	go bset.iterate(c)
	return c
}

func (bset bitset) tostring() string {
	str := "{"
	addcomma := false
	for member := range bset.iter() {
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
func equal(a, b *bitset) bool {
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
func subset(a, b *bitset) bool {
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

// Is set a a proper subset of set b
func propersubset(a, b *bitset) bool {
	if a.bitcount >= b.bitcount {
		return false
	}
	return subset(a, b)
}

// Is set a a superset of set b
func superset(a, b *bitset) bool {
	return subset(b, a)
}

// Is set a a proper superset of set b
func propersuperset(a, b *bitset) bool {
	return propersubset(b, a)
}

// Are a and b disjoint sets
func disjoint(a, b *bitset) bool {
	var smallest, other *bitset

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

func intersection(a, b *bitset) (bset *bitset) {
	var smallest, other *bitset

	if len(a.bits) < len(b.bits) {
		smallest = a
		other = b
	} else {
		smallest = b
		other = a
	}
	bset = makebitset()
	for key, schunk := range smallest.bits {
		chunk := schunk & other.bits[key]
		if chunk != 0 {
			bset.bits[key] = chunk
			bset.bitcount += uint64(bitcount(chunk))
		}
	}
	return
}

func bitsetcopy(a *bitset) (bset *bitset) {
	bset = makebitset()
	for akey, achunk := range a.bits {
		bset.bits[akey] = achunk
	}
	bset.bitcount = a.bitcount
	return
}

func union(a, b *bitset) (bset *bitset) {
	bset = bitsetcopy(a)
	for bkey, bchunk := range b.bits {
		bset.bits[bkey] |= bchunk
	}
	bset.bitcount = 0
	for _, chunk := range bset.bits {
		bset.bitcount += uint64(bitcount(chunk))
	}
	return
}

func difference(a, b *bitset) (bset *bitset) {
	bset = makebitset()
	for akey, achunk := range a.bits {
		var chunk bitchunk = achunk & (^b.bits[akey])
		if chunk != 0 {
			bset.bits[akey] = chunk
			bset.bitcount += uint64(bitcount(chunk))
		}
	}
	return
}

func symmetricdifference(a, b *bitset) (bset *bitset) {
	bset = union(difference(a, b), difference(b, a))
	return
}

