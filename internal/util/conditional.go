package util

func Iif(condition bool, then any, ele any) any {
	if condition {
		return then
	} else {
		return ele
	}
}

func GetOrDefault(m map[string]string, key string, defaultValue string) string {
	v, ok := m[key]
	if ok {
		return v
	} else {
		return defaultValue
	}
}
