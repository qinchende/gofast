package mapx

import (
	"reflect"
	"sync"

	"github.com/go-playground/validator/v10"
)

const validConfigTag = "valid" // 指定模型中验证字段的tag标记

type StructValidator interface {
	ValidateStruct(interface{}) error
	Engine() interface{}
}

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var validMaster StructValidator = &defaultValidator{}

// 初始化, add by cd.net on 20220316
func init() {
	validMaster.Engine()
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Validate(obj interface{}) error {
	if validMaster == nil {
		return nil
	}
	return validMaster.ValidateStruct(obj)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	value := reflect.ValueOf(obj)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	if valueType == reflect.Struct {
		//v.lazyInit()
		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}
	return nil
}

func (v *defaultValidator) Engine() interface{} {
	v.lazyInit()
	return v.validate
}

func (v *defaultValidator) lazyInit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName(validConfigTag)
	})
}
