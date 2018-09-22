# mfe
my fantasy extension  
mfe is simple tools for simplify your code  
## linq
#### value in slise: `InS, InI, In`
b in s[]

#### the ternary operator `IifS, IifV, Iif`
return first valur if condition true 
## Variant
Is a simple type for parse json or get data from database
#### example
```golang
v, err := mfe.VariantNewFromJSON(string(body))
if err != nil {
    log.Fatalln(err)
}

its := v.GE("response", "GeoObjectCollection", "featureMember")
if its.Count() > 0 {
    vp := its.GI(0).GE("GeoObject", "Point", "pos")
...
    fmt.Println(vp)
```