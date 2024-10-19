package bind

import (
	"errors"
	"github.com/qinchende/gofast/aid/validx"
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/core/dts"
	"reflect"
	"unsafe"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func ValidateStruct(dst any, opts *dts.BindOptions) error {
	if dst == nil || opts == nil {
		return errors.New("has nil param.")
	}
	if opts.UseValid == false {
		return nil
	}
	if dstType, ptr, err := checkStructDest(dst); err != nil {
		return err
	} else {
		return validStructIter(ptr, dstType, opts)
	}
}

func validStructIter(ptr unsafe.Pointer, dstT reflect.Type, opts *dts.BindOptions) (err error) {
	sm := dts.SchemaByType(dstT, opts)

	for i := 0; i < len(sm.Fields); i++ {
		fa := &sm.FieldsAttr[i] // 肯定不会是nil
		fKind := fa.Kind
		fPtr := fa.MyPtr(ptr)

		// 如果字段是结构体类型
		if fKind == reflect.Struct && fa.Type != cst.TypeTime {
			// Note：指向其它结构体的字段，不做处理，防止嵌套循环验证
			if fa.PtrLevel > 0 {
				continue
			}
			if err = validStructIter(fPtr, fa.Type, opts); err != nil {
				return
			}
			continue
		}

		vOpt := fa.Valid
		if vOpt == nil {
			continue
		}
		fPtr = dts.PeelPtr(fPtr, fa.PtrLevel)
		if err = validx.ValidateFieldPtr(fPtr, fKind, vOpt); err != nil {
			return
		}
	}
	return
}
