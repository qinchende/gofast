package sqlx

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/mapx"
	"github.com/qinchende/gofast/store/gson"
	"github.com/qinchende/gofast/store/jde"
	"github.com/qinchende/gofast/store/orm"
	"reflect"
)

type gsonResultOne struct {
	gson.GsonRow
	hasValue bool
}

// 缓存实体 gsonResult
type gsonResult struct {
	gson.GsonRows
	onlyGson bool
}

// 将GsonRow记录值（仅仅是Value部分），绑定到对象中
func bindFromGsonValueString(obj any, bs []byte, ts *orm.TableSchema) error {
	//var values []any
	//if err := jsonx.UnmarshalFromString(&values, data); err != nil {
	//	return err
	//}
	//
	//cls := ts.Columns()
	//recordKV := make(cst.KV, len(cls))
	//for j := 0; j < len(cls); j++ {
	//	recordKV[cls[j]] = values[j]
	//}
	//
	//return mapx.BindKV(obj, recordKV, mapx.LikeLoadDB)
	//return nil

	return jde.DecodeGsonRowFromValueBytes(obj, bs)
}

// GsonRows的序列字符串绑定到对象数组中
func loadRecordsFromGsonString(objs any, data string, gr *gsonResult) error {
	if err := jsonx.UnmarshalFromString(&gr.GsonRows, data); err != nil {
		return err
	}

	_, sliceType, recordType, isPtr, isKV := checkDestType(objs)
	tpRecords := make([]reflect.Value, 0, gr.Ct)

	// 循环解析每一条记录
	for i := int64(0); i < gr.Ct; i++ {
		row := gr.Rows[i]
		recordKV := make(cst.KV, len(gr.Cls))
		for j := 0; j < len(gr.Cls); j++ {
			recordKV[gr.Cls[j]] = row[j]
		}

		if isKV {
			tpRecords = append(tpRecords, reflect.ValueOf(recordKV))
		} else {
			recordPtr := reflect.New(recordType)
			recordVal := reflect.Indirect(recordPtr)

			if err := mapx.BindKV(recordVal.Addr().Interface(), recordKV, mapx.LikeLoadDB); err != nil {
				return err
			}

			if isPtr {
				tpRecords = append(tpRecords, recordPtr)
			} else {
				tpRecords = append(tpRecords, recordVal)
			}
		}
	}

	records := reflect.MakeSlice(sliceType, 0, len(tpRecords))
	records = reflect.Append(records, tpRecords...)
	reflect.ValueOf(objs).Elem().Set(records)
	return nil
}
