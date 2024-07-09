package back

//
//import (
//	"golang.org/x/exp/constraints"
//	"unsafe"
//)

//	func encU16By5(bf *[]byte, sym uint8, v uint64) {
//		*bf = encU16By5Ret(*bf, sym, v)
//	}
//
//	func encU16By5Ret(bs []byte, sym uint8, v uint64) []byte {
//		switch {
//		default:
//			panic(errOutRange)
//		case v <= 29:
//			bs = append(bs, sym|(uint8(v)))
//		case v <= MaxUint08:
//			bs = append(bs, sym|30, uint8(v))
//		case v <= MaxUint16:
//			bs = append(bs, sym|31, byte(v), byte(v>>8))
//		}
//		return bs
//	}

//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func encU64By7RetPart1(bs []byte, sym byte, v uint64) []byte {
//	switch {
//	case v <= 119:
//		bs = append(bs, sym|(byte(v)))
//	case v <= MaxUint08:
//		bs = append(bs, sym|120, byte(v))
//	case v <= MaxUint16:
//		bs = append(bs, sym|121, byte(v), byte(v>>8))
//	case v <= MaxUint24:
//		bs = append(bs, sym|122, byte(v), byte(v>>8), byte(v>>16))
//	}
//	return bs
//}
//
////go:noinline
//func encU64By7RetPart2(bs []byte, sym byte, v uint64) []byte {
//	switch {
//	case v <= MaxUint32:
//		bs = append(bs, sym|123, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
//	case v <= MaxUint40:
//		bs = append(bs, sym|124, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32))
//	case v <= MaxUint48:
//		bs = append(bs, sym|125, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40))
//	case v <= MaxUint56:
//		bs = append(bs, sym|126, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40), byte(v>>48))
//	case v <= MaxUint64:
//		bs = append(bs, sym|127, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40), byte(v>>48), byte(v>>56))
//	}
//	return bs
//}
//

//func encListVarInt[T constraints.Integer](se *subEncode, tLen int) {
//	bs := *se.bf
//	bs = append(encU24By5Ret(bs, TypeList, uint64(tLen)), ListVarInt)
//	for i := 0; i < tLen; i++ {
//		iPtr := unsafe.Add(se.srcPtr, i*se.em.itemMemSize)
//		v := *((*T)(iPtr))
//		if v >= 0 {
//			if uint64(v) <= MaxUint24 {
//				bs = encU64By7RetPart1(bs, ListVarIntPos, uint64(v))
//			} else {
//				bs = encU64By7RetPart2(bs, ListVarIntPos, uint64(v))
//			}
//		} else {
//			if uint64(-v) <= MaxUint24 {
//				bs = encU64By7RetPart1(bs, ListVarIntNeg, uint64(-v))
//			} else {
//				bs = encU64By7RetPart2(bs, ListVarIntNeg, uint64(-v))
//			}
//		}
//	}
//	*se.bf = bs
//}
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

//func scanU64ValBy7(s string) (byte, int, uint64) {
//	typ, v := listVarIntHead(s[0])
//	var off int
//	if v <= 122 {
//		off, v = scanU64ValBy7Part1(s, v)
//	} else {
//		off, v = scanU64ValBy7Part2(s, v)
//	}
//	return typ, off, v
//}
//
//func scanU64ValBy7Part1(s string, v uint64) (int, uint64) {
//	if v <= 119 {
//		return 1, v
//	}
//	switch v {
//	case 120:
//		return 2, uint64(s[1])
//	case 121:
//		return 3, uint64(s[1]) | uint64(s[2])<<8
//	case 122:
//		return 4, uint64(s[1]) | uint64(s[2])<<8 | uint64(s[3])<<16
//	}
//	panic(errChar)
//}
//
////go:noinline
//func scanU64ValBy7Part2(s string, v uint64) (int, uint64) {
//	switch v {
//	case 123:
//		return 5, uint64(s[1]) | uint64(s[2])<<8 | uint64(s[3])<<16 | uint64(s[4])<<24
//	case 124:
//		return 6, uint64(s[1]) | uint64(s[2])<<8 | uint64(s[3])<<16 | uint64(s[4])<<24 | uint64(s[5])<<32
//	case 125:
//		return 7, uint64(s[1]) | uint64(s[2])<<8 | uint64(s[3])<<16 | uint64(s[4])<<24 | uint64(s[5])<<32 | uint64(s[6])<<40
//	case 126:
//		return 8, uint64(s[1]) | uint64(s[2])<<8 | uint64(s[3])<<16 | uint64(s[4])<<24 | uint64(s[5])<<32 | uint64(s[6])<<40 | uint64(s[7])<<48
//	case 127:
//		return 9, uint64(s[1]) | uint64(s[2])<<8 | uint64(s[3])<<16 | uint64(s[4])<<24 | uint64(s[5])<<32 | uint64(s[6])<<40 | uint64(s[7])<<48 | uint64(s[8])<<56
//	}
//	panic(errChar)
//}

//func decInt64List(d *subDecode, tLen int) {
//	list := *(*[]int64)(unsafe.Pointer(&d.slice))
//	pos := validListItemType(d, ListVarInt)
//	for i := 0; i < len(list); i++ {
//		sym, v := listVarIntHead(d.str[pos])
//		var off int
//		if v <= 122 {
//			off, v = scanU64ValBy7Part1(d.str[pos:], v)
//		} else {
//			off, v = scanU64ValBy7Part2(d.str[pos:], v)
//		}
//		list[i] = toInt64(sym, v)
//		pos += off
//	}
//	d.scan = pos
//}

//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
//func scanVarInt(s string, size uint64) (int, uint64) {
//	switch size {
//	default:
//		return 1, uint64(s[0] & 0x1F)
//	case 1:
//		return 2, uint64(s[0]&0x1F)<<8 | uint64(s[1])
//	case 2:
//		return 3, uint64(s[0]&0x1F)<<16 | uint64(s[1])<<8 | uint64(s[2])
//	}
//}
//
//func scanVarInt2(s string, size uint64) (int, uint64) {
//	switch size {
//	default:
//		panic(errChar)
//	case 3:
//		return 4, uint64(s[1]) | uint64(s[2])<<8 | uint64(s[3])<<16
//	case 4:
//		return 5, uint64(s[1]) | uint64(s[2])<<8 | uint64(s[3])<<16 | uint64(s[4])<<24
//	case 5:
//		return 6, uint64(s[1]) | uint64(s[2])<<8 | uint64(s[3])<<16 | uint64(s[4])<<24 | uint64(s[5])<<32
//	case 6:
//		return 7, uint64(s[1]) | uint64(s[2])<<8 | uint64(s[3])<<16 | uint64(s[4])<<24 | uint64(s[5])<<32 | uint64(s[6])<<40
//	case 7:
//		return 8, uint64(s[1]) | uint64(s[2])<<8 | uint64(s[3])<<16 | uint64(s[4])<<24 | uint64(s[5])<<32 | uint64(s[6])<<40 | uint64(s[7])<<48
//	case 8:
//		return 9, uint64(s[1]) | uint64(s[2])<<8 | uint64(s[3])<<16 | uint64(s[4])<<24 | uint64(s[5])<<32 | uint64(s[6])<<40 | uint64(s[7])<<48 | uint64(s[8])<<56
//	}
//}
//
//func decIntList(d *subDecode, tLen int) {
//	pos := d.scan
//	if d.str[pos] != ListVarInt {
//		panic(errChar)
//	}
//	pos++
//	for i := 0; i < tLen; i++ {
//		iPtr := unsafe.Add(d.dstPtr, i*d.dm.itemMemSize)
//
//		//Part0
//		c := d.str[pos]
//		typ := c & 0x80
//		v := uint64(c>>5) & 0x03
//		var off int
//		if v <= 2 {
//			off, v = scanVarInt(d.str[pos:], v)
//		} else {
//			off, v = scanVarInt2(d.str[pos:], uint64(c&0x0F))
//		}
//
//		if typ == 0x00 {
//			bindInt(iPtr, int64(v))
//		} else if typ == 0x80 {
//			bindInt(iPtr, int64(-v))
//		} else {
//			panic(errChar)
//		}
//		pos += off
//
//		//// Part1
//		//typ, v := typeValue(d.str[pos])
//		//var off int
//		//if v <= 59 {
//		//	off, v = scanU64Part1(d.str[pos:], v)
//		//} else {
//		//	off, v = scanU64Part2(d.str[pos:], v)
//		//}
//		//
//		//// Part3
//		////off, typ, v := scanTypeLen8(d.str[pos:])
//		//
//		//pos += off
//		//if typ == TypePosInt {
//		//	bindInt(iPtr, int64(v))
//		//} else if typ == TypeNegInt {
//		//	bindInt(iPtr, int64(-v))
//		//} else {
//		//	panic(errChar)
//		//}
//	}
//	d.scan = pos
//}
