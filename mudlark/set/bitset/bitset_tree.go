// Copyright 2010 -- Peter Williams, all rights reserved
// Use of this source code is governed by the new BSD license.

// The set package implements integer sets
package set

import (
	"fmt"
	"os"
	"mudlark/tree/llrb_tree"
)

type bitchunk uint
type bitchunkkey int64
type bitrecord struct {
	key bitchunkkey
	chunk bitchunk
}

func (this *bitrecord) Precedes(other interface{}) bool {
	return this.key < other.(*bitrecord).key
}

const bitchunkSZ = (1 + ^bitchunk(0)>>32&1) * 32

// Bitset is a sparse representation of a bitset for use as a basis for
// integer sets
type bitset struct {
	// The number of bits in the set with a value of true
	bitcount uint64
	// A record of the bits in the bitset with a value of true
	// Bit i's value is stored in bit i % 32 of bits[i / 32]
	bits *llrb_tree.Tree
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
	record, found := bset.bits.Find(&bitrecord{key, 0})
	if found {
		chunk := record.(*bitrecord).chunk
		mask |= chunk
		if mask != chunk {
			record.(*bitrecord).chunk = mask
			bset.bitcount++
		}
	} else {
		bset.bits.Insert(&bitrecord{key, mask})
		bset.bitcount++
	}
}

// Clear the specified bit (i.e. set to false)
func (bset *bitset) clearBit(key bitchunkkey, mask bitchunk) {
	record, found := bset.bits.Find(&bitrecord{key, 0})
	if found {
		chunk := record.(*bitrecord).chunk
		newchunk := chunk & (^mask)
		if newchunk != chunk {
			if newchunk == 0 {
				bset.bits.Delete(&bitrecord{key, 0})
			} else {
				bset.bits.Insert(&bitrecord{key, newchunk})
			}
			bset.bitcount--
		}
	}
}

// Get the value for the specified bit
func (bset *bitset) getBit(key bitchunkkey, mask bitchunk) bool {
	record, found := bset.bits.Find(&bitrecord{key, 0})
	if found {
		return (record.(*bitrecord).chunk & mask) != 0
	}
	return false
}

func makebitset() (bset *bitset) {
	bset = new(bitset)
	bset.bits = llrb_tree.Make(true)
	return
}

func (bset *bitset) clear() {
	bset.bitcount = 0
	bset.bits = llrb_tree.Make(true) // let GC clean up after us
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
	for irecord := range bset.bits.Iter(llrb_tree.IN_ORDER) {
		record := irecord.(*bitrecord)
		for _, bit := range getbits(record.chunk) {
			c <- imemberval(record.key, bit)
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
	if a.bitcount != b.bitcount || a.bits.Len() != b.bits.Len() {
		return false
	} else {
		for iarecord := range a.bits.Iter(llrb_tree.IN_ORDER) {
			arecord := iarecord.(*bitrecord)
			ibrecord, found := b.bits.Find(arecord)
			if !found || arecord.chunk != ibrecord.(*bitrecord).chunk {
				return false
			}
		}
	}
	return true
}

// Is set a a subset of set b
func subset(a, b *bitset) bool {
	if a.bitcount > b.bitcount || a.bits.Len() > b.bits.Len() {
		return false
	} else {
		for iarecord := range a.bits.Iter(llrb_tree.IN_ORDER) {
			arecord := iarecord.(*bitrecord)
			ibrecord, found := b.bits.Find(arecord)
			if !found || (arecord.chunk & ibrecord.(*bitrecord).chunk) != arecord.chunk {
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

	if a.bits.Len() < b.bits.Len() {
		smallest = a
		other = b
	} else {
		smallest = b
		other = a
	}
	for ismrecord := range smallest.bits.Iter(llrb_tree.IN_ORDER) {
		smrecord := ismrecord.(*bitrecord)
		iotrecord, found := other.bits.Find(smrecord)
		if found && (smrecord.chunk & iotrecord.(*bitrecord).chunk) != 0 {
			return false
		}
	}
	return true
}

func intersection(a, b *bitset) (bset *bitset) {
	var smallest, other *bitset

	if a.bits.Len() < b.bits.Len() {
		smallest = a
		other = b
	} else {
		smallest = b
		other = a
	}
	bset = makebitset()
	for ismrecord := range smallest.bits.Iter(llrb_tree.IN_ORDER) {
		smrecord := ismrecord.(*bitrecord)
		iotrecord, found := other.bits.Find(smrecord)
		if found {
			chunk := smrecord.chunk & iotrecord.(*bitrecord).chunk
			if chunk != 0 {
				bset.bits.Insert(&bitrecord{smrecord.key, chunk})
				bset.bitcount += uint64(bitcount(chunk))
			}
		}
	}
	return
}

func bitsetcopy(a *bitset) (bset *bitset) {
	bset = makebitset()
	bset.bits = a.bits.Copy()
	bset.bitcount = a.bitcount
	return
}

func union(a, b *bitset) (bset *bitset) {
	bset = bitsetcopy(a)
	for ibrecord := range b.bits.Iter(llrb_tree.IN_ORDER) {
		brecord := ibrecord.(*bitrecord)
		irecord, found := bset.bits.Find(brecord)
		if found {
			record := irecord.(*bitrecord)
			newchunk := brecord.chunk | record.chunk
			bset.bits.Insert(&bitrecord{brecord.key, newchunk})
			bset.bitcount += uint64(bitcount(newchunk) - bitcount(record.chunk))
		} else {
			bset.bits.Insert(brecord)
			bset.bitcount += uint64(bitcount(brecord.chunk))
		}
	}
	return
}

func difference(a, b *bitset) (bset *bitset) {
	bset = makebitset()
	for iarecord := range a.bits.Iter(llrb_tree.IN_ORDER) {
		arecord := iarecord.(*bitrecord)
		ibrecord, found := b.bits.Find(arecord)
		if found {
			brecord := ibrecord.(*bitrecord)
			chunk := arecord.chunk & (^brecord.chunk)
			if chunk != 0 {
				bset.bits.Insert(&bitrecord{arecord.key, chunk})
				bset.bitcount += uint64(bitcount(chunk))
			}
		} else {
			bset.bits.Insert(arecord)
			bset.bitcount += uint64(bitcount(arecord.chunk))
		}
	}
	return
}

func symmetricdifference(a, b *bitset) (bset *bitset) {
	bset = union(difference(a, b), difference(b, a))
	return
}

