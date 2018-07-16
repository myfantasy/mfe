package mfe

//StringGetNotEmpty -- выдаёт первый не пустой результат
func StringGetNotEmpty(s ...string) string {

	for i := 0; i < len(s); i++ {
		if s[i] != "" {
			return s[i]
		}
	}
	return ""
}
