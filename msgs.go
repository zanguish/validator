package validator

var defMsg = "The %v value is invalid"

type msg map[string]string

var enMap = msg{
	"IsRequired":   "The %v field is required",
	"IsRequiredIf": "The %v field is required",
	"IsInteger":    "The %v value must be an integer type",
	"IsFloat":      "The %v value must be an float type",
	"IsBool":       "The %v value must be an bool type",
	"IsJson":       "The %v value must be a valid JSON string",
	"IsLen":        "The %v value length must be %v",
	"IsMinLen":     "The %v value length must be equal or greater than %v",
	"IsMaxLen":     "The %v value length must be equal or lesser than %v",
	"IsMin":        "The %v value must be equal or greater than %v",
	"IsMax":        "The %v value must be equal or lesser than %v",
	"IsBetWeen":    "The %v value must be between %v and %v",
	"IsIn":         "The %v value is not in acceptable range",
	"IsNotIn":      "The %v value is not in acceptable range",
	"IsContains":   "The %v value must be contains {%v}",
	"IsMatch":      "The %v value is not match by regex",
	"IsSameKey":    "The %v value must be the same as field %v",
	"IsDiffKey":    "The %v value must be different from field %v",
}

var zhCnMap = msg{
	"IsRequired":   "%v 字段不能为空",
	"IsRequiredIf": "%v 字段不能为空",
	"IsInteger":    "%v 字段应当为整数",
	"IsFloat":      "%v 字段应当为浮点数",
	"IsBool":       "%v 字段应当为布尔值",
	"IsJson":       "%v 字段应当为JSON格式",
	"IsLen":        "%v 字段长度应当为 %v 个字符",
	"IsMinLen":     "%v 字段最小长度应当为 %v 个字符",
	"IsMaxLen":     "%v 字段最大长度应当为 %v 个字符",
	"IsMin":        "%v 字段最小值应当当为 %v",
	"IsMax":        "%v 字段最大值应当为 %v",
	"IsBetWeen":    "%v 字段值大小应当在 %v 到 %v之间",
	"IsIn":         "%v 字段取值范围不合理",
	"IsNotIn":      "%v 字段取值范围不合理",
	"IsContains":   "%v 字段值应当包含 %v",
	"IsMatch":      "%v 字段值不满足规则",
	"IsSameKey":    "%v 字段值应该和字段 %v 相同",
	"IsDiffKey":    "%v 字段值不应该和字段 %v 相同",
}

var langMsgMap = map[string]msg{
	"en":    enMap,
	"zh_CN": zhCnMap,
}
