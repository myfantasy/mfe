package mfe

import (
	"reflect"
	"time"
)

// InS value in slise (string)
func InS(b string, s ...string) bool {

	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return true
		}
	}
	return false
}

// InI value in slise (int)
func InI(b int, s ...int) bool {

	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return true
		}
	}
	return false
}

// CoalesceS returns first not "" value
func CoalesceS(s ...string) string {

	for i := 0; i < len(s); i++ {
		if s[i] != "" {
			return s[i]
		}
	}
	return ""
}

// CoalesceI returns first not 0 value
func CoalesceI(s ...int) int {

	for i := 0; i < len(s); i++ {
		if s[i] != 0 {
			return s[i]
		}
	}
	return 0
}

// CoalesceD returns first not 0 value
func CoalesceD(s ...time.Duration) time.Duration {

	for i := 0; i < len(s); i++ {
		if s[i] != 0 {
			return s[i]
		}
	}
	return 0
}

// In value in slise (любые типы)
func In(a interface{}, b ...interface{}) bool {
	for _, c := range b {
		if reflect.DeepEqual(a, c) {
			return true
		}
	}
	return false
}

// IifS return first valur if condition true (the ternary operator)
func IifS(b bool, s1 string, s2 string) string {
	for b {
		return s1
	}
	return s2
}

// IifV return first valur if condition true (the ternary operator
func IifV(b bool, s1 Variant, s2 Variant) Variant {
	for b {
		return s1
	}
	return s2
}

// Iif return first valur if condition true (the ternary operator
func Iif(b bool, s1 interface{}, s2 interface{}) interface{} {
	for b {
		return s1
	}
	return s2
}

// JoinS - Join Strings with separator
func JoinS(separator string, vals ...string) (res string) {
	for i, s := range vals {
		if i > 0 {
			res += separator
		}
		res += s
	}
	return res
}
