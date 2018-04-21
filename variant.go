package mfe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

const (
	vmapType = iota
	listOfVariant
	decimalType
	stringType
	boolType
)

// Variant - Тип который будет анонимным для простых типов
type Variant struct {
	typeCode int
	isNull   bool
	value    interface{}
	valueVM  VMap
	valueSV  SV
}

// VariantNewNull create new null Variant
func VariantNewNull() Variant {
	var v Variant
	v.isNull = true
	return v
}

// VariantNew новый экземпляр
func VariantNew(i interface{}) Variant {
	switch i.(type) {
	case bool:
		return Variant{typeCode: boolType, value: i}
	case VMap:
		return Variant{typeCode: vmapType, valueVM: i.(VMap)}
	case SV:
		return Variant{typeCode: listOfVariant, valueSV: i.(SV)}
	case string:
		return Variant{typeCode: stringType, value: i}
	case float32:
		return Variant{typeCode: decimalType, value: decimal.NewFromFloat(float64(i.(float32)))}
	case float64:
		return Variant{typeCode: decimalType, value: decimal.NewFromFloat(i.(float64))}
	case int, int64, int32, int16, int8:
		d, _ := decimal.NewFromString(fmt.Sprintf("%v", i))
		return Variant{typeCode: decimalType, value: d}
	case []uint8:
		return Variant{typeCode: stringType, value: string(i.([]uint8))}
	}

	return Variant{typeCode: stringType, value: fmt.Sprintf("%v", i)}
}

// IsNull check value is null value
func (v Variant) IsNull() bool {
	return v.isNull
}

// IsVM check value is VMap
func (v Variant) IsVM() bool {
	return v.typeCode == vmapType
}

// IsSV check value is SV
func (v Variant) IsSV() bool {
	return v.typeCode == listOfVariant
}

// IsDecimal check value is decimal
func (v Variant) IsDecimal() bool {
	return v.typeCode == decimalType
}

// IsBool check value is boolean
func (v Variant) IsBool() bool {
	return v.typeCode == boolType
}

// IsString check value is string
func (v Variant) IsString() bool {
	return v.typeCode == stringType
}

// Dec - Get Decimal (decimal.Decimal) from Variant
func (v Variant) Dec() decimal.Decimal {
	if v.typeCode == decimalType {
		return v.value.(decimal.Decimal)
	}
	d := decimal.New(0, 0)
	return d
}

// Bool - Get bool from Variant
func (v Variant) Bool() bool {
	if v.typeCode == decimalType {
		return v.value.(bool)
	}
	return false
}

// Str - Get String Value from Variant
func (v Variant) Str() string {
	if v.typeCode == stringType {
		return v.value.(string)
	}
	return ""
}

// String - Get string display Variant
func (v Variant) String() string {
	if v.isNull {
		return "nil"
	}
	if v.typeCode == decimalType {
		return v.value.(decimal.Decimal).String()
	}
	if v.typeCode == boolType {
		return strconv.FormatBool(v.value.(bool))
	}
	if v.typeCode == stringType {
		s, _ := json.Marshal(v.value.(string))
		return string(s)
	}
	if v.typeCode == listOfVariant {
		var b bytes.Buffer
		b.WriteString("[")
		for key, value := range v.valueSV {
			if key > 0 {
				b.WriteString(",")
			}
			b.WriteString(value.String())
		}
		b.WriteString("]")
		return b.String()
	}
	if v.typeCode == vmapType {
		var b bytes.Buffer
		b.WriteString("{")
		fr := false
		for key, value := range v.valueVM {
			if fr {
				b.WriteString(",")
			}
			fr = true
			b.WriteString("\"" + key + "\"")
			b.WriteString(":")
			b.WriteString(value.String())
		}
		b.WriteString("}")
		return b.String()
	}
	panic("Unknown type or value")
}

// UnmarshalJSON is unmarchal into Variant
func (v *Variant) UnmarshalJSON(b []byte) error {
	ut := 0
	ks := []byte("[")[0]
	fs := []byte("{")[0]
	kk := []byte("\"")[0]
	sp := []byte(" ")[0]
	ni := []byte("n")[0]
	tr := []byte("t")[0]
	fa := []byte("f")[0]
	for i := 0; i < len(b); i++ {
		if ks == b[i] {
			ut = 1
			break
		}
		if fs == b[i] {
			ut = 2
			break
		}
		if kk == b[i] {
			ut = 3
			break
		}
		if sp != b[i] {
			if strings.Index("0987654321.", fmt.Sprintf("%s", b[i:i+1])) != -1 {
				ut = 4
				break
			}
		}
		if ni == b[i] {
			ut = 5
			break
		}
		if tr == b[i] {
			ut = 6
			break
		}
		if fa == b[i] {
			ut = 7
			break
		}
	}
	if ut == 7 { // false
		v.value = false
		v.typeCode = boolType
	}
	if ut == 6 { // true
		v.value = true
		v.typeCode = boolType
	}
	if ut == 5 { // nil
		v.isNull = true
	}
	if ut == 4 { // decimal
		var d decimal.Decimal
		if err := json.Unmarshal(b, &d); err != nil {
			return err
		}
		v.value = d
		v.typeCode = decimalType
	}
	if ut == 3 { // string
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		v.value = s
		v.typeCode = stringType
	}
	if ut == 2 { // map
		v.valueVM = VMap{}
		var m map[string]json.RawMessage
		if err := json.Unmarshal(b, &m); err != nil {
			return err
		}
		for key, value := range m {
			var vv Variant
			if err := json.Unmarshal(value, &vv); err != nil {
				return err
			}
			v.valueVM[key] = vv
		}
		v.typeCode = vmapType
	}
	if ut == 1 { // array
		var a []json.RawMessage
		if err := json.Unmarshal(b, &a); err != nil {
			return err
		}
		v.valueSV = make([]Variant, len(a))

		for key, value := range a {
			var vv Variant
			if err := json.Unmarshal(value, &vv); err != nil {
				return err
			}
			v.valueSV[key] = vv
		}
		v.typeCode = listOfVariant
	}
	if ut == 0 {
		v.isNull = true
	}

	return nil
}

// GetElement is getting element from Variant.VMap by iererchy
func (v Variant) GetElement(name ...string) (vo Variant, isOk bool) {
	if len(name) == 0 {
		vo, isOk = v, true
		return
	}
	if len(name) >= 0 && !v.IsVM() {
		vo, isOk = Variant{isNull: true}, false
		return
	}
	if len(name) == 1 {
		vo, isOk = v.valueVM[name[0]]
		return
	}

	vo, isOk = v.valueVM[name[0]].GetElement(name[1:len(name)]...)
	return

}
