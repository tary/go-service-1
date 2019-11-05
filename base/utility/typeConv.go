package utility

import (
	"bytes"
	"fmt"
	"strconv"

	log "github.com/cihub/seelog"
)

//TypeToString 几种类型转换为字符串
func TypeToString(v interface{}) string {
	var ret string

	switch paramType := v.(type) {
	case int32:
		ret = strconv.Itoa(int(paramType))
	case uint32:
		ret = strconv.Itoa(int(paramType))
	case int64:
		ret = strconv.Itoa(int(paramType))
	case uint64:
		ret = strconv.Itoa(int(paramType))
	case float32:
		ret = strconv.FormatFloat(float64(paramType), 'f', 6, 32)
	case float64:
		ret = strconv.FormatFloat(float64(paramType), 'f', 6, 64)
	case string:
		ret = paramType
	case bool:
		ret = strconv.FormatBool(paramType)
	case *int32:
		ret = strconv.Itoa(int(*paramType))
	case *uint32:
		ret = strconv.Itoa(int(*paramType))
	case *int64:
		ret = strconv.Itoa(int(*paramType))
	case *uint64:
		ret = strconv.Itoa(int(*paramType))
	case *float32:
		ret = strconv.FormatFloat(float64(*paramType), 'f', 6, 32)
	case *float64:
		ret = strconv.FormatFloat(float64(*paramType), 'f', 6, 64)
	case *string:
		ret = *paramType
	case *bool:
		ret = strconv.FormatBool(*paramType)
	default:
		panic(fmt.Errorf("typeToString unsupport type: %T", v))
	}

	return ret
}

// ArrayToStr 带有0的数组转换为字符串
func ArrayToStr(buff []byte) string {
	index := bytes.IndexByte(buff, 0)
	return string(buff[:index])
}

//Unquote 去掉双引号
func Unquote(str string) string {
	if len(str) > 1 && str[0] == '"' && str[len(str)-1] == '"' {
		return str[1 : len(str)-1]
	}

	return str
}

// Atof 字符串转float
func Atof(value string) float64 {
	if len(value) == 0 {
		return 0.0
	}

	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Error("value is not float: ", value, ", err: ", err)
		return 0.0
	}

	return v
}

// Atoi 字符串转int
func Atoi(value string) int {
	if len(value) == 0 {
		return 0
	}

	v, err := strconv.Atoi(value)
	if err != nil {
		log.Error("value is not int: ", value, ", err: ", err)
		return 0
	}

	return v
}

// Atob 字符串转bool
func Atob(value string) bool {
	if len(value) == 0 {
		return false
	}

	v, err := strconv.ParseBool(value)
	if err != nil {
		log.Error("value is not bool: ", value, ", err: ", err)
		return false
	}

	return v
}
