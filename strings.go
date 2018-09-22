package mfe

//StringGetNotEmpty return first not emplty result
func StringGetNotEmpty(s ...string) string {

	for i := 0; i < len(s); i++ {
		if s[i] != "" {
			return s[i]
		}
	}
	return ""
}
