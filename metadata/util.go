package metadata

import (
	"strconv"
	"time"
)

func GetBool(md Metadata, key string) (v bool) {
	if md == nil {
		return
	}

	if !md.IsExists(key) {
		return
	}
	switch vv := md.Get(key).(type) {
	case bool:
		v = vv
	case int:
		v = vv != 0
	case string:
		v, _ = strconv.ParseBool(vv)
	}

	return
}

func GetInt(md Metadata, key string) (v int) {
	if md == nil {
		return
	}

	if !md.IsExists(key) {
		return
	}
	switch vv := md.Get(key).(type) {
	case bool:
		if vv {
			v = 1
		}
	case int:
		v = vv
	case string:
		v, _ = strconv.Atoi(vv)
	}

	return
}

func GetFloat(md Metadata, key string) (v float64) {
	if md == nil {
		return
	}

	if !md.IsExists(key) {
		return
	}

	switch vv := md.Get(key).(type) {
	case float64:
		v = vv
	case int:
		v = float64(vv)
	case string:
		v, _ = strconv.ParseFloat(vv, 64)
	}

	return
}

func GetDuration(md Metadata, key string) (v time.Duration) {
	if md == nil {
		return
	}

	if !md.IsExists(key) {
		return
	}

	switch vv := md.Get(key).(type) {
	case int:
		v = time.Duration(vv) * time.Second
	case string:
		v, _ = time.ParseDuration(vv)
		if v == 0 {
			n, _ := strconv.Atoi(vv)
			v = time.Duration(n) * time.Second
		}
	}
	return
}

func GetString(md Metadata, key string) (v string) {
	if md == nil {
		return
	}

	if !md.IsExists(key) {
		return
	}

	switch vv := md.Get(key).(type) {
	case string:
		v = vv
	case int:
		v = strconv.FormatInt(int64(vv), 10)
	case int64:
		v = strconv.FormatInt(vv, 10)
	case uint:
		v = strconv.FormatUint(uint64(vv), 10)
	case uint64:
		v = strconv.FormatUint(uint64(vv), 10)
	case bool:
		v = strconv.FormatBool(vv)
	case float32:
		v = strconv.FormatFloat(float64(vv), 'f', -1, 32)
	case float64:
		v = strconv.FormatFloat(float64(vv), 'f', -1, 64)
	}

	return
}

func GetStrings(md Metadata, keys ...string) (ss []string) {
	if md == nil {
		return
	}

	for _, key := range keys {
		if !md.IsExists(key) {
			continue
		}

		switch v := md.Get(key).(type) {
		case []string:
			ss = v
		case []any:
			for _, vv := range v {
				if s, ok := vv.(string); ok {
					ss = append(ss, s)
				}
			}
		}
		break
	}
	return
}
