package Validator

import (
	"fmt"
	reflect "reflect"
	"strconv"
	"strings"
)

type Validator struct {
	errorArr []string
	IsOk     bool
	i18n     string
}

func (t *Validator) Errors() []string {
	return t.errorArr
}

func (t *Validator) ErrorOne() string {
	for _, err := range t.errorArr {
		return err
	}
	return ""
}

func New() *Validator {
	validator := &Validator{
		errorArr: []string{},
		IsOk:     true,
		i18n:     "en",
	}
	return validator
}

func (t *Validator) WithLocale(defaultLocale string) {
	if len(languageMap[defaultLocale]) == 0 {
		panic("temporarily unsupported languages")
	}
	t.i18n = defaultLocale
}

func (t *Validator) WithMessages(messages map[string]string) {
	for ruleName, errMsg := range messages {
		languageMap[t.i18n][ruleName] = errMsg
	}
}

func (t *Validator) Struct(s interface{}) *Validator {
	messages := languageMap[t.i18n]

	sValue := reflect.ValueOf(s)
	sType := sValue.Type()

	validatorValue := reflect.ValueOf(t)
	validatorType := validatorValue.Type()

	var validatorFunc map[string]int
	validatorFunc = make(map[string]int)
	for i := 0; i < validatorType.NumMethod(); i++ {
		name := validatorType.Method(i).Name
		if strings.HasPrefix(name, "Is") {
			validatorFunc[name] = i
		}
	}

	var sFunc map[string]int
	sFunc = make(map[string]int)
	for i := 0; i < sType.NumMethod(); i++ {
		name := sType.Method(i).Name
		if strings.HasPrefix(name, "Is") {
			sFunc[name] = i
		}
		if name == "Messages" {
			rs := sValue.MethodByName(name).Call([]reflect.Value{})
			if d, ok := rs[0].Interface().(map[string]string); ok {
				for key, value := range d {
					messages[key] = value
				}
			}
		}
	}

	for i := 0; i < sType.NumField(); i++ {
		tag := sType.Field(i).Tag
		name := sType.Field(i).Name
		fieldName := name
		if field, ok := tag.Lookup("alias"); ok {
			fieldName = field
		}
		if validTag, ok := tag.Lookup("valid"); ok {
			validTag = strings.Trim(validTag, "|")
			validArr := strings.Split(validTag, "|")
			for _, item := range validArr {
				itemArr := strings.Split(item, ":")
				methodName := itemArr[0]

				handleObj := sValue
				if _, ok := sFunc[methodName]; !ok {
					if _, ok := validatorFunc[methodName]; ok {
						handleObj = validatorValue
					} else {
						panic(fmt.Sprintf("%s is not a valid tag", methodName))
					}
				}

				argLen := 0
				var argArr []string
				if len(itemArr) >= 2 {
					argStr := itemArr[1]
					argStr = strings.Trim(argStr, ",")
					argArr = strings.Split(argStr, ",")
					argLen = len(argArr)
				}

				if methodName == "IsMatch" {
					if matchStr, ok := tag.Lookup("match"); ok {
						argArr = append(argArr, matchStr)
						argLen = len(argArr)
					} else {
						panic(fmt.Sprintf("%s must set a match tag", methodName))
					}
				}

				numIn := handleObj.MethodByName(methodName).Type().NumIn()
				if argLen != numIn-1 {
					panic(fmt.Sprintf("%s incorrect number of function parameters", methodName))
				}

				param := make([]reflect.Value, argLen+1)
				for j := 0; j < numIn; j++ {
					if j == 0 {
						param[0] = sValue.Field(i)
					} else {
						inType := handleObj.MethodByName(methodName).Type().In(j).Kind()
						cv := argArr[j-1]
						switch inType {
						case reflect.Int:
							sd, _ := strconv.Atoi(cv)
							param[j] = reflect.ValueOf(sd)
						case reflect.Int8:
							sd, _ := strconv.ParseInt(cv, 10, 8)
							param[j] = reflect.ValueOf(sd)
						case reflect.Int16:
							sd, _ := strconv.ParseInt(cv, 10, 16)
							param[j] = reflect.ValueOf(sd)
						case reflect.Int32:
							sd, _ := strconv.ParseInt(cv, 10, 32)
							param[j] = reflect.ValueOf(sd)
						case reflect.Int64:
							sd, _ := strconv.ParseInt(cv, 10, 64)
							param[j] = reflect.ValueOf(sd)
						case reflect.Uint:
							sd, _ := strconv.ParseUint(cv, 10, 0)
							param[j] = reflect.ValueOf(sd)
						case reflect.Uint8:
							sd, _ := strconv.ParseUint(cv, 10, 8)
							param[j] = reflect.ValueOf(sd)
						case reflect.Uint16:
							sd, _ := strconv.ParseUint(cv, 10, 16)
							param[j] = reflect.ValueOf(sd)
							param[j] = reflect.ValueOf(sd)
						case reflect.Uint32:
							sd, _ := strconv.ParseUint(cv, 10, 32)
							param[j] = reflect.ValueOf(sd)
						case reflect.Uint64:
							sd, _ := strconv.ParseUint(cv, 10, 64)
							param[j] = reflect.ValueOf(sd)
						case reflect.Float32:
							sd, _ := strconv.ParseFloat(cv, 32)
							param[j] = reflect.ValueOf(sd)
						case reflect.Float64:
							sd, _ := strconv.ParseFloat(cv, 64)
							param[j] = reflect.ValueOf(sd)
						case reflect.Bool:
							sd, _ := strconv.ParseBool(cv)
							param[j] = reflect.ValueOf(sd)
						case reflect.String:
							sd := argArr[j-1]
							param[j] = reflect.ValueOf(sd)
						default:
							sd := argArr[j-1]
							param[j] = reflect.ValueOf(sd)
						}
					}

				}
				rs := handleObj.MethodByName(methodName).Call(param)
				cs := rs[0].Interface()
				if valid, ok := cs.(bool); ok && !valid {
					cString := messages[methodName]
					if len(cString) == 0 {
						cString = messages["_"]
					}
					fs := strings.Count(cString, "%s")
					if fs > 0 {
						fs := strings.Count(cString, "%s")
						var formats []interface{}
						formats = append(formats, fieldName)
						for j := 0; j < fs-1; j++ {
							formats = append(formats, argArr[j])
						}
						cString = fmt.Sprintf(cString, formats...)
					}
					t.errorArr = append(t.errorArr, cString)
					t.IsOk = false
				}

			}
		}
	}
	return t
}
