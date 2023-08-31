package orm

import (
	"database/sql"
	"reflect"
	"time"
)

var tableAttrsList map[string]*TableAttrs

func ShareTableAttrs(list map[string]*TableAttrs) {
	tableAttrsList = list
}

// dbc: 数据库相关的配置参数
// dbf: 数据库字段的名称
// pms: 绑定数值时候的字段名称
// valid: 验证命令配置

// GoFast框架的ORM定义，所有Model必须公用的方法
type CommonFields struct {
	ID        int64     // `dbc:"primary_field"`
	Status    int16     // `dbc:"status_field"`
	CreatedAt time.Time // `dbc:"created_field"`
	UpdatedAt time.Time // `dbc:"updated_field"`
}

func (cf *CommonFields) GfAttrs(parent OrmStruct) (attr *TableAttrs) {
	if tableAttrsList != nil {
		fullName := ""
		if parent != nil {
			fullName = reflect.TypeOf(parent).Elem().String()
		}
		attr = tableAttrsList[fullName]
	}
	if attr == nil {
		attr = &TableAttrs{}
	}
	//_ = mapx.Optimize(attr, mapx.LikeConfig) // 添加默认值，验证字段
	return
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
