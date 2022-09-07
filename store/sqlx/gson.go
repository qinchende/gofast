package sqlx

import (
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/mapx"
	"reflect"
)

// 缓存实体 gsonResult
type gsonResult struct {
	Ct   int64
	Tt   int64
	Cls  []string
	Rows [][]any
}

func loadFromGsonString(dest any, data string, gr *gsonResult) error {
	if err := jsonx.UnmarshalFromString(gr, data); err != nil {
		return err
	}

	sliceType, recordType, isPtr, isKV := checkDestType(dest)
	tpRecords := make([]reflect.Value, 0, gr.Ct)

	// 循环解析每一条记录
	for i := int64(0); i < gr.Ct; i++ {
		row := gr.Rows[i]
		record := make(map[string]any, len(gr.Cls))
		for j := 0; j < len(gr.Cls); j++ {
			record[gr.Cls[j]] = row[j]
		}

		if isKV {
			tpRecords = append(tpRecords, reflect.ValueOf(record))
		} else {
			recordPtr := reflect.New(recordType)
			recordVal := reflect.Indirect(recordPtr)

			if err := mapx.ApplyKVOfData(recordVal.Addr().Interface(), record); err != nil {
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
