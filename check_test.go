package validator

import (
	"testing"
)

var v *Validator
var rules = []Rules{
	{
		ScopeField: "nickname",
		ScopeRules: []ScopeRule{
			{Rn: "selfRule"},
		},
	},
	{
		ScopeField: "password",
		ScopeRules: []ScopeRule{
			{Rn: "IsLen", Op: Option{Len: 6}},
			{Rn: "IsMinLen", Op: Option{MinLen: 6}},
			{Rn: "IsMaxLen", Op: Option{MaxLen: 20}},
			{Rn: "IsMatch", Op: Option{Regex: "[a-z]{3}\\d{3}"}},
			{Rn: "IsContains", Op: Option{SubStr: "abc"}},
			{Rn: "IsSameKey", Op: Option{SameKey: "password2"}},
			{Rn: "IsDiffKey", Op: Option{DiffKey: "nickname"}},
		},
	},
	{
		ScopeField: "score",
		ScopeRules: []ScopeRule{
			{Rn: "IsIn", Op: Option{In: []string{"1", "2", "3.1"}}},
			{Rn: "IsNotIn", Op: Option{In: []string{"5.20"}}},
			{Rn: "IsMin", Op: Option{Min: 1}},
			{Rn: "IsMax", Op: Option{Max: 5}},
			{Rn: "IsBetWeen", Op: Option{Min: 1, Max: 5}},
			{Rn: "IsFloat"},
		},
	},
	{
		ScopeField: "require",
		ScopeRules: []ScopeRule{
			{Rn: "IsRequired"},
			{Rn: "IsInteger"},
		},
	},
	{
		ScopeField: "require_if",
		ScopeRules: []ScopeRule{
			{Rn: "IsRequiredIf", Op: Option{RequireIf: []RequireIf{
				{Key: "nickname", Val: "abc"},
				{Key: "password", Val: "abc123"}},
			}},
		},
	},
	{
		ScopeField: "status",
		ScopeRules: []ScopeRule{
			{Rn: "IsBool"},
		},
	},
	{
		ScopeField: "json",
		ScopeRules: []ScopeRule{
			{Rn: "IsJson", Op: Option{Alias: "json_alias"}},
		},
	},
	{
		ScopeField: "sub.a",
		ScopeRules: []ScopeRule{
			{Rn: "IsJson"},
		},
	},
	{
		ScopeField: "sub.b.c",
		ScopeRules: []ScopeRule{
			{Rn: "IsBool"},
		},
	},
}

//Custom rule
func selfRule(op Option) error {
	return nil
}

func init() {
	v = New()
	//Register custom rules or override default rules
	v.RegRule("selfRule", selfRule)

	//Set lang，The default is en
	//v.SetLang("zh_CN")

	//Set the tag name, the default tag name is v
	//Only used in the case of check Struct
	//v.SetTagName("v")

	//Customize the error message or override the default error message
	v.WithMsg(map[string]string{})
}
func TestCheckMap(t *testing.T) {
	testMap := make(map[string]interface{})
	testMap["nickname"] = "abc"
	testMap["password"] = "abc123"
	testMap["password2"] = "abc123"
	testMap["score"] = 3.1
	testMap["require"] = 0
	testMap["require_if"] = "he"
	testMap["status"] = false
	testMap["json"] = "{}"
	testMap["sub"] = map[string]interface{}{
		"a": "{}",
		"b": map[string]interface{}{
			"c": true,
		},
	}

	v.CheckMap(testMap, rules)

	t.Logf("%t\r\n", v.IsOk)
	if !v.IsOk {
		t.Fatal(v.FirstErr())
	}
}

type td struct {
	Nickname  string                 `v:"nickname"`
	Password  string                 `v:"password"`
	Password2 string                 `v:"password2"`
	Score     float64                `v:"score"`
	Require   int                    `v:"require"`
	RequireIf string                 `v:"require_if"`
	Status    bool                   `v:"status"`
	Json      string                 `v:"json"`
	Sub       map[string]interface{} `v:"sub"`
}

func TestCheckStruct(t *testing.T) {
	testStruct := td{
		Nickname:  "abc",
		Password:  "abc123",
		Password2: "abc123",
		Score:     3.1,
		Require:   0,
		RequireIf: "he",
		Status:    false,
		Json:      "{}",
		Sub: map[string]interface{}{
			"a": "{}",
			"b": map[string]interface{}{
				"c": false,
			},
		},
	}
	v.CheckStruct(testStruct, rules)

	t.Logf("%t\r\n", v.IsOk)
	if !v.IsOk {
		t.Fatal(v.FirstErr())
	}
}
