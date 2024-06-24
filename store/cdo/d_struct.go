package cdo

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// int +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrIntValue(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		//off, v := scanVarInt(d.str[d.scan:])
		//d.scan += off
		//bindInt(fieldPtrDeep(d), v)
	}
}

func scanObjIntValue(d *subDecode) {
	off, typ, v := scanVarInt(d.str[d.scan:])
	d.scan += off
	bindInt(fieldPtr(d), typ, v)
}

// int8 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrInt8Value(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		//off, v := scanVarInt(d.str[d.scan:])
		//d.scan += off
		//bindInt8(fieldPtrDeep(d), v)
	}
}

func scanObjInt8Value(d *subDecode) {
	//off, v := scanVarInt(d.str[d.scan:])
	//d.scan += off
	//bindInt8(fieldPtr(d), v)

}

// int16 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrInt16Value(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		//off, v := scanVarInt(d.str[d.scan:])
		//d.scan += off
		//bindInt16(fieldPtrDeep(d), v)
	}
}

func scanObjInt16Value(d *subDecode) {
	//off, v := scanVarInt(d.str[d.scan:])
	//d.scan += off
	//bindInt16(fieldPtr(d), v)
}

// int32 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrInt32Value(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		//off, v := scanVarInt(d.str[d.scan:])
		//d.scan += off
		//bindInt32(fieldPtrDeep(d), v)
	}
}

func scanObjInt32Value(d *subDecode) {
	//off, v := scanVarInt(d.str[d.scan:])
	//d.scan += off
	//bindInt32(fieldPtr(d), v)
}

// int64 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrInt64Value(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		//off, v := scanVarInt(d.str[d.scan:])
		//d.scan += off
		//bindInt64(fieldPtrDeep(d), v)
	}
}

func scanObjInt64Value(d *subDecode) {
	//off, v := scanVarInt(d.str[d.scan:])
	//d.scan += off
	//bindInt64(fieldPtr(d), v)
}

// uint +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrUintValue(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		//off, v := scanUint64(d.str[d.scan:])
		//d.scan += off
		//bindUint(fieldPtrDeep(d), v)
	}
}

func scanObjUintValue(d *subDecode) {
	//off, v := scanUint64(d.str[d.scan:])
	//d.scan += off
	//bindUint(fieldPtr(d), v)
}

// uint8 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrUint8Value(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		//off, v := scanUint64(d.str[d.scan:])
		//d.scan += off
		//bindUint8(fieldPtrDeep(d), v)
	}
}

func scanObjUint8Value(d *subDecode) {
	//off, v := scanUint64(d.str[d.scan:])
	//d.scan += off
	//bindUint8(fieldPtr(d), v)
}

// uint16 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrUint16Value(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		//off, v := scanUint64(d.str[d.scan:])
		//d.scan += off
		//bindUint16(fieldPtrDeep(d), v)
	}
}

func scanObjUint16Value(d *subDecode) {
	//off, v := scanUint64(d.str[d.scan:])
	//d.scan += off
	//bindUint16(fieldPtr(d), v)
}

// uint32 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrUint32Value(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		//off, v := scanUint64(d.str[d.scan:])
		//d.scan += off
		//bindUint32(fieldPtrDeep(d), v)
	}
}

func scanObjUint32Value(d *subDecode) {
	//off, v := scanUint64(d.str[d.scan:])
	//d.scan += off
	//bindUint32(fieldPtr(d), v)
}

// uint64 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrUint64Value(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		//off, v := scanUint64(d.str[d.scan:])
		//d.scan += off
		//bindUint64(fieldPtrDeep(d), v)
	}
}

func scanObjUint64Value(d *subDecode) {
	//off, v := scanUint64(d.str[d.scan:])
	//d.scan += off
	//bindUint64(fieldPtr(d), v)
}

// float32
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrF32Value(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		v := scanF32Val(d.str[d.scan:])
		d.scan += 4
		bindF32(fieldPtrDeep(d), v)
	}
}

func scanObjF32Value(d *subDecode) {
	v := scanF32Val(d.str[d.scan:])
	d.scan += 4
	bindF32(fieldPtr(d), v)
}

// float64
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrF64Value(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		v := scanF64Val(d.str[d.scan:])
		d.scan += 4
		bindF64(fieldPtrDeep(d), v)
	}
}

func scanObjF64Value(d *subDecode) {
	v := scanF64Val(d.str[d.scan:])
	d.scan += 4
	bindF64(fieldPtr(d), v)
}

// string +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrStrValue(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		off, str := scanString(d.str[d.scan:])
		d.scan += off
		bindString(fieldPtrDeep(d), str)
	}
}

func scanObjStrValue(d *subDecode) {
	off, str := scanString(d.str[d.scan:])
	d.scan += off
	bindString(fieldPtr(d), str)
}

// []byte +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjBytesValue(d *subDecode) {
	off, bs := scanBytes(d.str[d.scan:])
	d.scan += off
	bindBytes(fieldPtr(d), bs)
}

// time.Time +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjTimeValue(d *subDecode) {
	switch d.str[d.scan] {
	case '"':
		start := d.scan + 1
		//d.scanQuoteStr()
		bindTime(fieldPtr(d), d.str[start:d.scan-1])
	default:
	}
}

func scanObjPtrTimeValue(d *subDecode) {
	switch d.str[d.scan] {
	case '"':
		start := d.scan + 1
		//d.scanQuoteStr()
		bindTime(fieldPtrDeep(d), d.str[start:d.scan-1])
	default:
		fieldSetNil(d)
	}
}

func scanArrTimeValue(d *subDecode) {
	v := ""
	switch d.str[d.scan] {
	case '"':
		start := d.scan + 1
		//d.scanQuoteStr()
		v = d.str[start : d.scan-1]
	default:
	}
	bindTime(arrItemPtr(d), v)
}

//func scanListTimeValue(d *subDecode) {
//	v := false
//	switch d.str[d.scan] {
//	case 't':
//		//d.skipTrue()
//		v = true
//	case 'f':
//		//d.skipFalse()
//	default:
//		//d.skipNull()
//		d.pl.nulPos = append(d.pl.nulPos, len(d.pl.bufBol))
//	}
//	d.pl.bufBol = append(d.pl.bufBol, v)
//}

// bool +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjBoolValue(d *subDecode) {
	v := scanBoolVal(d.str[d.scan:])
	d.scan += 1
	bindBool(fieldPtr(d), v)
}

func scanObjPtrBoolValue(d *subDecode) {

}

// any +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjAnyValue(d *subDecode) {
}

func scanObjPtrAnyValue(d *subDecode) {
}

func scanArrAnyValue(d *subDecode) {
}

func scanListAnyValue(d *subDecode) {
}
