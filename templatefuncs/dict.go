package templatefuncs

import "reflect"

func dict(values ...interface{}) map[string]interface{} {
	if len(values)%2 != 0 {
		panic("无效的字典调用")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			panic("字典键必须是字符串")
		}
		dict[key] = values[i+1]
	}
	return dict
}

// InSlice 检查值是否存在于切片中（通用版本）
func InSlice(item interface{}, slice interface{}) bool {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return false
	}

	for i := 0; i < s.Len(); i++ {
		if reflect.DeepEqual(s.Index(i).Interface(), item) {
			return true
		}
	}
	return false
}

// InStringSlice 检查字符串是否存在于字符串切片中
func InStringSlice(item string, slice []string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// InIntSlice 检查整数是否存在于整数切片中
func InIntSlice(item int, slice []int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
