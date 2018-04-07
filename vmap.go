package mfe

// VMap is map of string -> variant
type VMap map[string]Variant

// SVM is slice of VMap
type SVM []VMap

// GetElement is getting element from VMap by iererchy
func (vm VMap) GetElement(name ...string) (v Variant, isOk bool) {
	if len(name) == 0 {
		v, isOk = Variant{isNull: true}, false
		return
	}
	if len(name) == 1 {
		v, isOk = vm[name[0]]
		return
	}

	v, isOk = vm.GetElement(name[1:len(name)]...)
	return

}
