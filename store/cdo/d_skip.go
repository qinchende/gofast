package cdo

func (d *decoder) skipListStr(size int) {
	for i := 0; i < size; i++ {
		d.skipList()
	}
}

func (d *decoder) skipListAny(size int) {
	for i := 0; i < size; i++ {
		d.skipList()
	}
}

func (d *decoder) skipListVarInt(size int) {
	for i := 0; i < size; i++ {
		d.skipList()
	}
}

func (d *decoder) skipListMap(size int) {
	//off1, _, size := scanTypeU16(d.str[d.scan:])
	//if typ != TypeMap {
	//	panic(errKV)
	//}

	//for i := 0; i < int(size); i++ {
	//	off2 := skipString(d.str[off1:])
	//	d.scan += off2 + off1
	//	d.skipOneValue()
	//}
}

func (d *decoder) skipListList(size int) {
	for i := 0; i < size; i++ {
		d.skipList()
	}
}

func (d *decoder) skipListStruct(size int) {
	for i := 0; i < size; i++ {
		d.skipList()
	}
}

// skip items
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (d *decoder) skipList() {
	// 解析List长度
	off1, tLen := decListTypeU24(d.str[d.scan:])
	lSize := int(tLen)
	d.scan += off1
	// 解析List项目类型 和 相应值
	off2, typ, val := decListSubtypeU16(d.str[d.scan:])
	d.scan += off2

	switch typ {
	default:
		panic(errCdoChar)
	case ListBaseType:
		switch val {
		default:
			panic(errListType)
		case ListVarInt:
			d.skipListVarInt(lSize)
		case ListBool:
			d.scan += 1 * lSize
		case ListF32:
			d.scan += 4 * lSize
		case ListF64, ListTime:
			d.scan += 8 * lSize
		case ListStr:
			d.skipListStr(lSize)
		case ListAny:
			d.skipListAny(lSize)
		case ListKV:
			d.skipListMap(lSize)
		case ListList:
			d.skipListList(lSize)
		}
	case ListObjFields:
		d.skipListStruct(lSize)
		//case ListObjIndex:
		//case ListExt:
	}
}

func (d *decoder) skipOneValue() {
	str := d.str[d.scan:]
	c := str[0]
	typ := c & TypeMask
	val := c & TypeValMask
	off := 0

	switch typ {
	case TypeFixed:
		switch val {
		default:
			if val >= TypeList {
				d.skipList()
				return
			}
			panic(errCdoChar)
		case FixNil, FixNilMixed, FixTrue, FixFalse:
			off = 1
		case FixF32:
			off = 5
		case FixF64:
			off = 9
		case FixTime:
			off = 5
		}
	case TypeVarIntPos, TypeVarIntNeg:
		off = 1
		if val > 55 {
			off += int(val - 55)
		}
	case TypeStr:
		off = scanStringLen(str)
	}

	d.scan += off // 往前跳过相应字节
}
