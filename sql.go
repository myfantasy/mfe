package mfe

import (
	"database/sql"
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
				vm[ct[i].Name()] = VariantNew(*(vals[i].(*interface{})))
			}

			svDataRes = append(svDataRes, VariantNew(vm))
		}

		svAll = append(svAll, VariantNew(svDataRes))
	}

	return VariantNew(svAll), nil
}
