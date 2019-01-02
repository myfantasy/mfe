package mfe

import (
	"testing"
)

var (
	tempJSON = `
{"hello":"world", 
"":"empty",
"date":"20180402 12:35:34",
"numb":43543.345343,
"long_numb":5555.55555555555544444444444444444333333333333333222222222222222221111111111111111100000000000001,
"null":null,
"array_empty":[],
"array":[5, {"some":"struct"}, "234", {}, {"some_big" : {"struct":4353453, "345" : "dffg"}}],
"struct":{"struct1":{"s21":{},"s22":{"array":[34543, "dfgfdg", {}]}}}
}
`
)

func Test_String(t *testing.T) {

	v, e := VariantNewFromJSON(tempJSON)
	if e != nil {
		t.Fatal(e.Error())
	}

	s := v.String()

	if len(s) < 10 {
		t.Fatal("something wrong")
	}
}

func Test_VariantNewFromJSON(t *testing.T) {

	_, e := VariantNewFromJSON(tempJSON)

	if e != nil {
		t.Fatal(e.Error())
	}
}

func Benchmark_String(b *testing.B) {

	v, _ := VariantNewFromJSON(tempJSON)

	for i := 0; i < b.N; i++ {
		_ = v.String()
	}

}

func Benchmark_VariantNewFromJSON(b *testing.B) {

	for i := 0; i < b.N; i++ {
		_, _ = VariantNewFromJSON(tempJSON)
	}

}
