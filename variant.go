package mfe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

const (
	vmapType = iota
	listOfVariant
	decimalType
	stringType
	boolType
	timeType
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

//TimeFromString -- Создаёт и парсит дату из строки
func TimeFromString(s string) Variant {

	s = strings.Replace(s, "Z", "", -1)
	s = strings.Replace(s, "\"", "", -1)

	lo := ""
	d := strings.Index(s, ":")

	sd := s[d:]

	ip := strings.Index(sd, "+") + d
	im := strings.Index(sd, "-") + d

	z := ip
	if z == -1 {
		z = im
	}

	if z != -1 {
		lo = lo + "Z07:00"
	}

	i := strings.Index(s, ".")
	if i != -1 {
		if z == -1 {
			z = len(s)
		}
		lo = "." + strings.Repeat("9", z-i-1) + lo
	}

	if d != -1 {
		d2 := strings.Index(sd[1:], ":")
		if d2 == 2 {
			lo = "15:04:05" + lo
		} else {
			lo = "15:04" + lo
		}
		if strings.Contains(s, "T") {
			lo = "T" + lo
		} else {
			lo = " " + lo
		}
	}
	if d > 2 {
		if strings.Contains(s, "-") {
			lo = "2006-01-02" + lo
		} else {
			lo = "20060102" + lo
		}
	}

	t, err := time.Parse(lo, s)

	if err != nil {
		panic(err)
	}
	return Variant{typeCode: timeType, value: t}
}

// ToDecimal -- Создаёт VariantDecimal из VariantString
func (v Variant) ToDecimal() Variant {
	vr, _ := decimal.NewFromString(strings.Replace(v.String(), "\"", "", -1))
	return VariantNew(vr)
}

// ToTime -- Создаёт VariantTime из VariantString
func (v Variant) ToTime() Variant {
	return TimeFromString(strings.Replace(v.String(), "\"", "", -1))
}

// VariantNew новый экземпляр
func VariantNew(i interface{}) Variant {
	if i == nil {
		return VariantNewNull()
	}

	switch i.(type) {
	case bool:
		return Variant{typeCode: boolType, value: i}
	case time.Time:
		return Variant{typeCode: timeType, value: i}
	case VMap:
		return Variant{typeCode: vmapType, valueVM: i.(VMap)}
	case SV:
		return Variant{typeCode: listOfVariant, valueSV: i.(SV)}
	case string:
		return Variant{typeCode: stringType, value: i}
	case float32:
		return Variant{typeCode: decimalType, value: decimal.NewFromFloat(float64(i.(float32)))}
	case decimal.Decimal:
		return Variant{typeCode: decimalType, value: i}
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

// IsTime check value is time
func (v Variant) IsTime() bool {
	return v.typeCode == timeType
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

// Time - Get time.Time from Variant
func (v Variant) Time() time.Time {
	if v.typeCode == timeType {
		return v.value.(time.Time)
	}
	return time.Time{}
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
		return "null"
	}
	if v.typeCode == decimalType {
		return v.value.(decimal.Decimal).String()
	}
	if v.typeCode == boolType {
		return strconv.FormatBool(v.value.(bool))
	}
	if v.typeCode == timeType {
		s, _ := json.Marshal(v.value.(time.Time))
		return string(s)
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

// MarshalJSON is marshal json from Variant
func (v Variant) MarshalJSON() ([]byte, error) {
	s := v.String()
	b := []byte(s)
	return b, nil
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
