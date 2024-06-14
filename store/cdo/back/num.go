package back

//
//import (
//	"golang.org/x/exp/constraints"
//	"unsafe"
//)
//
//// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
//func encU64VarUint(bs []byte, v uint64) []byte {
//	switch {
//	case v <= MaxUint05:
//		bs = append(bs, 0x00|(uint8(v)))
//	case v <= MaxUint13:
//		bs = append(bs, 0x20|byte(v>>8), uint8(v))
//	case v <= MaxUint21:
//		bs = append(bs, 0x40|byte(v>>16), byte(v>>8), byte(v))
//	}
//	return bs
//}
//
////go:noinline
//func encU64VarUint2(bs []byte, v uint64) []byte {
//	switch {
//	case v <= MaxUint24:
//		bs = append(bs, 0x63, byte(v), byte(v>>8), byte(v>>16))
//	case v <= MaxUint32:
//		bs = append(bs, 0x64, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
//	case v <= MaxUint40:
//		bs = append(bs, 0x65, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32))
//	case v <= MaxUint48:
//		bs = append(bs, 0x66, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40))
//	case v <= MaxUint56:
//		bs = append(bs, 0x67, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40), byte(v>>48))
//	case v <= MaxUint64:
//		bs = append(bs, 0x68, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40), byte(v>>48), byte(v>>56))
//	}
//	return bs
//}
//
//func encListInt[T constraints.Integer](se *subEncode, listSize int) {
//	encU24By5(se.bf, TypeList, uint64(listSize))
//
//	bs := *se.bf
//	bs = append(bs, ListVarInt)
//	for i := 0; i < listSize; i++ {
//		iPtr := unsafe.Add(se.srcPtr, i*se.em.itemMemSize)
//		v := *((*T)(iPtr))
//		if v >= 0 {
//			if uint64(v) <= MaxUint21 {
//				bs = encU64VarUint(bs, uint64(v))
//			} else {
//				bs = encU64VarUint2(bs, uint64(v))
//			}
//		} else {
//			if uint64(-v) <= MaxUint21 {
//				bs = encU64VarUint(bs, uint64(-v))
//			} else {
//				bs = encU64VarUint2(bs, uint64(-v))
//			}
//		}
//	}
//	*se.bf = bs
//}
//
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
//func decIntList(d *subDecode, listSize int) {
//	offS := d.scan
//	if d.str[offS] != ListVarInt {
//		panic(errChar)
//	}
//	offS++
//	for i := 0; i < listSize; i++ {
//		iPtr := unsafe.Add(d.dstPtr, i*d.dm.itemMemSize)
//
//		//Part0
//		c := d.str[offS]
//		typ := c & 0x80
//		v := uint64(c>>5) & 0x03
//		var off int
//		if v <= 2 {
//			off, v = scanVarInt(d.str[offS:], v)
//		} else {
//			off, v = scanVarInt2(d.str[offS:], uint64(c&0x0F))
//		}
//
//		if typ == 0x00 {
//			bindInt(iPtr, int64(v))
//		} else if typ == 0x80 {
//			bindInt(iPtr, int64(-v))
//		} else {
//			panic(errChar)
//		}
//		offS += off
//
//		//// Part1
//		//typ, v := typeValue(d.str[offS])
//		//var off int
//		//if v <= 59 {
//		//	off, v = scanU64Part1(d.str[offS:], v)
//		//} else {
//		//	off, v = scanU64Part2(d.str[offS:], v)
//		//}
//		//
//		//// Part3
//		////off, typ, v := scanTypeLen8(d.str[offS:])
//		//
//		//offS += off
//		//if typ == TypePosInt {
//		//	bindInt(iPtr, int64(v))
//		//} else if typ == TypeNegInt {
//		//	bindInt(iPtr, int64(-v))
//		//} else {
//		//	panic(errChar)
//		//}
//	}
//	d.scan = offS
//}
