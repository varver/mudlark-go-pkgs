// Copyright 2010 -- Peter Williams, all rights reserved
// Use of this source code is governed by the new BSD license.

package bitset

import (
	"testing"
	"rand"
	"reflect"
	//"mudlark/tree/llrb_tree"
)

func TestMakeSet(t *testing.T) {
	set := Make()
	if reflect.Typeof(set).String() != "*bitset.Set" {
		t.Errorf("Expected type \"*bitset.Set\": got %v", reflect.Typeof(set).String())
	}
	if set.bitcount != 0 {
		t.Errorf("Expected bitcount 0: got %v", set.bitcount)
	}
	if set.bits == nil {
		t.Errorf("Bit map unitialized")
	}
	//if set.bits.Len() != 0 {
	if len(set.bits) != 0 {
		//t.Errorf("Expected len(bits) 0: got %v", set.bits.Len())
		t.Errorf("Expected len(bits) 0: got %v", len(set.bits))
	}
}

func TestMakeSetWithMembers(t *testing.T) {
	set := Make(-1, 28, 18, 28, 9)
	if reflect.Typeof(set).String() != "*bitset.Set" {
		t.Errorf("Expected type \"*bitset.Set\": got %v", reflect.Typeof(set).String())
	}
	if set.bitcount != 4 {
		t.Errorf("Expected bitcount 4: got %v", set.bitcount)
	}
	if set.bits == nil {
		t.Errorf("Bit map unitialized")
	}
}

func TestKeyMappingInt64(t *testing.T) {
	for i := 0; i < 10000; i++ {
		num := rand.Int63()
		if i % 5 != 0 {
			num = -num
		}
		key, mask := sbitlocation(num)
		if bitcount(mask) != 1 {
			t.Errorf("Expected exactly one bit in mask; found %v", bitcount(mask))
		}
		dcnum := imemberval(key, getbits(mask)[0])
		switch tp := dcnum.(type) {
		case int64:
			if num >= 0 {
				t.Errorf("Expected type \"uint64\": got %v", reflect.Typeof(tp))
			}
			if num != tp {
				t.Errorf("Expected type %v: got %v (%v,%v)", num, tp, key, mask)
			}
		case uint64:
			if num < 0 {
				t.Errorf("Expected type \"int64\": got %v", reflect.Typeof(tp))
			}
			if uint64(num) != tp {
				t.Errorf("Expected type %v: got %v", num, tp)
			}
		default:
			t.Errorf("Expected type \"(u)int64\": got %v", reflect.Typeof(tp))
		}
	}
}

func TestKeyMappingUint64(t *testing.T) {
	for i := 0; i < 10000; i++ {
		num := uint64(rand.Int63())
		key, mask := ubitlocation(num)
		if bitcount(mask) != 1 {
			t.Errorf("Expected exactly one bit in mask; found %v", bitcount(mask))
		}
		dcnum := imemberval(key, getbits(mask)[0])
		switch tp := dcnum.(type) {
		case uint64:
			if num < 0 {
				t.Errorf("Expected type \"int64\": got %v", reflect.Typeof(tp))
			}
			if uint64(num) != tp {
				t.Errorf("Expected type %v: got %v", num, tp)
			}
		default:
			t.Errorf("Expected type \"(u)int64\": got %v", reflect.Typeof(tp))
		}
	}
}

func checkkmn(key bitchunkkey, mask bitchunk, num int64, t *testing.T) {
		if bitcount(mask) != 1 {
			t.Errorf("Expected exactly one bit in mask; found %v", bitcount(mask))
		}
		dcnum := imemberval(key, getbits(mask)[0])
		switch tp := dcnum.(type) {
		case int64:
			if num >= 0 {
				t.Errorf("Expected type \"uint64\": got %v", reflect.Typeof(tp))
			}
			if num != tp {
				t.Errorf("Expected num %v: got %v (%v,%v)", num, tp, key, mask)
			}
		case uint64:
			if num < 0 {
				t.Errorf("Expected type \"int64\": got %v", reflect.Typeof(tp))
			}
			if uint64(num) != tp {
				t.Errorf("Expected num %v: got %v", num, tp)
			}
		default:
			t.Errorf("Expected type \"(u)int64\": got %v", reflect.Typeof(tp))
		}
}

func TestKeyMappingInterface(t *testing.T) {
	const uintsz = (1 + ^uint(0)>>32&1) * 32
	var key bitchunkkey
	var mask bitchunk
	for i := 0; i < 1000; i++ {
		num := rand.Int63n(1<<8)
		key, mask = ibitlocation(uint8(num))
		checkkmn(key, mask, num, t)
		num = rand.Int63n(1<<16)
		key, mask = ibitlocation(uint16(num))
		checkkmn(key, mask, num, t)
		num = rand.Int63n(1<<32)
		key, mask = ibitlocation(uint32(num))
		checkkmn(key, mask, num, t)
		num = rand.Int63n(1<<63-1)
		key, mask = ibitlocation(uint64(num))
		checkkmn(key, mask, num, t)
		num = rand.Int63n(1<<8) - (1 << 7)
		key, mask = ibitlocation(int8(num))
		checkkmn(key, mask, num, t)
		num = rand.Int63n(1<<16) - (1 << 15)
		key, mask = ibitlocation(int16(num))
		checkkmn(key, mask, num, t)
		num = rand.Int63n(1<<32) - (1 << 31)
		key, mask = ibitlocation(int32(num))
		checkkmn(key, mask, num, t)
		num = rand.Int63n(1<<63-1) - (1 << 62)
		key, mask = ibitlocation(int64(num))
		checkkmn(key, mask, num, t)
		num = rand.Int63n(1<<uintsz)
		key, mask = ibitlocation(uint(num))
		checkkmn(key, mask, num, t)
		num = rand.Int63n(1<<uintsz) - (1 << (uintsz - 1))
		key, mask = ibitlocation(int(num))
		checkkmn(key, mask, num, t)
	}
}

func checkbitcount(bset *Set, str string, t *testing.T) {
	var count uint64 = 0
	//for record := range bset.bits.Iter(llrb_tree.IN_ORDER) {
	//	count += uint64(bitcount(record.(*bitrecord).chunk))
	for _, chunk := range bset.bits {
		count += uint64(bitcount(chunk))
	}
	if count != bset.bitcount {
		t.Errorf("Bit count %s. Expected: %v got: %v", str, bset.bitcount, count)
	}
}

func TestKeyBitcountAddAndRemove(t *testing.T) {
	const loopsz = 1000
	bset := Make()
	for i := 0; i < loopsz; i++ {
		bset.Add(i)
		checkbitcount(bset, "add(sequence)", t)
	}
	for i := 0; i < loopsz; i++ {
		bset.Add(rand.Int63())
		checkbitcount(bset, "add(random(spread))", t)
	}
	for i := 0; i < loopsz; i++ {
		bset.Add(rand.Int63n(loopsz * 2))
		checkbitcount(bset, "add(random(focussed))", t)
	}
	for i := 0; i < loopsz; i++ {
		bset.Remove(rand.Int63())
		checkbitcount(bset, "remove(random(spread))", t)
	}
	for i := 0; i < loopsz; i++ {
		bset.Remove(rand.Int63n(loopsz * 2))
		checkbitcount(bset, "remove(random(focussed))", t)
	}
	for i := 0; i < loopsz; i++ {
		bset.Remove(i)
		checkbitcount(bset, "remove(sequence)", t)
	}
}

func BenchmarkMakeEmptySet(b *testing.B) {
	b.SetBytes(1)
	for i := 0; i < b.N; i++ {
		set := Make()
		b.StopTimer()
		if set.bitcount > 0 {
			// This is just here to stop the compiler complaining
		}
		b.StartTimer()
	}
}

func BenchmarkInsertRandom(b *testing.B) {
	const N = 50000
	b.SetBytes(N)
	for ib := 0; ib < b.N; ib++ {
		b.StopTimer()
		set := Make()
		b.StartTimer()
		for i := 0; i < N; i++ {
			set.Add(rand.Int())
		}
	}
}

func BenchmarkInsertSerial(b *testing.B) {
	const N = 50000
	b.SetBytes(N)
	for ib := 0; ib < b.N; ib++ {
		b.StopTimer()
		set := Make()
		b.StartTimer()
		for i := 0; i < N; i++ {
			set.Add(i)
		}
	}
}

func TestIterate(t *testing.T) {
	set := Make()
	var count int
	for i := 0; i < 10000; i++ {
		set.Add(rand.Int())
		count++
	}
	for _ = range set.Iter() {
		count--
	}
	if count != 0 {
		t.Errorf("%v count", count)
	}
}

func make_set_serial(begin, end int64) (set *Set) {
	set = Make()
	for i := begin; i <= end; i++ {
		set.Add(i)
	}
	return
}

func TestDisjointIntersect(t *testing.T) {
	setA := make_set_serial(-100, 0)
	setB := make_set_serial(1, 100)
	setC := make_set_serial(-50, 50)
	if Intersect(setA, setB) {
		t.Errorf("setA and setB should be disjoint")
	}
	if !Intersect(setA, setC) {
		t.Errorf("setA and setC should intersect")
	}
	if Intersect(setA, setB) && Disjoint(setA, setB) {
		t.Errorf("Intersect(A, B) and Disjoint(A, B) should be mutually exclusive")
	}
	if Intersect(setA, setC) && Disjoint(setA, setC) {
		t.Errorf("Intersect(A, B) and Disjoint(A, B) should be mutually exclusive")
	}
	if Intersect(setA, setB) != Intersect(setB, setA) {
		t.Errorf("Intersect(A, B) and Intersect(B, A) should be equal")
	}
	if Intersect(setA, setC) != Intersect(setC, setA) {
		t.Errorf("Intersect(A, C) and Intersect(C, A) should be equal")
	}
	if Disjoint(setA, setB) != Disjoint(setB, setA) {
		t.Errorf("Disjoint(A, B) and Disjoint(B, A) should be equal")
	}
	if Disjoint(setA, setC) != Disjoint(setC, setA) {
		t.Errorf("Disjoint(A, C) and Disjoint(C, A) should be equal")
	}
}

func TestUnion(t *testing.T) {
	setA := make_set_serial(-100, 0)
	setB := make_set_serial(1, 100)
	setC := make_set_serial(-50, 50)
	setAuB := Union(setA, setB)
	setAuC := Union(setA, setC)
	if !Intersect(setA, setAuB) || !Intersect(setB, setAuB) {
		t.Errorf("setAuB should intersect with both setA and SetB")
	}
	if setAuB.Cardinality() != setA.Cardinality() + setB.Cardinality() {
		t.Errorf("Cardinality of a union of disjoint sets should be the sum of their cardinalities")
	}
	for item := range setAuB.Iter() {
		if !setA.Has(item) && !setB.Has(item) {
			t.Errorf("Items in setAuB should be in either setA or setB")
		}
	}
	if !Intersect(setA, setAuC) || !Intersect(setC, setAuC) {
		t.Errorf("setAuC should intersect with both setA and SetC")
	}
	var incommon uint64
	for item := range setAuC.Iter() {
		if setA.Has(item) {
			if setC.Has(item) {
				incommon++
			}
		} else if !setC.Has(item) {
			t.Errorf("Items in setAuC should be in either setA or setC")
		}
	}
	if setAuC.Cardinality() != setA.Cardinality() + setC.Cardinality() - incommon {
		t.Errorf("Cardinality of a union of intesecting sets should be the sum of their cardinalities minus the size of their intersection")
	}
}

func TestIntersection(t *testing.T) {
	setA := make_set_serial(-100, 0)
	setB := make_set_serial(1, 100)
	setC := make_set_serial(-50, 50)
	setAiB := Intersection(setA, setB)
	setAiC := Intersection(setA, setC)
	if (Intersect(setA, setAiB) && !Intersect(setB, setAiB)) || (!Intersect(setA, setAiB) && Intersect(setB, setAiB)) {
		t.Errorf("setAiB should intersect with both setA and SetB or neither")
	}
	if setAiB.Cardinality() != 0 {
		t.Errorf("Cardinality of an intersection of disjoint sets should be 0")
	}
	for item := range setAiB.Iter() {
		if !setA.Has(item) || !setB.Has(item) {
			t.Errorf("Items in setAiB should be in both setA and setB")
		}
	}
	if !Intersect(setA, setAiC) || !Intersect(setC, setAiC) {
		t.Errorf("setAiC should intersect with both setA and SetC")
	}
	for item := range setAiC.Iter() {
		if setA.Has(item) {
			if !setC.Has(item) {
				t.Errorf("Items in setAiC should be in both setA and setC")
			}
		} else if setC.Has(item) {
			t.Errorf("Items in setAiC should be in both setA and setC")
		}
	}
	if setAiC.Cardinality() != 51 {
		t.Errorf("Cardinality of an intersection of intesecting sets should be the size of their intersection")
	}
}

