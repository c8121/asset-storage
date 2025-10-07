package util

func Iif(condition bool, then any, ele any) any {
	if condition {
		return then
	} else {
		return ele
	}
}
