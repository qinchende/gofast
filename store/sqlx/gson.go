package sqlx

import (
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/mapx"
	"github.com/qinchende/gofast/store/gson"
	"github.com/qinchende/gofast/store/orm"
	"reflect"
)

type gsonResultOne struct {
	gson.GsonOne
	hasValue bool
}

// 缓存实体 gsonResult
type gsonResult struct {
	gson.Gson
	onlyGson bool
}

func loadRecordFromGsonString(dest any, data string, sm *orm.ModelSchema) error {
	var values []any
	if err := jsonx.UnmarshalFromString(&values, data); err != nil {
		return err
	}

	cls := sm.Columns()
	recordKV := make(map[string]any, len(cls))
	for j := 0; j < len(cls); j++ {
		recordKV[cls[j]] = values[j]
	}

	return mapx.ApplyKVOfData(dest, recordKV)
}

func loadRecordsFromGsonString(dest any, data string, gr *gsonResult) error {
	if err := jsonx.UnmarshalFromString(&gr.Gson, data); err != nil {
		return err
	}

	_, sliceType, recordType, isPtr, isKV := checkDestType(dest)
	tpRecords := make([]reflect.Value, 0, gr.Ct)

	// 循环解析每一条记录
	for i := int64(0); i < gr.Ct; i++ {
		row := gr.Rows[i]
		recordKV := make(map[string]any, len(gr.Cls))
		for j := 0; j < len(gr.Cls); j++ {
			recordKV[gr.Cls[j]] = row[j]
		}

		if isKV {
			tpRecords = append(tpRecords, reflect.ValueOf(recordKV))
		} else {
			recordPtr := reflect.New(recordType)
			recordVal := reflect.Indirect(recordPtr)

			if err := mapx.ApplyKVOfData(recordVal.Addr().Interface(), recordKV); err != nil {
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
	reflect.ValueOf(dest).Elem().Set(records)
	return nil
}
