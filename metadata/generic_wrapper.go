package metadata

import (
	"time"

	"github.com/meta-apex/gopkg/cast"
)

// GetString returns string value for the given key
func (m Generic[K]) GetString(key K, def ...string) (val string, ok bool) {
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

// GetToString returns string value with type conversion for the given key
func (m Generic[K]) GetToString(key K, def ...string) (val string, ok bool, err error) {
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

// GetInt returns int value for the given key
func (m Generic[K]) GetInt(key K, def ...int) (val int, ok bool) {
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

// GetToInt returns int value with type conversion for the given key
func (m Generic[K]) GetToInt(key K, def ...int) (val int, ok bool, err error) {
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

// GetUInt returns uint value for the given key
func (m Generic[K]) GetUInt(key K, def ...uint) (val uint, ok bool) {
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

// GetToUInt returns uint value with type conversion for the given key
func (m Generic[K]) GetToUInt(key K, def ...uint) (val uint, ok bool, err error) {
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

// GetInt16 returns int16 value for the given key
func (m Generic[K]) GetInt16(key K, def ...int16) (val int16, ok bool) {
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

// GetToInt16 returns int16 value with type conversion for the given key
func (m Generic[K]) GetToInt16(key K, def ...int16) (val int16, ok bool, err error) {
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

// GetUInt16 returns uint16 value for the given key
func (m Generic[K]) GetUInt16(key K, def ...uint16) (val uint16, ok bool) {
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

// GetToUInt16 returns uint16 value with type conversion for the given key
func (m Generic[K]) GetToUInt16(key K, def ...uint16) (val uint16, ok bool, err error) {
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

// GetInt32 returns int32 value for the given key
func (m Generic[K]) GetInt32(key K, def ...int32) (val int32, ok bool) {
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

// GetToInt32 returns int32 value with type conversion for the given key
func (m Generic[K]) GetToInt32(key K, def ...int32) (val int32, ok bool, err error) {
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

// GetUInt32 returns uint32 value for the given key
func (m Generic[K]) GetUInt32(key K, def ...uint32) (val uint32, ok bool) {
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

// GetToUInt32 returns uint32 value with type conversion for the given key
func (m Generic[K]) GetToUInt32(key K, def ...uint32) (val uint32, ok bool, err error) {
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

// GetInt64 returns int64 value for the given key
func (m Generic[K]) GetInt64(key K, def ...int64) (val int64, ok bool) {
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

// GetToInt64 returns int64 value with type conversion for the given key
func (m Generic[K]) GetToInt64(key K, def ...int64) (val int64, ok bool, err error) {
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

// GetUInt64 returns uint64 value for the given key
func (m Generic[K]) GetUInt64(key K, def ...uint64) (val uint64, ok bool) {
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

// GetToUInt64 returns uint64 value with type conversion for the given key
func (m Generic[K]) GetToUInt64(key K, def ...uint64) (val uint64, ok bool, err error) {
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

// GetBool returns bool value for the given key
func (m Generic[K]) GetBool(key K, def ...bool) (val bool, ok bool) {
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

// GetToBool returns bool value with type conversion for the given key
func (m Generic[K]) GetToBool(key K, def ...bool) (val bool, ok bool, err error) {
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

// GetFloat32 returns float32 value for the given key
func (m Generic[K]) GetFloat32(key K, def ...float32) (val float32, ok bool) {
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

// GetToFloat32 returns float32 value with type conversion for the given key
func (m Generic[K]) GetToFloat32(key K, def ...float32) (val float32, ok bool, err error) {
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

// GetFloat64 returns float64 value for the given key
func (m Generic[K]) GetFloat64(key K, def ...float64) (val float64, ok bool) {
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

// GetToFloat64 returns float64 value with type conversion for the given key
func (m Generic[K]) GetToFloat64(key K, def ...float64) (val float64, ok bool, err error) {
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

// GetDuration returns time.Duration value for the given key
func (m Generic[K]) GetDuration(key K, def ...time.Duration) (val time.Duration, ok bool) {
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

// GetToDuration returns time.Duration value with type conversion for the given key
func (m Generic[K]) GetToDuration(key K, def ...time.Duration) (val time.Duration, ok bool, err error) {
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

// GetTime returns time.Time value for the given key
func (m Generic[K]) GetTime(key K, def ...time.Time) (val time.Time, ok bool) {
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

// GetToTime returns time.Time value with type conversion for the given key
func (m Generic[K]) GetToTime(key K, def ...time.Time) (val time.Time, ok bool, err error) {
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
