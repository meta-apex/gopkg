package zlog

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/netip"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"
)

// Entry represents a zlog entry. It is instanced by one of the level method of Logger and finalized by the Msg or Msgf method.
type Entry struct {
	buf   []byte
	Level Level
	w     Writer
}

// Writer defines an entry writer interface.
type Writer interface {
	WriteEntry(*Entry) (int, error)
}

// Time append t formated as string using time.RFC3339Nano.
func (e *Entry) Time(key string, t time.Time) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '"')
	e.buf = t.AppendFormat(e.buf, "2006-01-02T15:04:05.999Z07:00")
	e.buf = append(e.buf, '"')
	return e
}

// TimeFormat append t formated as string using timefmt.
func (e *Entry) TimeFormat(key string, timefmt string, t time.Time) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	switch timefmt {
	case TimeFormatUnix:
		e.buf = strconv.AppendInt(e.buf, t.Unix(), 10)
	case TimeFormatUnixMs:
		e.buf = strconv.AppendInt(e.buf, t.UnixNano()/1000000, 10)
	case TimeFormatUnixWithMs:
		e.buf = strconv.AppendInt(e.buf, t.Unix(), 10)
		e.buf = append(e.buf, '.')
		e.buf = strconv.AppendInt(e.buf, t.UnixNano()/1000000%1000, 10)
	default:
		e.buf = append(e.buf, '"')
		e.buf = t.AppendFormat(e.buf, timefmt)
		e.buf = append(e.buf, '"')
	}
	return e
}

// Times append a formated as string array using time.RFC3339Nano.
func (e *Entry) Times(key string, a []time.Time) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '[')
	for i, t := range a {
		if i != 0 {
			e.buf = append(e.buf, ',')
		}
		e.buf = append(e.buf, '"')
		e.buf = t.AppendFormat(e.buf, time.RFC3339Nano)
		e.buf = append(e.buf, '"')
	}
	e.buf = append(e.buf, ']')

	return e
}

// TimesFormat append a formated as string array using timefmt.
func (e *Entry) TimesFormat(key string, timefmt string, a []time.Time) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '[')
	for i, t := range a {
		if i != 0 {
			e.buf = append(e.buf, ',')
		}
		switch timefmt {
		case TimeFormatUnix:
			e.buf = strconv.AppendInt(e.buf, t.Unix(), 10)
		case TimeFormatUnixMs:
			e.buf = strconv.AppendInt(e.buf, t.UnixNano()/1000000, 10)
		case TimeFormatUnixWithMs:
			e.buf = strconv.AppendInt(e.buf, t.Unix(), 10)
			e.buf = append(e.buf, '.')
			e.buf = strconv.AppendInt(e.buf, t.UnixNano()/1000000%1000, 10)
		default:
			e.buf = append(e.buf, '"')
			e.buf = t.AppendFormat(e.buf, timefmt)
			e.buf = append(e.buf, '"')
		}
	}
	e.buf = append(e.buf, ']')

	return e
}

// Bool append the val as a bool to the entry.
func (e *Entry) Bool(key string, b bool) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	e.buf = strconv.AppendBool(e.buf, b)
	return e
}

// Bools adds the field key with val as a []bool to the entry.
func (e *Entry) Bools(key string, b []bool) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '[')
	for i, a := range b {
		if i != 0 {
			e.buf = append(e.buf, ',')
		}
		e.buf = strconv.AppendBool(e.buf, a)
	}
	e.buf = append(e.buf, ']')
	return e
}

// Dur adds the field key with duration d to the entry.
func (e *Entry) Dur(key string, d time.Duration) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	if d < 0 {
		d = -d
		e.buf = append(e.buf, '-')
	}
	e.buf = strconv.AppendInt(e.buf, int64(d/time.Millisecond), 10)
	if n := (d % time.Millisecond); n != 0 {
		var tmp [7]byte
		b := n % 100 * 2
		n /= 100
		tmp[6] = smallsString[b+1]
		tmp[5] = smallsString[b]
		b = n % 100 * 2
		n /= 100
		tmp[4] = smallsString[b+1]
		tmp[3] = smallsString[b]
		b = n % 100 * 2
		tmp[2] = smallsString[b+1]
		tmp[1] = smallsString[b]
		tmp[0] = '.'
		e.buf = append(e.buf, tmp[:]...)
	}
	return e
}

// TimeDiff adds the field key with positive duration between time t and start.
// If time t is not greater than start, duration will be 0. Duration format follows the same principle as Dur().
func (e *Entry) TimeDiff(key string, t time.Time, start time.Time) *Entry {
	if e == nil {
		return nil
	}

	var d time.Duration
	if t.After(start) {
		d = t.Sub(start)
	}
	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	e.buf = strconv.AppendInt(e.buf, int64(d/time.Millisecond), 10)
	if n := (d % time.Millisecond); n != 0 {
		var tmp [7]byte
		b := n % 100 * 2
		n /= 100
		tmp[6] = smallsString[b+1]
		tmp[5] = smallsString[b]
		b = n % 100 * 2
		n /= 100
		tmp[4] = smallsString[b+1]
		tmp[3] = smallsString[b]
		b = n % 100 * 2
		tmp[2] = smallsString[b+1]
		tmp[1] = smallsString[b]
		tmp[0] = '.'
		e.buf = append(e.buf, tmp[:]...)
	}
	return e
}

// Durs adds the field key with val as a []time.Duration to the entry.
func (e *Entry) Durs(key string, d []time.Duration) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '[')
	for i, a := range d {
		if i != 0 {
			e.buf = append(e.buf, ',')
		}
		if a < 0 {
			a = -a
			e.buf = append(e.buf, '-')
		}
		e.buf = strconv.AppendInt(e.buf, int64(a/time.Millisecond), 10)
		if n := (a % time.Millisecond); n != 0 {
			var tmp [7]byte
			b := n % 100 * 2
			n /= 100
			tmp[6] = smallsString[b+1]
			tmp[5] = smallsString[b]
			b = n % 100 * 2
			n /= 100
			tmp[4] = smallsString[b+1]
			tmp[3] = smallsString[b]
			b = n % 100 * 2
			tmp[2] = smallsString[b+1]
			tmp[1] = smallsString[b]
			tmp[0] = '.'
			e.buf = append(e.buf, tmp[:]...)
		}
	}
	e.buf = append(e.buf, ']')
	return e
}

// Err adds the field "error" with serialized err to the entry.
func (e *Entry) Err(err error) *Entry {
	return e.AnErr("error", err)
}

// AnErr adds the field key with serialized err to the zlog context.
func (e *Entry) AnErr(key string, err error) *Entry {
	if e == nil {
		return nil
	}

	if err == nil {
		e.buf = append(e.buf, ',', '"')
		e.buf = append(e.buf, key...)
		e.buf = append(e.buf, "\":null"...)
		return e
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	if o, ok := err.(ObjectMarshaler); ok {
		o.MarshalObject(e)
	} else {
		e.buf = append(e.buf, '"')
		e.string(err.Error())
		e.buf = append(e.buf, '"')
	}
	return e
}

// Errs adds the field key with errs as an array of serialized errors to the entry.
func (e *Entry) Errs(key string, errs []error) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '[')
	for i, err := range errs {
		if i != 0 {
			e.buf = append(e.buf, ',')
		}
		if err == nil {
			e.buf = append(e.buf, "null"...)
		} else {
			e.buf = append(e.buf, '"')
			e.string(err.Error())
			e.buf = append(e.buf, '"')
		}
	}
	e.buf = append(e.buf, ']')
	return e
}

func appendFloat(b []byte, f float64, bits int) []byte {
	switch {
	case math.IsNaN(f):
		return append(b, `"NaN"`...)
	case math.IsInf(f, 1):
		return append(b, `"+Inf"`...)
	case math.IsInf(f, -1):
		return append(b, `"-Inf"`...)
	}
	abs := math.Abs(f)
	fmt := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if bits == 64 && (abs < 1e-6 || abs >= 1e21) || bits == 32 && (float32(abs) < 1e-6 || float32(abs) >= 1e21) {
			fmt = 'e'
		}
	}
	b = strconv.AppendFloat(b, f, fmt, -1, int(bits))
	if fmt == 'e' {
		// clean up e-09 to e-9
		n := len(b)
		if n >= 4 && b[n-4] == 'e' && b[n-3] == '-' && b[n-2] == '0' {
			b[n-2] = b[n-1]
			b = b[:n-1]
		}
	}
	return b
}

// Float64 adds the field key with f as a float64 to the entry.
func (e *Entry) Float64(key string, f float64) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	e.buf = appendFloat(e.buf, f, 64)
	return e
}

// Float32 adds the field key with f as a float32 to the entry.
func (e *Entry) Float32(key string, f float32) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	e.buf = appendFloat(e.buf, float64(f), 32)
	return e
}

// Floats64 adds the field key with f as a []float64 to the entry.
func (e *Entry) Floats64(key string, f []float64) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '[')
	for i, a := range f {
		if i != 0 {
			e.buf = append(e.buf, ',')
		}
		e.buf = appendFloat(e.buf, a, 64)
	}
	e.buf = append(e.buf, ']')
	return e
}

// Floats32 adds the field key with f as a []float32 to the entry.
func (e *Entry) Floats32(key string, f []float32) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '[')
	for i, a := range f {
		if i != 0 {
			e.buf = append(e.buf, ',')
		}
		e.buf = appendFloat(e.buf, float64(a), 32)
	}
	e.buf = append(e.buf, ']')
	return e
}

// Int64 adds the field key with i as a int64 to the entry.
func (e *Entry) Int64(key string, i int64) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	e.buf = strconv.AppendInt(e.buf, i, 10)
	return e
}

// Uint adds the field key with i as a uint to the entry.
func (e *Entry) Uint(key string, i uint) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	e.buf = strconv.AppendUint(e.buf, uint64(i), 10)
	return e
}

// Uint64 adds the field key with i as a uint64 to the entry.
func (e *Entry) Uint64(key string, i uint64) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	e.buf = strconv.AppendUint(e.buf, i, 10)
	return e
}

// Int adds the field key with i as a int to the entry.
func (e *Entry) Int(key string, i int) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	e.buf = strconv.AppendInt(e.buf, int64(i), 10)
	return e
}

// Int32 adds the field key with i as a int32 to the entry.
func (e *Entry) Int32(key string, i int32) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	e.buf = strconv.AppendInt(e.buf, int64(i), 10)
	return e
}

// Int16 adds the field key with i as a int16 to the entry.
func (e *Entry) Int16(key string, i int16) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	e.buf = strconv.AppendInt(e.buf, int64(i), 10)
	return e
}

// Int8 adds the field key with i as a int8 to the entry.
func (e *Entry) Int8(key string, i int8) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	e.buf = strconv.AppendInt(e.buf, int64(i), 10)
	return e
}

// Uint32 adds the field key with i as a uint32 to the entry.
func (e *Entry) Uint32(key string, i uint32) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	e.buf = strconv.AppendUint(e.buf, uint64(i), 10)
	return e
}

// Uint16 adds the field key with i as a uint16 to the entry.
func (e *Entry) Uint16(key string, i uint16) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	e.buf = strconv.AppendUint(e.buf, uint64(i), 10)
	return e
}

// Uint8 adds the field key with i as a uint8 to the entry.
func (e *Entry) Uint8(key string, i uint8) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	e.buf = strconv.AppendUint(e.buf, uint64(i), 10)
	return e
}

// Ints64 adds the field key with i as a []int64 to the entry.
func (e *Entry) Ints64(key string, a []int64) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '[')
	for i, n := range a {
		if i != 0 {
			e.buf = append(e.buf, ',')
		}
		e.buf = strconv.AppendInt(e.buf, n, 10)
	}
	e.buf = append(e.buf, ']')
	return e
}

// Ints32 adds the field key with i as a []int32 to the entry.
func (e *Entry) Ints32(key string, a []int32) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '[')
	for i, n := range a {
		if i != 0 {
			e.buf = append(e.buf, ',')
		}
		e.buf = strconv.AppendInt(e.buf, int64(n), 10)
	}
	e.buf = append(e.buf, ']')
	return e
}

// Ints16 adds the field key with i as a []int16 to the entry.
func (e *Entry) Ints16(key string, a []int16) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '[')
	for i, n := range a {
		if i != 0 {
			e.buf = append(e.buf, ',')
		}
		e.buf = strconv.AppendInt(e.buf, int64(n), 10)
	}
	e.buf = append(e.buf, ']')
	return e
}

// Ints8 adds the field key with i as a []int8 to the entry.
func (e *Entry) Ints8(key string, a []int8) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '[')
	for i, n := range a {
		if i != 0 {
			e.buf = append(e.buf, ',')
		}
		e.buf = strconv.AppendInt(e.buf, int64(n), 10)
	}
	e.buf = append(e.buf, ']')
	return e
}

// Ints adds the field key with i as a []int to the entry.
func (e *Entry) Ints(key string, a []int) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '[')
	for i, n := range a {
		if i != 0 {
			e.buf = append(e.buf, ',')
		}
		e.buf = strconv.AppendInt(e.buf, int64(n), 10)
	}
	e.buf = append(e.buf, ']')
	return e
}

// Uints64 adds the field key with i as a []uint64 to the entry.
func (e *Entry) Uints64(key string, a []uint64) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '[')
	for i, n := range a {
		if i != 0 {
			e.buf = append(e.buf, ',')
		}
		e.buf = strconv.AppendUint(e.buf, n, 10)
	}
	e.buf = append(e.buf, ']')
	return e
}

// Uints32 adds the field key with i as a []uint32 to the entry.
func (e *Entry) Uints32(key string, a []uint32) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '[')
	for i, n := range a {
		if i != 0 {
			e.buf = append(e.buf, ',')
		}
		e.buf = strconv.AppendUint(e.buf, uint64(n), 10)
	}
	e.buf = append(e.buf, ']')
	return e
}

// Uints16 adds the field key with i as a []uint16 to the entry.
func (e *Entry) Uints16(key string, a []uint16) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '[')
	for i, n := range a {
		if i != 0 {
			e.buf = append(e.buf, ',')
		}
		e.buf = strconv.AppendUint(e.buf, uint64(n), 10)
	}
	e.buf = append(e.buf, ']')
	return e
}

// Uints8 adds the field key with i as a []uint8 to the entry.
func (e *Entry) Uints8(key string, a []uint8) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '[')
	for i, n := range a {
		if i != 0 {
			e.buf = append(e.buf, ',')
		}
		e.buf = strconv.AppendUint(e.buf, uint64(n), 10)
	}
	e.buf = append(e.buf, ']')
	return e
}

// Uints adds the field key with i as a []uint to the entry.
func (e *Entry) Uints(key string, a []uint) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '[')
	for i, n := range a {
		if i != 0 {
			e.buf = append(e.buf, ',')
		}
		e.buf = strconv.AppendUint(e.buf, uint64(n), 10)
	}
	e.buf = append(e.buf, ']')
	return e
}

// RawJSON adds already encoded JSON to the zlog line under key.
func (e *Entry) RawJSON(key string, b []byte) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	e.buf = append(e.buf, b...)
	return e
}

// RawJSONStr adds already encoded JSON String to the zlog line under key.
func (e *Entry) RawJSONStr(key string, s string) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	e.buf = append(e.buf, s...)
	return e
}

// Str adds the field key with val as a string to the entry.
func (e *Entry) Str(key string, val string) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '"')
	e.string(val)
	e.buf = append(e.buf, '"')
	return e
}

// StrInt adds the field key with integer val as a string to the entry.
func (e *Entry) StrInt(key string, val int64) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '"')
	e.buf = strconv.AppendInt(e.buf, val, 10)
	e.buf = append(e.buf, '"')
	return e
}

// Stringer adds the field key with val.String() to the entry.
func (e *Entry) Stringer(key string, val fmt.Stringer) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	if val != nil {
		e.buf = append(e.buf, '"')
		e.string(val.String())
		e.buf = append(e.buf, '"')
	} else {
		e.buf = append(e.buf, "null"...)
	}
	return e
}

// GoStringer adds the field key with val.GoStringer() to the entry.
func (e *Entry) GoStringer(key string, val fmt.GoStringer) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	if val != nil {
		e.buf = append(e.buf, '"')
		e.string(val.GoString())
		e.buf = append(e.buf, '"')
	} else {
		e.buf = append(e.buf, "null"...)
	}
	return e
}

// Strs adds the field key with vals as a []string to the entry.
func (e *Entry) Strs(key string, vals []string) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '[')
	for i, val := range vals {
		if i != 0 {
			e.buf = append(e.buf, ',')
		}
		e.buf = append(e.buf, '"')
		e.string(val)
		e.buf = append(e.buf, '"')
	}
	e.buf = append(e.buf, ']')
	return e
}

// Byte adds the field key with val as a byte to the entry.
func (e *Entry) Byte(key string, val byte) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	switch val {
	case '"':
		e.buf = append(e.buf, "\"\\\"\""...)
	case '\\':
		e.buf = append(e.buf, "\"\\\\\""...)
	case '\n':
		e.buf = append(e.buf, "\"\\n\""...)
	case '\r':
		e.buf = append(e.buf, "\"\\r\""...)
	case '\t':
		e.buf = append(e.buf, "\"\\t\""...)
	case '\f':
		e.buf = append(e.buf, "\"\\u000c\""...)
	case '\b':
		e.buf = append(e.buf, "\"\\u0008\""...)
	case '<':
		e.buf = append(e.buf, "\"\\u003c\""...)
	case '\'':
		e.buf = append(e.buf, "\"\\u0027\""...)
	case 0:
		e.buf = append(e.buf, "\"\\u0000\""...)
	default:
		e.buf = append(e.buf, '"', val, '"')
	}
	return e
}

// Bytes adds the field key with val as a string to the entry.
func (e *Entry) Bytes(key string, val []byte) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '"')
	e.bytes(val)
	e.buf = append(e.buf, '"')
	return e
}

// BytesOrNil adds the field key with val as a string or nil to the entry.
func (e *Entry) BytesOrNil(key string, val []byte) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	if val == nil {
		e.buf = append(e.buf, "null"...)
	} else {
		e.buf = append(e.buf, '"')
		e.bytes(val)
		e.buf = append(e.buf, '"')
	}
	return e
}

const hex = "0123456789abcdef"

// Hex adds the field key with val as a hex string to the entry.
func (e *Entry) Hex(key string, val []byte) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '"')
	for _, v := range val {
		e.buf = append(e.buf, hex[v>>4], hex[v&0x0f])
	}
	e.buf = append(e.buf, '"')
	return e
}

// Encode encodes bytes using enc.AppendEncode to the entry.
func (e *Entry) Encode(key string, val []byte, enc interface {
	AppendEncode(dst, src []byte) []byte
}) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '"')
	e.buf = enc.AppendEncode(e.buf, val)
	e.buf = append(e.buf, '"')
	return e
}

// IPAddr adds IPv4 or IPv6 Address to the entry.
func (e *Entry) IPAddr(key string, ip net.IP) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '"')
	if ip4 := ip.To4(); ip4 != nil {
		_ = ip4[3]
		e.buf = strconv.AppendInt(e.buf, int64(ip4[0]), 10)
		e.buf = append(e.buf, '.')
		e.buf = strconv.AppendInt(e.buf, int64(ip4[1]), 10)
		e.buf = append(e.buf, '.')
		e.buf = strconv.AppendInt(e.buf, int64(ip4[2]), 10)
		e.buf = append(e.buf, '.')
		e.buf = strconv.AppendInt(e.buf, int64(ip4[3]), 10)
	} else if a, ok := netip.AddrFromSlice(ip); ok {
		e.buf = a.AppendTo(e.buf)
	}
	e.buf = append(e.buf, '"')
	return e
}

// IPPrefix adds IPv4 or IPv6 Prefix (address and mask) to the entry.
func (e *Entry) IPPrefix(key string, pfx net.IPNet) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '"')
	e.buf = append(e.buf, pfx.String()...)
	e.buf = append(e.buf, '"')
	return e
}

// MACAddr adds MAC address to the entry.
func (e *Entry) MACAddr(key string, ha net.HardwareAddr) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '"')
	for i, c := range ha {
		if i > 0 {
			e.buf = append(e.buf, ':')
		}
		e.buf = append(e.buf, hex[c>>4])
		e.buf = append(e.buf, hex[c&0xF])
	}
	e.buf = append(e.buf, '"')
	return e
}

// NetIPAddr adds IPv4 or IPv6 Address to the entry.
func (e *Entry) NetIPAddr(key string, ip netip.Addr) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '"')
	e.buf = ip.AppendTo(e.buf)
	e.buf = append(e.buf, '"')
	return e
}

// NetIPAddrs adds IPv4 or IPv6 Addresses to the entry.
func (e *Entry) NetIPAddrs(key string, ips []netip.Addr) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '[')
	for i, ip := range ips {
		if i > 0 {
			e.buf = append(e.buf, ',')
		}
		e.buf = append(e.buf, '"')
		e.buf = ip.AppendTo(e.buf)
		e.buf = append(e.buf, '"')
	}
	e.buf = append(e.buf, ']')
	return e
}

// NetIPAddrPort adds IPv4 or IPv6 with Port Address to the entry.
func (e *Entry) NetIPAddrPort(key string, ipPort netip.AddrPort) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '"')
	e.buf = ipPort.AppendTo(e.buf)
	e.buf = append(e.buf, '"')
	return e
}

// NetIPPrefix adds IPv4 or IPv6 Prefix (address and mask) to the entry.
func (e *Entry) NetIPPrefix(key string, pfx netip.Prefix) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '"')
	e.buf = pfx.AppendTo(e.buf)
	e.buf = append(e.buf, '"')
	return e
}

// Type adds type of the key using reflection to the entry.
func (e *Entry) Type(key string, v any) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '"')
	e.buf = append(e.buf, reflect.TypeOf(v).String()...)
	e.buf = append(e.buf, '"')
	return e
}

// Caller adds the file:line of the "caller" key.
// If depth is negative, adds the full /path/to/file:line of the "caller" key.
func (e *Entry) Caller(depth int) *Entry {
	if e == nil {
		return nil
	}

	var full bool
	var pc uintptr
	if depth < 0 {
		depth, full = -depth, true
	}
	e.caller(caller1(depth, &pc, 1, 1), pc, full)
	return e
}

// Stack enables stack trace printing for the error passed to Err().
func (e *Entry) Stack() *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ",\"stack\":\""...)
	e.bytes(stacks(false))
	e.buf = append(e.buf, '"')
	return e
}

// Enabled return false if the entry is going to be filtered out by zlog level.
func (e *Entry) Enabled() bool {
	return e != nil
}

// Discard disables the entry so Msg(f) won't print it.
func (e *Entry) Discard() *Entry {
	if e == nil {
		return e
	}

	if cap(e.buf) <= bbcap {
		epool.Put(e)
	}
	return nil
}

var notTest = true

// Msg sends the entry with msg added as the message field if not empty.
func (e *Entry) Msg(msg string) {
	if e == nil {
		return
	}

	if msg != "" {
		e.buf = append(e.buf, ",\"message\":\""...)
		e.string(msg)
		e.buf = append(e.buf, "\"}\n"...)
	} else {
		e.buf = append(e.buf, '}', '\n')
	}
	_, _ = e.w.WriteEntry(e)
	if (e.Level == FatalLevel) && notTest {
		os.Exit(255)
	}
	if (e.Level == PanicLevel) && notTest {
		panic(msg)
	}
	if cap(e.buf) <= bbcap {
		epool.Put(e)
	}
}

type bb struct {
	B []byte
}

func (b *bb) Write(p []byte) (int, error) {
	b.B = append(b.B, p...)
	return len(p), nil
}

var bbpool = sync.Pool{
	New: func() any {
		return new(bb)
	},
}

// Msgf sends the entry with formatted msg added as the message field if not empty.
func (e *Entry) Msgf(format string, v ...any) {
	if e == nil {
		return
	}

	b := bbpool.Get().(*bb)
	b.B = b.B[:0]
	e.buf = append(e.buf, ",\"message\":\""...)
	_, _ = fmt.Fprintf(b, format, v...)
	e.bytes(b.B)
	e.buf = append(e.buf, '"')
	if cap(b.B) <= bbcap {
		bbpool.Put(b)
	}
	e.Msg("")
}

// Msgs sends the entry with msgs added as the message field if not empty.
func (e *Entry) Msgs(args ...any) {
	if e == nil {
		return
	}

	b := bbpool.Get().(*bb)
	b.B = b.B[:0]
	e.buf = append(e.buf, ",\"message\":\""...)
	_, _ = fmt.Fprint(b, args...)
	e.bytes(b.B)
	e.buf = append(e.buf, '"')
	if cap(b.B) <= bbcap {
		bbpool.Put(b)
	}
	e.Msg("")
}

func (e *Entry) caller(n int, pc uintptr, fullpath bool) {
	if n < 1 {
		return
	}

	file, line, name := pcFileLineName(pc)
	if !fullpath {
		var i, j int
		for i = len(file) - 1; i >= 0; i-- {
			if file[i] == '/' {
				break
			}
		}
		if i > 0 {
			for j = i - 1; j >= 0; j-- {
				if file[j] == '/' {
					break
				}
			}
			if j > 0 {
				i = j
			}
			file = file[i+1:]
		}
		if i = strings.LastIndexByte(name, '/'); i > 0 {
			name = name[i+1:]
		}
	}

	e.buf = append(e.buf, ",\"caller\":\""...)
	e.buf = append(e.buf, file...)
	e.buf = append(e.buf, ':')
	e.buf = strconv.AppendInt(e.buf, int64(line), 10)
	e.buf = append(e.buf, "\",\"callerfunc\":\""...)
	e.buf = append(e.buf, name...)
	e.buf = append(e.buf, "\",\"goid\":"...)
	e.buf = strconv.AppendInt(e.buf, int64(goid()), 10)
}

var escapes = [256]bool{
	'"':  true,
	'<':  true,
	'\'': true,
	'\\': true,
	'\b': true,
	'\f': true,
	'\n': true,
	'\r': true,
	'\t': true,
}

func (e *Entry) escapeb(b []byte) {
	n := len(b)
	j := 0
	if n > 0 {
		// Hint the compiler to remove bounds checks in the loop below.
		_ = b[n-1]
	}
	for i := 0; i < n; i++ {
		switch b[i] {
		case '"':
			e.buf = append(e.buf, b[j:i]...)
			e.buf = append(e.buf, '\\', '"')
			j = i + 1
		case '\\':
			e.buf = append(e.buf, b[j:i]...)
			e.buf = append(e.buf, '\\', '\\')
			j = i + 1
		case '\n':
			e.buf = append(e.buf, b[j:i]...)
			e.buf = append(e.buf, '\\', 'n')
			j = i + 1
		case '\r':
			e.buf = append(e.buf, b[j:i]...)
			e.buf = append(e.buf, '\\', 'r')
			j = i + 1
		case '\t':
			e.buf = append(e.buf, b[j:i]...)
			e.buf = append(e.buf, '\\', 't')
			j = i + 1
		case '\f':
			e.buf = append(e.buf, b[j:i]...)
			e.buf = append(e.buf, '\\', 'u', '0', '0', '0', 'c')
			j = i + 1
		case '\b':
			e.buf = append(e.buf, b[j:i]...)
			e.buf = append(e.buf, '\\', 'u', '0', '0', '0', '8')
			j = i + 1
		case '<':
			e.buf = append(e.buf, b[j:i]...)
			e.buf = append(e.buf, '\\', 'u', '0', '0', '3', 'c')
			j = i + 1
		case '\'':
			e.buf = append(e.buf, b[j:i]...)
			e.buf = append(e.buf, '\\', 'u', '0', '0', '2', '7')
			j = i + 1
		case 0:
			e.buf = append(e.buf, b[j:i]...)
			e.buf = append(e.buf, '\\', 'u', '0', '0', '0', '0')
			j = i + 1
		}
	}
	e.buf = append(e.buf, b[j:]...)
}

func (e *Entry) escapes(s string) {
	n := len(s)
	j := 0
	if n > 0 {
		// Hint the compiler to remove bounds checks in the loop below.
		_ = s[n-1]
	}
	for i := 0; i < n; i++ {
		switch s[i] {
		case '"':
			e.buf = append(e.buf, s[j:i]...)
			e.buf = append(e.buf, '\\', '"')
			j = i + 1
		case '\\':
			e.buf = append(e.buf, s[j:i]...)
			e.buf = append(e.buf, '\\', '\\')
			j = i + 1
		case '\n':
			e.buf = append(e.buf, s[j:i]...)
			e.buf = append(e.buf, '\\', 'n')
			j = i + 1
		case '\r':
			e.buf = append(e.buf, s[j:i]...)
			e.buf = append(e.buf, '\\', 'r')
			j = i + 1
		case '\t':
			e.buf = append(e.buf, s[j:i]...)
			e.buf = append(e.buf, '\\', 't')
			j = i + 1
		case '\f':
			e.buf = append(e.buf, s[j:i]...)
			e.buf = append(e.buf, '\\', 'u', '0', '0', '0', 'c')
			j = i + 1
		case '\b':
			e.buf = append(e.buf, s[j:i]...)
			e.buf = append(e.buf, '\\', 'u', '0', '0', '0', '8')
			j = i + 1
		case '<':
			e.buf = append(e.buf, s[j:i]...)
			e.buf = append(e.buf, '\\', 'u', '0', '0', '3', 'c')
			j = i + 1
		case '\'':
			e.buf = append(e.buf, s[j:i]...)
			e.buf = append(e.buf, '\\', 'u', '0', '0', '2', '7')
			j = i + 1
		case 0:
			e.buf = append(e.buf, s[j:i]...)
			e.buf = append(e.buf, '\\', 'u', '0', '0', '0', '0')
			j = i + 1
		}
	}
	e.buf = append(e.buf, s[j:]...)
}

func (e *Entry) string(s string) {
	for _, c := range []byte(s) {
		if escapes[c] {
			e.escapes(s)
			return
		}
	}
	e.buf = append(e.buf, s...)
}

func (e *Entry) bytes(b []byte) {
	for _, c := range b {
		if escapes[c] {
			e.escapeb(b)
			return
		}
	}
	e.buf = append(e.buf, b...)
}

// Interface adds the field key with i marshaled using reflection.
func (e *Entry) Interface(key string, i any) *Entry {
	if e == nil {
		return nil
	}

	if o, ok := i.(ObjectMarshaler); ok {
		return e.Object(key, o)
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	b := bbpool.Get().(*bb)
	b.B = b.B[:0]
	enc := json.NewEncoder(b)
	enc.SetEscapeHTML(false)
	err := enc.Encode(i)
	if err != nil {
		b.B = b.B[:0]
		_, _ = fmt.Fprintf(b, `marshaling error: %+v`, err)
		e.buf = append(e.buf, '"')
		e.bytes(b.B)
		e.buf = append(e.buf, '"')
	} else {
		b.B = b.B[:len(b.B)-1]
		e.buf = append(e.buf, b.B...)
	}

	return e
}

// Object marshals an object that implement the ObjectMarshaler interface.
func (e *Entry) Object(key string, obj ObjectMarshaler) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':')
	if obj == nil || (*[2]uintptr)(unsafe.Pointer(&obj))[1] == 0 {
		e.buf = append(e.buf, "null"...)
		return e
	}

	n := len(e.buf)
	obj.MarshalObject(e)
	if n < len(e.buf) {
		e.buf[n] = '{'
		e.buf = append(e.buf, '}')
	} else {
		e.buf = append(e.buf, "null"...)
	}

	return e
}

// Objects marshals a slice of objects that implement the ObjectMarshaler interface.
func (e *Entry) Objects(key string, objects any) *Entry {
	if e == nil {
		return nil
	}

	values := reflect.ValueOf(objects)
	if values.Kind() != reflect.Slice {
		e.buf = append(e.buf, ',', '"')
		e.buf = append(e.buf, key...)
		e.buf = append(e.buf, `":null`...)
		return e
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '[')
	for i := 0; i < values.Len(); i++ {
		if i != 0 {
			e.buf = append(e.buf, ',')
		}
		value := values.Index(i)
		if value.Kind() == reflect.Ptr && value.IsNil() {
			e.buf = append(e.buf, "null"...)
		} else if obj, ok := value.Interface().(ObjectMarshaler); ok {
			i := len(e.buf)
			obj.MarshalObject(e)
			e.buf[i] = '{'
			e.buf = append(e.buf, '}')
		} else {
			e.buf = append(e.buf, `null`...)
		}
	}
	e.buf = append(e.buf, ']')
	return e
}

// Func allows an anonymous func to run only if the entry is enabled.
//
//go:nonline
func (e *Entry) Func(f func(e *Entry)) *Entry {
	if e != nil && f != nil {
		f(e)
	}
	return e
}

// EmbedObject marshals and Embeds an object that implement the ObjectMarshaler interface.
func (e *Entry) EmbedObject(obj ObjectMarshaler) *Entry {
	if e == nil {
		return nil
	}

	if obj != nil && (*[2]uintptr)(unsafe.Pointer(&obj))[1] != 0 {
		obj.MarshalObject(e)
	}
	return e
}

// Any adds the field key with f as an any value to the entry.
func (e *Entry) Any(key string, value any) *Entry {
	if e == nil {
		return nil
	}

	if value == nil || (*[2]uintptr)(unsafe.Pointer(&value))[1] == 0 {
		e.buf = append(e.buf, ',', '"')
		e.buf = append(e.buf, key...)
		e.buf = append(e.buf, '"', ':')
		e.buf = append(e.buf, "null"...)
		return e
	}
	switch value := value.(type) {
	case ObjectMarshaler:
		e.buf = append(e.buf, ',', '"')
		e.buf = append(e.buf, key...)
		e.buf = append(e.buf, '"', ':')
		value.MarshalObject(e)
	case Context:
		e.Dict(key, value)
	case []time.Duration:
		e.Durs(key, value)
	case time.Duration:
		e.Dur(key, value)
	case time.Time:
		e.Time(key, value)
	case net.HardwareAddr:
		e.MACAddr(key, value)
	case net.IP:
		e.IPAddr(key, value)
	case net.IPNet:
		e.IPPrefix(key, value)
	case json.RawMessage:
		e.buf = append(e.buf, ',', '"')
		e.buf = append(e.buf, key...)
		e.buf = append(e.buf, '"', ':')
		e.buf = append(e.buf, value...)
	case []bool:
		e.Bools(key, value)
	case []byte:
		e.Bytes(key, value)
	case []error:
		e.Errs(key, value)
	case []float32:
		e.Floats32(key, value)
	case []float64:
		e.Floats64(key, value)
	case []string:
		e.Strs(key, value)
	case string:
		e.Str(key, value)
	case bool:
		e.Bool(key, value)
	case error:
		e.AnErr(key, value)
	case float32:
		e.Float32(key, value)
	case float64:
		e.Float64(key, value)
	case int16:
		e.Int16(key, value)
	case int32:
		e.Int32(key, value)
	case int64:
		e.Int64(key, value)
	case int8:
		e.Int8(key, value)
	case int:
		e.Int(key, value)
	case uint16:
		e.Uint16(key, value)
	case uint32:
		e.Uint32(key, value)
	case uint64:
		e.Uint64(key, value)
	case uint8:
		e.Uint8(key, value)
	case fmt.GoStringer:
		e.GoStringer(key, value)
	case fmt.Stringer:
		e.Stringer(key, value)
	default:
		e.buf = append(e.buf, ',', '"')
		e.buf = append(e.buf, key...)
		e.buf = append(e.buf, '"', ':')
		b := bbpool.Get().(*bb)
		b.B = b.B[:0]
		enc := json.NewEncoder(b)
		enc.SetEscapeHTML(false)
		err := enc.Encode(value)
		if err != nil {
			b.B = b.B[:0]
			fmt.Fprintf(b, `marshaling error: %+v`, err)
			e.buf = append(e.buf, '"')
			e.bytes(b.B)
			e.buf = append(e.buf, '"')
		} else {
			b.B = b.B[:len(b.B)-1]
			e.buf = append(e.buf, b.B...)
		}
		if cap(b.B) <= bbcap {
			bbpool.Put(b)
		}
	}
	return e
}

// KeysAndValues sends keysAndValues to Entry
func (e *Entry) KeysAndValues(keysAndValues ...any) *Entry {
	if e == nil {
		return nil
	}

	var key string
	for i, v := range keysAndValues {
		if i%2 == 0 {
			key, _ = v.(string)
			continue
		}
		e.Any(key, v)
	}
	return e
}

// Fields type, used to pass to `Fields`.
type Fields map[string]any

// Fields is a helper function to use a map to set fields using type assertion.
func (e *Entry) Fields(fields Fields) *Entry {
	if e == nil {
		return nil
	}

	for key, value := range fields {
		e.Any(key, value)
	}
	return e
}

// Context represents contextual fields.
type Context []byte

// NewContext starts a new contextual entry.
func NewContext(dst []byte) (e *Entry) {
	e = new(Entry)
	e.buf = dst
	return
}

// Value builds the contextual fields.
func (e *Entry) Value() Context {
	if e == nil {
		return nil
	}
	return e.buf
}

// Context sends the contextual fields to entry.
func (e *Entry) Context(ctx Context) *Entry {
	if e == nil {
		return nil
	}

	if len(ctx) != 0 {
		e.buf = append(e.buf, ctx...)
	}
	return e
}

// Dict sends the contextual fields with key to entry.
func (e *Entry) Dict(key string, ctx Context) *Entry {
	if e == nil {
		return nil
	}

	e.buf = append(e.buf, ',', '"')
	e.buf = append(e.buf, key...)
	e.buf = append(e.buf, '"', ':', '{')
	if len(ctx) > 0 {
		e.buf = append(e.buf, ctx[1:]...)
	}
	e.buf = append(e.buf, '}')
	return e
}

// stacks is a wrapper for runtime.Stack that attempts to recover the data for all goroutines.
func stacks(all bool) (trace []byte) {
	// We don't know how big the traces are, so grow a few times if they don't fit. Start large, though.
	n := 10000
	if all {
		n = 100000
	}
	for i := 0; i < 5; i++ {
		trace = make([]byte, n)
		nbytes := runtime.Stack(trace, all)
		if nbytes < len(trace) {
			trace = trace[:nbytes]
			break
		}
		n *= 2
	}
	return
}

// wlprintf is a helper function for tests
func wlprintf(w Writer, level Level, format string, args ...any) (int, error) {
	return w.WriteEntry(&Entry{
		Level: level,
		buf:   []byte(fmt.Sprintf(format, args...)),
	})
}

func b2s(b []byte) string { return *(*string)(unsafe.Pointer(&b)) }

//go:noescape
//go:linkname now time.now
func now() (sec int64, nsec int32, mono int64)

//go:noescape
//go:linkname absDate time.absDate
func absDate(abs uint64, full bool) (year int, month time.Month, day int, yday int)

//go:noescape
//go:linkname absClock time.absClock
func absClock(abs uint64) (hour, min, sec int)

//go:noescape
//go:linkname caller1 runtime.callers
func caller1(skip int, pc *uintptr, len, cap int) int
