package Validator

import "github.com/zanguish/validator/i18n"

type mt map[string]string

var languageMap = map[string]mt{
	"zh_CN": i18n.ZhCNBuiltinMessages,
	"en":    i18n.EnBuiltinMessages,
}
