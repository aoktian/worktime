package templatefuncs

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
