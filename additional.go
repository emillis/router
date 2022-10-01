package veryFastRouter

//===========[STRUCTS]====================================================================================================

//AdditionalInfo defines a struct which gets passed into each handler
type AdditionalInfo struct {
	Variables variables
}

//variables contain
type variables struct {
	data []string
}

//GetVar returns path segment if it was defined as a variable
func (v variables) GetVar(key string) string {
	for i := 0; i < len(v.data); i = i + 2 {
		if v.data[i] != key {
			continue
		}

		return v.data[i+1]
	}

	return ""
}

//===========[FUNCTIONALITY]====================================================================================================

//newAdditionalInfo returns pointer to new instance of AdditionalInfo
func newAdditionalInfo() *AdditionalInfo {
	return &AdditionalInfo{variables{[]string{}}}
}
