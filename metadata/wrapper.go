package metadata

import (
	"time"

	"github.com/meta-apex/gopkg/cast"
)

func (m Metadata) GetString(key string, def ...string) (val string, ok bool) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, ok = v.(string)
	return
}

func (m Metadata) GetToString(key string, def ...string) (val string, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, err = cast.ToStringE(v)
	if err != nil {
		return
	}
	return
}

func (m Metadata) GetInt(key string, def ...int) (val int, ok bool) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, ok = v.(int)
	return
}

func (m Metadata) GetToInt(key string, def ...int) (val int, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, err = cast.ToIntE(v)
	if err != nil {
		return
	}
	return
}

func (m Metadata) GetUInt(key string, def ...uint) (val uint, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, ok = v.(uint)
	return
}

func (m Metadata) GetToUInt(key string, def ...uint) (val uint, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, err = cast.ToUintE(v)
	if err != nil {
		return
	}
	return
}

func (m Metadata) GetInt16(key string, def ...int16) (val int16, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, ok = v.(int16)
	return
}

func (m Metadata) GetToInt16(key string, def ...int16) (val int16, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, err = cast.ToInt16E(v)
	if err != nil {
		return
	}
	return
}

func (m Metadata) GetUInt16(key string, def ...uint16) (val uint16, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, ok = v.(uint16)
	return
}

func (m Metadata) GetToUInt16(key string, def ...uint16) (val uint16, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, err = cast.ToUint16E(v)
	if err != nil {
		return
	}
	return
}

func (m Metadata) GetInt32(key string, def ...int32) (val int32, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, ok = v.(int32)
	return
}

func (m Metadata) GetToInt32(key string, def ...int32) (val int32, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, err = cast.ToInt32E(v)
	if err != nil {
		return
	}
	return
}

func (m Metadata) GetUInt32(key string, def ...uint32) (val uint32, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, ok = v.(uint32)
	return
}

func (m Metadata) GetToUInt32(key string, def ...uint32) (val uint32, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, err = cast.ToUint32E(v)
	if err != nil {
		return
	}
	return
}

func (m Metadata) GetInt64(key string, def ...int64) (val int64, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, ok = v.(int64)
	return
}

func (m Metadata) GetToInt64(key string, def ...int64) (val int64, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, err = cast.ToInt64E(v)
	if err != nil {
		return
	}
	return
}

func (m Metadata) GetUInt64(key string, def ...uint64) (val uint64, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, ok = v.(uint64)
	return
}

func (m Metadata) GetToUInt64(key string, def ...uint64) (val uint64, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, err = cast.ToUint64E(v)
	if err != nil {
		return
	}
	return
}

func (m Metadata) GetBool(key string, def ...bool) (val bool, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, ok = v.(bool)
	return
}

func (m Metadata) GetToBool(key string, def ...bool) (val bool, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, err = cast.ToBoolE(v)
	if err != nil {
		return
	}
	return
}

func (m Metadata) GetFloat32(key string, def ...float32) (val float32, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, ok = v.(float32)
	return
}

func (m Metadata) GetToFloat32(key string, def ...float32) (val float32, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, err = cast.ToFloat32E(v)
	if err != nil {
		return
	}
	return
}

func (m Metadata) GetFloat64(key string, def ...float64) (val float64, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, ok = v.(float64)
	return
}

func (m Metadata) GetToFloat64(key string, def ...float64) (val float64, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, err = cast.ToFloat64E(v)
	if err != nil {
		return
	}
	return
}

func (m Metadata) GetDuration(key string, def ...time.Duration) (val time.Duration, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, ok = v.(time.Duration)
	return
}

func (m Metadata) GetToDuration(key string, def ...time.Duration) (val time.Duration, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, err = cast.ToDurationE(v)
	if err != nil {
		return
	}
	return
}

func (m Metadata) GetTime(key string, def ...time.Time) (val time.Time, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, ok = v.(time.Time)
	return
}

func (m Metadata) GetToTime(key string, def ...time.Time) (val time.Time, ok bool, err error) {
	if len(def) > 0 {
		val = def[0]
	}

	v, ok := m[key]
	if !ok {
		return
	}

	val, err = cast.ToTimeE(v)
	if err != nil {
		return
	}
	return
}
