// Copyright 2010 -- Peter Williams, all rights reserved
// Use of this source code is governed by the new BSD license.

package set

import (
	"testing"
	"rand"
	"reflect"
	"mudlark/tree/llrb_tree"
)

func TestMakeBitSet(t *testing.T) {
	set := makebitset()
	if reflect.Typeof(set).String() != "*set.bitset" {
		t.Errorf("Expected type \"*bitset\": got %v", reflect.Typeof(set).String())
	}
	if set.bitcount != 0 {
		t.Errorf("Expected bitcount 0: got %v", set.bitcount)
	}
	if set.bits == nil {
		t.Errorf("Bit map unitialized")
	}
	if set.bits.Len() != 0 {
		t.Errorf("Expected len(bits) 0: got %v", set.bits.Len())
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

func checkbitcount(bset *bitset, str string, t *testing.T) {
	var count uint64 = 0
	for record := range bset.bits.Iter(llrb_tree.IN_ORDER) {
		count += uint64(bitcount(record.(*bitrecord).chunk))
	}
	if count != bset.bitcount {
		t.Errorf("Bit count %s. Expected: %v got: %v", str, bset.bitcount, count)
	}
}

func TestKeyBitcountAddAndRemove(t *testing.T) {
	const loopsz = 1000
	bset := makebitset()
	for i := 0; i < loopsz; i++ {
		bset.setBit(ibitlocation(i))
		checkbitcount(bset, "add(sequence)", t)
	}
	for i := 0; i < loopsz; i++ {
		bset.setBit(ibitlocation(rand.Int63()))
		checkbitcount(bset, "add(random(spread))", t)
	}
	for i := 0; i < loopsz; i++ {
		bset.setBit(ibitlocation(rand.Int63n(loopsz * 2)))
		checkbitcount(bset, "add(random(focussed))", t)
	}
	for i := 0; i < loopsz; i++ {
		bset.clearBit(ibitlocation(rand.Int63()))
		checkbitcount(bset, "remove(random(spread))", t)
	}
	for i := 0; i < loopsz; i++ {
		bset.clearBit(ibitlocation(rand.Int63n(loopsz * 2)))
		checkbitcount(bset, "remove(random(focussed))", t)
	}
	for i := 0; i < loopsz; i++ {
		bset.clearBit(ibitlocation(i))
		checkbitcount(bset, "remove(sequence)", t)
	}
}

func BenchmarkMakeEmptyBitSet(b *testing.B) {
	b.SetBytes(1)
	for i := 0; i < b.N; i++ {
		set := makebitset()
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
		set := makebitset()
		b.StartTimer()
		for i := 0; i < N; i++ {
			set.setBit(ibitlocation(rand.Int()))
		}
	}
}

func BenchmarkInsertSerial(b *testing.B) {
	const N = 50000
	b.SetBytes(N)
	for ib := 0; ib < b.N; ib++ {
		b.StopTimer()
		set := makebitset()
		b.StartTimer()
		for i := 0; i < N; i++ {
			set.setBit(ibitlocation(i))
		}
	}
}

