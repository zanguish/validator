package validator

import (
	"encoding/json"
	"github.com/spf13/cast"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Fn func(op Option) error

var fnMap = map[string]Fn{
	"IsRequired":   IsRequired,
	"IsRequiredIf": IsRequiredIf,
	"IsInteger":    IsInteger,
	"IsFloat":      IsFloat,
	"IsBool":       IsBool,
	"IsJson":       IsJson,
	"IsLen":        IsLen,
	"IsMinLen":     IsMinLen,
	"IsMaxLen":     IsMaxLen,
	"IsMin":        IsMin,
	"IsMax":        IsMax,
	"IsBetWeen":    IsBetWeen,
	"IsIn":         IsIn,
	"IsNotIn":      IsNotIn,
	"IsContains":   IsContains,
	"IsMatch":      IsMatch,
	"IsSameKey":    IsSameKey,
	"IsDiffKey":    IsDiffKey,
}

func (v *Validator) RegRule(name string, fn Fn) {
	fnMap[name] = fn
}

func IsRequired(op Option) error {
	reflectValue := reflect.ValueOf(op.val)
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	switch reflectValue.Kind() {
	case reflect.String, reflect.Map, reflect.Array, reflect.Slice:
		if reflectValue.Len() != 0 {
			return nil
		}
		return ErrBuild(op)
	}
	if cast.ToString(op.val) != "" {
		return nil
	}
	return ErrBuild(op)
}

func IsRequiredIf(op Option) error {
	match := true
	for _, condition := range op.RequireIf {
		if op.kv[condition.Key] != condition.Val {
			match = false
		}
	}
	if match {
		return IsRequired(op)
	}
	return nil
}

func IsInteger(op Option) error {
	switch op.val.(type) {
	case int, int8, int16, int32, int64:
		return nil
	case uint, uint8, uint16, uint32, uint64:
		return nil
	}
	return ErrBuild(op)
}

func IsFloat(op Option) error {
	switch op.val.(type) {
	case float32, float64:
		return nil
	}
	return ErrBuild(op)
}

func IsBool(op Option) error {
	switch op.val.(type) {
	case bool:
		return nil
	}
	return ErrBuild(op)
}

func IsJson(op Option) error {
	value := cast.ToString(op.val)
	if json.Valid([]byte(value)) {
		return nil
	}
	return ErrBuild(op)
}

func IsLen(op Option) error {
	value := cast.ToString(op.val)
	valueLen := utf8.RuneCountInString(value)
	if op.Len == valueLen {
		return nil
	}
	return ErrBuild(op, op.Len)
}

func IsMinLen(op Option) error {
	value := cast.ToString(op.val)
	valueLen := utf8.RuneCountInString(value)
	if valueLen >= op.MinLen {
		return nil
	}
	return ErrBuild(op, op.MinLen)
}

func IsMaxLen(op Option) error {
	value := cast.ToString(op.val)
	valueLen := utf8.RuneCountInString(value)
	if valueLen <= op.MaxLen {
		return nil
	}
	return ErrBuild(op, op.MaxLen)
}

func IsMin(op Option) error {
	value := cast.ToFloat64(op.val)
	if value >= op.Min {
		return nil
	}
	return ErrBuild(op, op.Min)
}

func IsMax(op Option) error {
	value := cast.ToFloat64(op.val)
	if value <= op.Max {
		return nil
	}
	return ErrBuild(op, op.Max)
}

func IsBetWeen(op Option) error {
	value := cast.ToFloat64(op.val)
	if value >= op.Min && value <= op.Max {
		return nil
	}
	return ErrBuild(op, op.Min, op.Max)
}

func IsIn(op Option) error {
	valueStr := cast.ToString(op.val)
	for _, a := range op.In {
		if a == valueStr {
			return nil
		}
	}
	return ErrBuild(op)
}

func IsNotIn(op Option) error {
	valueStr := cast.ToString(op.val)
	for _, a := range op.In {
		if a == valueStr {
			return ErrBuild(op)
		}
	}
	return nil
}

func IsContains(op Option) error {
	if ok := strings.Contains(cast.ToString(op.val), op.SubStr); ok {
		return nil
	}
	return ErrBuild(op, op.SubStr)
}

func IsMatch(op Option) error {
	if ok := regexp.MustCompile(op.Regex).MatchString(cast.ToString(op.val)); ok {
		return nil
	}
	return ErrBuild(op)
}

func IsSameKey(op Option) error {
	if op.val == op.kv[op.SameKey] {
		return nil
	}
	return ErrBuild(op, op.SameKey)
}

func IsDiffKey(op Option) error {
	if op.val != op.kv[op.DiffKey] {
		return nil
	}
	return ErrBuild(op, op.DiffKey)
}
