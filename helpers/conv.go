package helpers

import "encoding/json"

// Conversion functions

func ToMap(v interface{}) map[string]interface{} {
	// This is called a marshal-unmarshal cycle
	// needed to convert struct => map[string]iface{}
	var m map[string]interface{}
	byt, _ := json.Marshal(v)
	json.Unmarshal(byt, &m)
	return m
}

func FromMap[S interface{}](m map[string]interface{}) *S {
	// This is called a marshal-unmarshal cycle;
	// needed to convert map[string]iface{} => struct
	var s S
	byt, _ := json.Marshal(s)
	json.Unmarshal(byt, &s)
	return &s
}
