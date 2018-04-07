package mfe

const (
	vmapType       = 0
	listOfVmapType = 1
	longType       = 2
	stringType     = 3
)

// Variant - Тип который будет анонимным для простых типов
type Variant struct {
	typeName int
	isNull   bool
	value    interface{}
	valueVM  VMap
	valueSVM SVM
}

// IsNull check value is null value
func (v Variant) IsNull() bool {
	return v.isNull
}

// IsVM check value is VMap
func (v Variant) IsVM() bool {
	return v.typeName == vmapType
}

// IsSVM check value is SVM
func (v Variant) IsSVM() bool {
	return v.typeName == listOfVmapType
}

// IsInt64 check value is int64
func (v Variant) IsInt64() bool {
	return v.typeName == longType
}

// IsString check value is string
func (v Variant) IsString() bool {
	return v.typeName == longType
}

// Int64 - Get int64 (long) from Variant
func (v Variant) Int64() int64 {
	if v.typeName == longType {
		return v.value.(int64)
	}
	return 0
}

// String - Get string from Variant
func (v Variant) String() string {
	return v.value.(string)
}
