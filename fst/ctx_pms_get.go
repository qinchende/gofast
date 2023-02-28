// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"errors"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/lang"
	"time"
)

var (
	errorKeyNotExist = errors.New("找不到参数值")
)

func (c *Context) Set(key string, val any) {
	if c.Pms == nil {
		c.Pms = make(cst.KV)
	}
	c.Pms[key] = val
}

func (c *Context) Get(key string) (val any, ok bool) {
	val, ok = c.Pms[key]
	return
}

func (c *Context) GetMust(key string) any {
	if val, ok := c.Pms[key]; ok {
		return val
	}
	cst.PanicIfErr(errorKeyNotExist)
	return nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (c *Context) GetString(key string) (string, error) {
	if v, ok := c.Pms[key]; ok {
		return lang.ToString2(v)
	}
	return "", errorKeyNotExist
}

func (c *Context) GetStringDef(key string, def string) string {
	if v, ok := c.Pms[key]; ok && v != nil {
		v2, err2 := lang.ToString2(v)
		cst.PanicIfErr(err2)
		return v2
	}
	return def
}

func (c *Context) GetStringMust(key string) string {
	v, err := lang.ToString2(c.GetMust(key))
	cst.PanicIfErr(err)
	return v
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (c *Context) GetBool(key string) (bool, error) {
	if v, ok := c.Pms[key]; ok {
		return lang.ToBool(v)
	}
	return false, errorKeyNotExist
}

func (c *Context) GetBoolDef(key string, def bool) bool {
	if v, ok := c.Pms[key]; ok && v != nil {
		v2, err2 := lang.ToBool(v)
		cst.PanicIfErr(err2)
		return v2
	}
	return def
}

func (c *Context) GetBoolMust(key string) bool {
	v, err := lang.ToBool(c.GetMust(key))
	cst.PanicIfErr(err)
	return v
}

func (c *Context) GetInt64(key string) (int64, error) {
	if v, ok := c.Pms[key]; ok {
		return lang.ToInt64(v)
	}
	return 0, errorKeyNotExist
}

func (c *Context) GetInt64Def(key string, def int64) int64 {
	if v, ok := c.Pms[key]; ok && v != nil {
		v2, err2 := lang.ToInt64(v)
		cst.PanicIfErr(err2)
		return v2
	}
	return def
}

func (c *Context) GetInt64Must(key string) int64 {
	v, err := lang.ToInt64(c.GetMust(key))
	cst.PanicIfErr(err)
	return v
}

func (c *Context) GetInt(key string) (int, error) {
	if v, ok := c.Pms[key]; ok {
		return lang.ToInt(v)
	}
	return 0, errorKeyNotExist
}

func (c *Context) GetIntDef(key string, def int) int {
	if v, ok := c.Pms[key]; ok && v != nil {
		v2, err2 := lang.ToInt(v)
		cst.PanicIfErr(err2)
		return v2
	}
	return def
}

func (c *Context) GetIntMust(key string) int {
	v, err := lang.ToInt(c.GetMust(key))
	cst.PanicIfErr(err)
	return v
}

func (c *Context) GetInt32(key string) (int32, error) {
	if v, ok := c.Pms[key]; ok {
		return lang.ToInt32(v)
	}
	return 0, errorKeyNotExist
}

func (c *Context) GetInt32Def(key string, def int32) int32 {
	if v, ok := c.Pms[key]; ok && v != nil {
		v2, err2 := lang.ToInt32(v)
		cst.PanicIfErr(err2)
		return v2
	}
	return def
}

func (c *Context) GetInt32Must(key string) int32 {
	v, err := lang.ToInt32(c.GetMust(key))
	cst.PanicIfErr(err)
	return v
}

func (c *Context) GetInt16(key string) (int16, error) {
	if v, ok := c.Pms[key]; ok {
		return lang.ToInt16(v)
	}
	return 0, errorKeyNotExist
}

func (c *Context) GetInt16Def(key string, def int16) int16 {
	if v, ok := c.Pms[key]; ok && v != nil {
		v2, err2 := lang.ToInt16(v)
		cst.PanicIfErr(err2)
		return v2
	}
	return def
}

func (c *Context) GetInt16Must(key string) int16 {
	v, err := lang.ToInt16(c.GetMust(key))
	cst.PanicIfErr(err)
	return v
}

func (c *Context) GetInt8(key string) (int8, error) {
	if v, ok := c.Pms[key]; ok {
		return lang.ToInt8(v)
	}
	return 0, errorKeyNotExist
}

func (c *Context) GetInt8Def(key string, def int8) int8 {
	if v, ok := c.Pms[key]; ok && v != nil {
		v2, err2 := lang.ToInt8(v)
		cst.PanicIfErr(err2)
		return v2
	}
	return def
}

func (c *Context) GetInt8Must(key string) int8 {
	v, err := lang.ToInt8(c.GetMust(key))
	cst.PanicIfErr(err)
	return v
}

func (c *Context) GetUint64(key string) (uint64, error) {
	if v, ok := c.Pms[key]; ok {
		return lang.ToUint64(v)
	}
	return 0, errorKeyNotExist
}

func (c *Context) GetUint64Def(key string, def uint64) uint64 {
	if v, ok := c.Pms[key]; ok && v != nil {
		v2, err2 := lang.ToUint64(v)
		cst.PanicIfErr(err2)
		return v2
	}
	return def
}

func (c *Context) GetUint64Must(key string) uint64 {
	v, err := lang.ToUint64(c.GetMust(key))
	cst.PanicIfErr(err)
	return v
}

func (c *Context) GetUint(key string) (uint, error) {
	if v, ok := c.Pms[key]; ok {
		return lang.ToUint(v)
	}
	return 0, errorKeyNotExist
}

func (c *Context) GetUintDef(key string, def uint) uint {
	if v, ok := c.Pms[key]; ok && v != nil {
		v2, err2 := lang.ToUint(v)
		cst.PanicIfErr(err2)
		return v2
	}
	return def
}

func (c *Context) GetUintMust(key string) uint {
	v, err := lang.ToUint(c.GetMust(key))
	cst.PanicIfErr(err)
	return v
}

func (c *Context) GetUint32(key string) (uint32, error) {
	if v, ok := c.Pms[key]; ok {
		return lang.ToUint32(v)
	}
	return 0, errorKeyNotExist
}

func (c *Context) GetUint32Def(key string, def uint32) uint32 {
	if v, ok := c.Pms[key]; ok && v != nil {
		v2, err2 := lang.ToUint32(v)
		cst.PanicIfErr(err2)
		return v2
	}
	return def
}

func (c *Context) GetUint32Must(key string) uint32 {
	v, err := lang.ToUint32(c.GetMust(key))
	cst.PanicIfErr(err)
	return v
}

func (c *Context) GetUint16(key string) (uint16, error) {
	if v, ok := c.Pms[key]; ok {
		return lang.ToUint16(v)
	}
	return 0, errorKeyNotExist
}

func (c *Context) GetUint16Def(key string, def uint16) uint16 {
	if v, ok := c.Pms[key]; ok && v != nil {
		v2, err2 := lang.ToUint16(v)
		cst.PanicIfErr(err2)
		return v2
	}
	return def
}

func (c *Context) GetUint16Must(key string) uint16 {
	v, err := lang.ToUint16(c.GetMust(key))
	cst.PanicIfErr(err)
	return v
}

func (c *Context) GetUint8(key string) (uint8, error) {
	if v, ok := c.Pms[key]; ok {
		return lang.ToUint8(v)
	}
	return 0, errorKeyNotExist
}

func (c *Context) GetUint8Def(key string, def uint8) uint8 {
	if v, ok := c.Pms[key]; ok && v != nil {
		v2, err2 := lang.ToUint8(v)
		cst.PanicIfErr(err2)
		return v2
	}
	return def
}

func (c *Context) GetUint8Must(key string) uint8 {
	v, err := lang.ToUint8(c.GetMust(key))
	cst.PanicIfErr(err)
	return v
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (c *Context) GetFloat64(key string) (float64, error) {
	if v, ok := c.Pms[key]; ok {
		return lang.ToFloat64(v)
	}
	return 0.0, errorKeyNotExist
}

func (c *Context) GetFloat64Def(key string, def float64) float64 {
	if v, ok := c.Pms[key]; ok && v != nil {
		v2, err2 := lang.ToFloat64(v)
		cst.PanicIfErr(err2)
		return v2
	}
	return def
}

func (c *Context) GetFloat64Must(key string) float64 {
	v, err := lang.ToFloat64(c.GetMust(key))
	cst.PanicIfErr(err)
	return v
}

func (c *Context) GetFloat32(key string) (float32, error) {
	if v, ok := c.Pms[key]; ok {
		return lang.ToFloat32(v)
	}
	return 0.0, errorKeyNotExist
}

func (c *Context) GetFloat32Def(key string, def float32) float32 {
	if v, ok := c.Pms[key]; ok && v != nil {
		v2, err2 := lang.ToFloat32(v)
		cst.PanicIfErr(err2)
		return v2
	}
	return def
}

func (c *Context) GetFloat32Must(key string) float32 {
	v, err := lang.ToFloat32(c.GetMust(key))
	cst.PanicIfErr(err)
	return v
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func (c *Context) GetTime(key string) (*time.Time, error) {
	if v, ok := c.Pms[key]; ok {
		return lang.ToTime("", v)
	}
	return nil, errorKeyNotExist
}

func (c *Context) GetTimeDef(key string, def *time.Time) *time.Time {
	if v, ok := c.Pms[key]; ok && v != nil {
		v2, err2 := lang.ToTime("", v)
		cst.PanicIfErr(err2)
		return v2
	}
	return def
}

func (c *Context) GetTimeMust(key string) *time.Time {
	v, err := lang.ToTime("", c.GetMust(key))
	cst.PanicIfErr(err)
	return v
}

func (c *Context) GetDuration(key string) (time.Duration, error) {
	if v, ok := c.Pms[key]; ok {
		return lang.ToDuration(v)
	}
	return 0, errorKeyNotExist
}

func (c *Context) GetDurationDef(key string, def time.Duration) time.Duration {
	if v, ok := c.Pms[key]; ok && v != nil {
		v2, err2 := lang.ToDuration(v)
		cst.PanicIfErr(err2)
		return v2
	}
	return def
}

func (c *Context) GetDurationMust(key string) time.Duration {
	v, err := lang.ToDuration(c.GetMust(key))
	cst.PanicIfErr(err)
	return v
}

//// GetDuration returns the value associated with the key as a duration.
//func (c *Context) GetDuration(key string) (d time.Duration) {
//	if val, ok := c.Pms[key]; ok && val != nil {
//		d, _ = val.(time.Duration)
//	}
//	return
//}

//// GetStringSlice returns the value associated with the key as a slice of strings.
//func (c *Context) GetStringSlice(key string) (ss []string) {
//	if val, ok := c.Pms[key]; ok && val != nil {
//		ss, _ = val.([]string)
//	}
//	return
//}
//
//// GetStringMap returns the value associated with the key as a map of interfaces.
//func (c *Context) GetStringMap(key string) (sm map[string]any) {
//	if val, ok := c.Pms[key]; ok && val != nil {
//		sm, _ = val.(map[string]any)
//	}
//	return
//}
//
//// GetStringMapString returns the value associated with the key as a map of strings.
//func (c *Context) GetStringMapString(key string) (sms map[string]string) {
//	if val, ok := c.Pms[key]; ok && val != nil {
//		sms, _ = val.(map[string]string)
//	}
//	return
//}
//
//// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
//func (c *Context) GetStringMapStringSlice(key string) (smss map[string][]string) {
//	if val, ok := c.Pms[key]; ok && val != nil {
//		smss, _ = val.(map[string][]string)
//	}
//	return
//}
