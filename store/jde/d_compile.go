package jde

import (
	"github.com/qinchende/gofast/store/dts"
	"reflect"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 如果不是map和*GsonRow，只能是 Array|Slice|Struct
func (sd *subDecode) buildMeta(dst any) (err error) {
	sd.dm = &destMeta{}

	rfVal := reflect.ValueOf(dst)
	if rfVal.Kind() != reflect.Pointer {
		return errValueMustPtr
	}
	rfVal = reflect.Indirect(rfVal)

	rfTyp := rfVal.Type()
	switch kd := rfTyp.Kind(); kd {
	case reflect.Struct:
		if rfTyp.String() == "time.Time" {
			return errValueType
		}
		if err = sd.initStructMeta(rfTyp); err != nil {
			return
		}
		// 初始化解析函数 ++++++++++++++++++
		sd.dm.ssFunc = make([]structDecode, len(sd.dm.ss.FieldsKind))
		for i := 0; i < len(sd.dm.ss.FieldsKind); i++ {
			switch sd.dm.ss.FieldsKind[i] {
			case reflect.Int:
				sd.dm.ssFunc[i] = scanObjIntValue
			case reflect.Int8:
				sd.dm.ssFunc[i] = scanObjInt8Value
			case reflect.Int16:
				sd.dm.ssFunc[i] = scanObjInt16Value
			case reflect.Int32:
				sd.dm.ssFunc[i] = scanObjInt32Value
			case reflect.Int64:
				sd.dm.ssFunc[i] = scanObjInt64Value
			case reflect.Uint:
				sd.dm.ssFunc[i] = scanObjUintValue
			case reflect.Uint8:
				sd.dm.ssFunc[i] = scanObjUint8Value
			case reflect.Uint16:
				sd.dm.ssFunc[i] = scanObjUint16Value
			case reflect.Uint32:
				sd.dm.ssFunc[i] = scanObjUint32Value
			case reflect.Uint64:
				sd.dm.ssFunc[i] = scanObjUint64Value
			case reflect.Float32:
				sd.dm.ssFunc[i] = scanObjFloat32Value
			case reflect.Float64:
				sd.dm.ssFunc[i] = scanObjFloat64Value
			case reflect.String:
				sd.dm.ssFunc[i] = scanObjStrValue
			case reflect.Bool:
				sd.dm.ssFunc[i] = scanObjBolValue
			case reflect.Interface:

			case reflect.Map:

			case reflect.Slice:

			case reflect.Array:

			case reflect.Struct:

			}
		}
		sd.dm.nextPtr = fieldPtr
	case reflect.Array, reflect.Slice:
		if err = sd.initListMeta(rfTyp); err != nil {
			return
		}

		switch sd.dm.itemKind {
		case reflect.Int:
			sd.dm.itemFunc = scanObjIntValue
		case reflect.Int8:
			sd.dm.itemFunc = scanObjInt8Value
		case reflect.Int16:
			sd.dm.itemFunc = scanObjInt16Value
		case reflect.Int32:
			sd.dm.itemFunc = scanObjInt32Value
		case reflect.Int64:
			sd.dm.itemFunc = scanObjInt64Value
		case reflect.Uint:
			sd.dm.itemFunc = scanObjUintValue
		case reflect.Uint8:
			sd.dm.itemFunc = scanObjUint8Value
		case reflect.Uint16:
			sd.dm.itemFunc = scanObjUint16Value
		case reflect.Uint32:
			sd.dm.itemFunc = scanObjUint32Value
		case reflect.Uint64:
			sd.dm.itemFunc = scanObjUint64Value
		case reflect.Float32:
			sd.dm.itemFunc = scanObjFloat32Value
		case reflect.Float64:
			sd.dm.itemFunc = scanObjFloat64Value
		case reflect.String:
			sd.dm.itemFunc = scanObjStrValue
		case reflect.Bool:
			sd.dm.itemFunc = scanObjBolValue
		case reflect.Interface:

		case reflect.Map:

		case reflect.Slice:

		case reflect.Array:

		case reflect.Struct:

		}
		sd.dm.nextPtr = listItemPtr

		// 进一步初始化数组
		if kd == reflect.Array {
			sd.initArrayMeta()
			sd.dm.arrLen = rfVal.Len()
		}
	default:
		return errValueType
	}
	return nil
}

func (sd *subDecode) initStructMeta(rfType reflect.Type) error {
	sd.dm.isStruct = true
	sd.dm.ss = dts.SchemaForInputByType(rfType)
	return nil
}

func (sd *subDecode) initListMeta(rfType reflect.Type) error {
	sd.dm.isList = true

	sd.dm.listType = rfType
	sd.dm.listKind = rfType.Kind()
	sd.dm.itemType = sd.dm.listType.Elem()
	sd.dm.itemKind = sd.dm.itemType.Kind()

peelPtr:
	if sd.dm.itemKind == reflect.Pointer {
		sd.dm.itemType = sd.dm.itemType.Elem()
		sd.dm.itemKind = sd.dm.itemType.Kind()
		sd.dm.isPtr = true
		sd.dm.ptrLevel++
		// TODO：指针嵌套不能超过3层，这种很少见，真遇到，自己后期再处理
		if sd.dm.ptrLevel > 3 {
			return errPtrLevel
		}
		goto peelPtr
	}

	// 是否是interface类型
	if sd.dm.itemKind == reflect.Interface {
		sd.dm.isAny = true
	}
	return nil
}

func (sd *subDecode) initArrayMeta() {
	sd.dm.isArray = true
	if sd.dm.isPtr {
		return
	}
	sd.dm.isArrBind = true
	sd.dm.itemSize = int(sd.dm.itemType.Size())
}
