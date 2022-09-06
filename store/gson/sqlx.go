package gson

import (
	"github.com/qinchende/gofast/store/sqlx"
)

func LoadFromString(data string, pet *sqlx.SelectPet) int64 {

	//dest := cst.KV{"count": 0, "records": cst.KV{}}
	//if err = jsonx.UnmarshalFromString(&dest, cValStr); err == nil {
	//	err := mapx.ApplySliceOfConfig(pet.Target, dest["records"])
	//	sqlx.ErrPanic(err)
	//	ct, _ := lang.ToInt64(dest["count"])
	//	return ct
	//}
	return 0
}
