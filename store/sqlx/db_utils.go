package sqlx

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/store/orm"
	"github.com/samber/lo"
	"reflect"
)

func panicIfErr(err error) {
	if err != nil {
		logx.Error("sqlx: " + err.Error())
		cst.PanicIfErr(err)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Utils
func checkDestType(dest any) (*orm.TableSchema, reflect.Type, reflect.Type, bool, bool) {
	dTyp := reflect.TypeOf(dest)
	if dTyp.Kind() != reflect.Pointer {
		cst.PanicString("Target object must be pointer.")
	}
	sliceType := dTyp.Elem()
	if sliceType.Kind() != reflect.Slice {
		cst.PanicString("Target object must be slice.")
	}
	ts := orm.SchemaByType(dTyp)

	isPtr := false
	isKV := false
	recordType := sliceType.Elem()
	// 推荐: dest 传入的 slice 类型为指针类型，这样将来就不涉及变量值拷贝了。
	if recordType.Kind() == reflect.Pointer {
		isPtr = true
		recordType = recordType.Elem()
	} else {
		typName := recordType.Name()
		if lo.Contains[string]([]string{"KV", "cst.KV"}, typName) {
			isKV = true
		}
		//if typName == "cst.KV" || typName == "KV" {
		//	isKV = true
		//}
	}

	return ts, sliceType, recordType, isPtr, isKV
}
