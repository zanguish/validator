package Validator

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	Email = `^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$`
	Url   = `^[a-zA-z]+://[^\s]*`
)

/*Global*/
func (t *Validator) IsRequired(val interface{}) bool {
	return !t.IsZeroVal(val)
}
func (t *Validator) IsZeroVal(val interface{}) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

/*String*/
func (t *Validator) IsContains(val string, s string) bool {
	return strings.Contains(val, s)
}
func (t *Validator) IsJson(val string) bool {
	if val == "" {
		return false
	}
	var js json.RawMessage
	return json.Unmarshal([]byte(val), &js) == nil
}

/*String Size */
func (t *Validator) IsStrSize(val string, wantLen int) bool {
	len := utf8.RuneCountInString(val)
	return wantLen == len
}
func (t *Validator) IsStrMinSize(val string, min int) bool {
	len := utf8.RuneCountInString(val)
	return min <= len
}
func (t *Validator) IsStrMaxSize(val string, max int) bool {
	len := utf8.RuneCountInString(val)
	return len <= max
}
func (t *Validator) IsStrSizeRange(val string, min int, max int) bool {
	len := utf8.RuneCountInString(val)
	return min <= len && len <= max
}

/*Number Range*/
func (t *Validator) IsIntEqual(val int, want int) bool {
	return val == want
}
func (t *Validator) IsIntNotEqual(val int, want int) bool {
	return val != want
}
func (t *Validator) IsIntMin(val int, dstVal int) bool {
	return val >= dstVal
}
func (t *Validator) IsIntMax(val int, dstVal int) bool {
	return val <= dstVal
}
func (t *Validator) IsIntBetWeen(val int, min int, max int) bool {
	return min <= val && val <= max
}

/*Regexp*/
func (t *Validator) IsEmail(val string) bool {
	return regexp.MustCompile(Email).MatchString(val)
}
func (t *Validator) IsUrl(val string) bool {
	return regexp.MustCompile(Url).MatchString(val)
}
func (t *Validator) IsMatch(val string, regular string) bool {
	return regexp.MustCompile(regular).MatchString(val)
}
