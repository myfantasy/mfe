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

// Variant - Is a simple type for parse json or get data from database
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

//TimeFromString create new VariantDate from string
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

	if z > d {
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

// ToDecimal create VariantDecimal from VariantString
func (v Variant) ToDecimal() Variant {
	vr, _ := decimal.NewFromString(strings.Replace(v.String(), "\"", "", -1))
	return VariantNew(vr)
}

// ToTime create VariantTime from VariantString
func (v Variant) ToTime() Variant {
	return TimeFromString(strings.Replace(v.String(), "\"", "", -1))
}

// IsSimpleValue return true if type is not Slise or VMap
func (v Variant) IsSimpleValue() bool {
	return InI(v.typeCode, decimalType, stringType, boolType, timeType)
}

// Value return value as Interface
func (v Variant) Value() interface{} {
	return v.Value
}

// VariantNew create new instance from interface{}
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

// VariantNewSV new item variant (SV type)
func VariantNewSV() Variant {
	return Variant{typeCode: listOfVariant, valueSV: make([]Variant, 0)}
}

// VariantNewVM new item variant (VMap type)
func VariantNewVM() Variant {
	return Variant{typeCode: vmapType, valueVM: VMap{}}
}

// VariantNewFromJSON create new variant from json
func VariantNewFromJSON(s string) (v Variant, e error) {
	v = Variant{}
	e = (&v).UnmarshalJSON([]byte(s))

	return v, e
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

// VMap - Get slice of Variant
func (v Variant) VMap() (vm VMap) {
	if v.typeCode == vmapType {
		return v.valueVM
	}
	return nil
}

// SV - Get slice of Variant
func (v Variant) SV() (sv SV) {
	if v.typeCode == listOfVariant {
		return v.valueSV
	}
	return nil
}

// Interface - Get interface from Variant
func (v Variant) Interface() (i interface{}) {
	if v.IsNull() {
		return nil
	}
	if v.typeCode == listOfVariant {
		return v.valueSV
	}
	if v.typeCode == vmapType {
		return v.valueVM
	}

	return v.value

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
func (v *Variant) GetElement(name ...string) (vo *Variant, isOk bool) {
	if len(name) == 0 {
		vo, isOk = v, true
		return
	}
	if len(name) >= 0 && !v.IsVM() {
		vr := VariantNewNull()
		vo, isOk = &vr, false
		return
	}
	if len(name) == 1 {
		ve, ok := v.valueVM[name[0]]
		if ok {
			return &ve, ok
		}
		vr := VariantNewNull()
		vo, isOk = &vr, false
		return
	}

	ve, ok := v.valueVM[name[0]]
	if ok {
		return (&ve).GetElement(name[1:len(name)]...)
	}
	vr := VariantNewNull()
	vo, isOk = &vr, false
	return
}

// GE try get element (GetElement and ignor error)
func (v *Variant) GE(name ...string) (vo *Variant) {
	vo, _ = v.GetElement(name...)
	return
}

// Count return count of element in Variant if it is SV else 0
func (v *Variant) Count() (c int) {
	if v.IsNull() {
		return 0
	}

	if v.IsSV() {
		sv := v.SV()
		if sv == nil {
			return 0
		}

		return len(sv)
	}
	return 0

}

// GI try get item
func (v *Variant) GI(i int) (vo *Variant) {
	if v.IsNull() {
		vr := VariantNewNull()
		return &vr
	}

	if v.IsSV() {
		sv := v.SV()
		if sv == nil || len(sv) <= i {
			vr := VariantNewNull()
			return &vr
		}

		return &sv[i]
	}
	vr := VariantNewNull()
	return &vr
}

// Add Add in Variant
func (v *Variant) Add(vi *Variant) (ok bool) {
	if v.IsSV() {
		v.valueSV = append(v.valueSV, *vi)
		v.isNull = false
		return true
	}

	return false
}

// AddOrUpdate Add or Update element in Variant
func (v *Variant) AddOrUpdate(vi *Variant, name ...string) (ok bool) {
	if len(name) >= 1 {
		if v.IsNull() || !v.IsVM() {
			v.typeCode = vmapType
			v.isNull = false
			v.valueVM = VMap{}
		}
	}
	if len(name) > 1 {
		vp, r := v.valueVM[name[0]]
		if !r {
			vp = VariantNewVM()
			v.valueVM[name[0]] = vp
		}
		return (&vp).AddOrUpdate(vi, name[1:len(name)]...)
	}
	if len(name) == 1 {
		v.valueVM[name[0]] = *vi
		return true
	}
	return false
}

// AddIfNotExists Add element in Variant if not exists
func (v *Variant) AddIfNotExists(vi *Variant, name ...string) (ok bool) {
	if len(name) >= 1 {
		if v.IsNull() || !v.IsVM() {
			v.typeCode = vmapType
			v.isNull = false
			v.valueVM = VMap{}
		}
	}
	if len(name) > 1 {
		vp, r := v.valueVM[name[0]]
		if !r {
			vp = VariantNewVM()
			v.valueVM[name[0]] = vp
		}
		return (&vp).AddOrUpdate(vi, name[1:len(name)]...)
	}
	if len(name) == 1 {
		_, r := v.valueVM[name[0]]
		if !r {
			v.valueVM[name[0]] = *vi
			return true
		}
	}
	return false
}

// Foreach - do f for each item in Variant if it is Slise ("" name) of Variant or Map (-1 index) of variant or do on it self if it is not (Slice or Map)
func (v *Variant) Foreach(f func(v *Variant, name string, index int)) {
	if v.IsSV() {
		if !v.IsNull() {
			for i, vl := range v.SV() {
				f(&vl, "", i)
			}
		}
		return
	}
	if v.IsVM() {
		if !v.IsNull() {
			for n, vl := range v.VMap() {
				f(&vl, n, -1)
			}
		}
		return
	}
	f(v, "", -1)
}

// Keys - for Map of Variant
func (v *Variant) Keys() (ks []string) {
	if v.IsVM() {
		if !v.IsNull() {
			for n := range v.VMap() {
				ks = append(ks, n)
			}
		}
	}

	return ks
}

// SplitBy Split Slise Varioant to Slise of []Variant
func (v *Variant) SplitBy(i int) (sv []Variant) {
	if v.IsSV() {
		if !v.IsNull() {
			b := i
			csv := VariantNewSV()
			for k, vi := range v.SV() {
				if b < k {
					b = b + i
					sv = append(sv, csv)
					csv = VariantNewSV()
				}
				csv.Add(&vi)
			}
			if csv.Count() > 0 {
				sv = append(sv, csv)
			}
		}
	}
	return sv
}
