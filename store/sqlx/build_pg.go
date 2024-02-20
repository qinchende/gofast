// Copyright 2023 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package sqlx

import (
	"github.com/qinchende/gofast/store/orm"
	"reflect"
)

// PostgreSQL
type PgBuilder struct {
}

func (*PgBuilder) Insert(ts *orm.TableSchema) string {
	return ""
}

func (*PgBuilder) Delete(ts *orm.TableSchema) string {
	return ""
}

func (*PgBuilder) Update(ts *orm.TableSchema) string {
	return ""
}

// 更新特定字段
func (*PgBuilder) UpdateColumns(ts *orm.TableSchema, rVal *reflect.Value, cNames ...string) (string, []any) {
	return "", nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 查询 select * from
func (*PgBuilder) SelectPrimary(ts *orm.TableSchema) string {
	return ""
}

func (*PgBuilder) SelectRow(ts *orm.TableSchema, columns string, where string) string {
	return ""
}

func (*PgBuilder) SelectRows(ts *orm.TableSchema, columns string, where string) string {
	return ""
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (*PgBuilder) InitPet(pet *SelectPet, ts *orm.TableSchema) {

}

func (*PgBuilder) SelectByPet(pet *SelectPet) string {
	return ""
}

func (*PgBuilder) SelectCountByPet(pet *SelectPet) string {
	return ""
}

func (*PgBuilder) SelectPagingByPet(pet *SelectPet) string {
	return ""
}
