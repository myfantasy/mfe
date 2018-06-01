package mfe

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// Execute some query in db
func Execute(db *sql.DB, query string) (v Variant, e error) {
	svAll := make(SV, 0)

	r, e := db.Query(query)

	if e != nil {
		return VariantNewNull(), e
	}

	for r.NextResultSet() {
		svDataRes := make(SV, 0)

		cols, _ := r.Columns()
		ct, _ := r.ColumnTypes()

		vals := make([]interface{}, len(cols))

		for i := range vals {
			vals[i] = new(interface{})
		}

		for r.Next() {
			vm := VMap{}

			er := r.Scan(vals...)
			if er != nil {
				return VariantNewNull(), e
			}
			for i := range vals {
				vv := VariantNew(*(vals[i].(*interface{})))
				dtn := ct[i].DatabaseTypeName()
				if dtn == "NUMERIC" && !vv.IsNull() {
					vm[ct[i].Name()] = vv.ToDecimal()
				} else {
					vm[ct[i].Name()] = vv
				}

			}

			svDataRes = append(svDataRes, VariantNew(vm))
		}

		svAll = append(svAll, VariantNew(svDataRes))
	}

	return VariantNew(svAll), nil
}

// SF Преобразование переменной для запроса в PG
func SF(i interface{}) (s string) {
	if i == nil {
		return "null"
	}

	switch i.(type) {
	case Variant:
		v := i.(Variant)
		if v.IsNull() {
			return "null"
		}
		if InI(v.typeCode, decimalType, stringType, boolType, timeType) {
			return SF(v.value)
		}
		return SF(v.String())

	case bool:
		if i.(bool) {
			return "true"
		}
		return "false"
	case time.Time:
		return "'" + i.(time.Time).Format("20060102 150405.999999") + "'"
	case string:
		return "'" + strings.Replace(i.(string), "'", "''", -1) + "'"
	case float32:
		return fmt.Sprintf("%f", i)
	case float64:
		return fmt.Sprintf("%f", i)
	case decimal.Decimal:
		return i.(decimal.Decimal).String()
	case int, int64, int32, int16, int8:
		return fmt.Sprintf("%v", i)
	case []uint8:
		return string(i.([]uint8))
	}

	return SF(fmt.Sprintf("%v", i))
}

// SFMS Преобразование переменной для запроса в MS
func SFMS(i interface{}) (s string) {
	if i == nil {
		return "null"
	}

	switch i.(type) {
	case Variant:
		v := i.(Variant)
		if v.IsNull() {
			return "null"
		}
		if InI(v.typeCode, decimalType, stringType, boolType, timeType) {
			return SFMS(v.value)
		}
		return SFMS(v.String())

	case bool:
		if i.(bool) {
			return "1"
		}
		return "0"
	}
	return SF(i)
}
