package mfe

import "reflect"

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

// IifS Если истина первое значение иначе второе
func IifS(b bool, s1 string, s2 string) string {
	for b {
		return s1
	}
	return s2
}

// IifV Если истина первое значение иначе второе
func IifV(b bool, s1 Variant, s2 Variant) Variant {
	for b {
		return s1
	}
	return s2
}

// Iif Если истина первое значение иначе второе
func Iif(b bool, s1 interface{}, s2 interface{}) interface{} {
	for b {
		return s1
	}
	return s2
}
