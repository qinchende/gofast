package sqlx

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/store/orm"
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
func checkDestType(dest any) (*orm.ModelSchema, reflect.Type, reflect.Type, bool, bool) {
	dTyp := reflect.TypeOf(dest)
	if dTyp.Kind() != reflect.Ptr {
		cst.PanicString("Target object must be pointer.")
	}
	sliceType := dTyp.Elem()
	if sliceType.Kind() != reflect.Slice {
		cst.PanicString("Target object must be slice.")
	}
	ms := orm.SchemaOfType(dTyp)

	isPtr := false
	isKV := false
	recordType := sliceType.Elem()
	// 推荐: dest 传入的 slice 类型为指针类型，这样将来就不涉及变量值拷贝了。
	if recordType.Kind() == reflect.Ptr {
		isPtr = true
		recordType = recordType.Elem()
	} else {
		typName := recordType.Name()
		if typName == "cst.KV" || typName == "fst.KV" || typName == "KV" {
			isKV = true
		}
	}

	return ms, sliceType, recordType, isPtr, isKV
}
