package Validator

import (
	"encoding/json"
	"testing"
)

//valid tag : set verification rules
//alias tag : set field alias
//match tag : set custom regular rules,use with IsMatch
type Person struct {
	Age     int    `valid:"IsIntBetWeen:10,20"`
	Name    string `valid:"IsRequired|IsContains:z" alias:"user_name"`
	Address string `valid:"IsValidAddress|IsRequired"`
	OpenId  string `valid:"IsStrSizeRange:20,32"`
	Json    string `valid:"IsJson"`
	Email   string `valid:"IsEmail"`
	Match   string `valid:"IsMatch" match:"\\d{4}\\w{2}"`
}

//Method name must start with Is
//Add custom validator func
func (t Person) IsValidAddress(a interface{}) bool {
	return true
}

//Custom error message
//Can override the default error message
func (t Person) Messages() map[string]string {
	return map[string]string{
		"IsValidAddress": "%s not a valid value",
	}
}

func TestValidator_Struct(t *testing.T) {
	one := Person{
		Age:     11,
		Name:    "zanguish",
		Address: "china",
		OpenId:  "qwertyuiopasdfghjklz",
		Json:    `[{"a":"b"}]`,
		Email:   "123@123.com",
		Match:   "1234d2",
	}
	v := New()

	//i18n,built in `en`, `zh_CN`
	v.WithLocale("en")

	//override the default error message
	v.WithMessages(map[string]string{})

	//validate Struct data
	v.Struct(one)
	if !v.IsOk {
		json, _ := json.MarshalIndent(v.Errors(), "", "	")
		t.Errorf("result: %s", json)
	}
}
