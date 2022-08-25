package orm

import (
	"database/sql"
	"reflect"
	"time"
)

var modelAttrsList map[string]*ModelAttrs

func ShareModelAttrs(list map[string]*ModelAttrs) {
	modelAttrsList = list
}

// dbc: 数据库相关的配置参数
// dbf: 数据库字段的名称
// pms: 绑定数值时候的字段名称
// valid: 验证命令配置

// GoFast框架的ORM定义，所有Model必须公用的方法
type CommonFields struct {
	ID        int64 `dbc:"primary_field"`
	Status    int16
	CreatedAt time.Time `dbc:"created_field"`
	UpdatedAt time.Time `dbc:"updated_field"`
}

func (cf *CommonFields) GfAttrs(parent OrmStruct) *ModelAttrs {
	if modelAttrsList != nil {
		fullName := ""
		if parent != nil {
			fullName = reflect.TypeOf(parent).Elem().String()
		}
		if attr := modelAttrsList[fullName]; attr != nil {
			return attr
		}
	}
	return &ModelAttrs{}
}

// 万一更新失败，这里的值已经修改，需要回滚吗？？？
func (cf *CommonFields) BeforeSave() {
	if cf.ID == 0 || cf.CreatedAt.IsZero() {
		cf.CreatedAt = time.Now()
	}
	cf.UpdatedAt = time.Now()
}

func (cf *CommonFields) AfterInsert(result sql.Result) {
	lstId, err := result.LastInsertId()
	if err == nil {
		cf.ID = lstId
	} else {
		cf.CreatedAt = time.Time{}
		cf.UpdatedAt = time.Time{}
	}
}
