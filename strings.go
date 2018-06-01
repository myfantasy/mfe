package mfe

import "reflect"

//StringGetNotEmpty -- выдаёт первый не пустой результат
func StringGetNotEmpty(s ...string) string {

	for i := 0; i < len(s); i++ {
		if s[i] != "" {
			return s[i]
		}
	}
	return ""
}

// InS Значения в списке (string)
func InS(b string, s ...string) bool {

	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return true
		}
	}
	return false
}

// InI Значения в списке (int)
func InI(b int, s ...int) bool {

	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return true
		}
	}
	return false
}

// In Значения в списке (любые типы)
func In(a interface{}, b ...interface{}) bool {
	for _, c := range b {
		if reflect.DeepEqual(a, c) {
			return true
		}
	}
	return false
}
