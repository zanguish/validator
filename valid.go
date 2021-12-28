package validator

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	defLang    = "en"
	defTagName = "v"
)

type Validator struct {
	ErrSlice []errItem
	IsOk     bool
	bail     bool
	lang     string
	tagName  string
}
type errItem struct {
	key string
	val string
}
type Rules struct {
	ScopeField string
	ScopeRules []ScopeRule
}

type ScopeRule struct {
	Rn string
	Op Option
}

type Option struct {
	key       string
	val       interface{}
	kv        map[string]interface{}
	rn        string
	lang      string
	Alias     string
	SubStr    string
	Len       int
	MinLen    int
	MaxLen    int
	Max       float64
	Min       float64
	Regex     string
	SameKey   string
	DiffKey   string
	In        []string
	NotIn     []string
	RequireIf []RequireIf
	MapStrAny []Rules
}

type RequireIf struct {
	Key string
	Val interface{}
}

func New() *Validator {
	return &Validator{
		IsOk:    true,
		lang:    defLang,
		tagName: defTagName,
	}
}

func (v *Validator) SetLang(lang string) {
	v.lang = lang
}

func (v *Validator) WithMsg(msg map[string]string) {
	for errKey, errMsg := range msg {
		langMsgMap[v.lang][errKey] = errMsg
	}
}

func (v *Validator) SetTagName(tagName string) {
	v.tagName = tagName
}

func Struct2Map(obj interface{}, TagName string) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag
		fieldName := t.Field(i).Name
		if fieldName[0] >= 'A' && fieldName[0] <= 'Z' {
			if field, ok := tag.Lookup(TagName); ok {
				fieldName = field
			}
			data[fieldName] = v.Field(i).Interface()
		}
	}
	return data
}

func GetMapValue(m map[string]interface{}, k string) interface{} {
	ks := strings.Split(strings.Trim(k, "."), ".")
	var s interface{}
	for _, k := range ks {
		if v, ok := m[k]; ok {
			if v, ok := v.(map[string]interface{}); ok {
				m = v
				s = v
			} else {
				s = m[k]
			}
		} else {
			s = nil
		}
	}
	return s
}

func (v *Validator) CheckMap(m map[string]interface{}, r []Rules) {
	for _, rule := range r {
		if v.bail && len(v.ErrSlice) == 1 {
			return
		}
		value := GetMapValue(m, rule.ScopeField)
		for _, s := range rule.ScopeRules {
			s.Op.kv = m
			s.Op.key = rule.ScopeField
			s.Op.val = value
			s.Op.rn = s.Rn
			s.Op.lang = v.lang
			if len(s.Op.Alias) > 0 {
				s.Op.key = s.Op.Alias
			}
			if fn, ok := fnMap[s.Rn]; ok {
				err := fn(s.Op)
				if err != nil {
					v.ErrSlice = append(v.ErrSlice, errItem{
						key: rule.ScopeField,
						val: err.Error(),
					})
					v.IsOk = false
					if v.bail {
						break
					}
				}
			}
		}
	}
}

func (v *Validator) CheckStruct(s interface{}, r []Rules) {
	m := Struct2Map(s, v.tagName)
	v.CheckMap(m, r)
}

func ErrBuild(op Option, args ...interface{}) error {
	if len(args) == 0 {
		args = append(args, op.key)
	}
	var format string
	langMsg := langMsgMap[op.lang]
	if msg, ok := langMsg[op.rn]; ok {
		format = msg
	} else {
		format = defMsg
	}
	vc := strings.Count(format, "%v")
	argLen := len(args)
	if argLen < vc {
		vc = argLen
	}
	args = args[:vc]
	return fmt.Errorf(format, args...)
}

func (v *Validator) FirstErr() string {
	for _, errItem := range v.ErrSlice {
		return errItem.val
	}
	return ""
}
