package veryFastRouter

//===========[STRUCTS]====================================================================================================

//route contains all the needed information to handle a single url path
type route struct {
	originalPattern string
	segments        []segment
	hasVariables    bool
	methods         []Method
	handler         HandlerFunc
}

//compare path supplied as *pathDetails to the route and returns whether it matches
func (r *route) compare(pd *pathDetails) bool {
	if pd.count != len(r.segments) {
		return false
	}

	for i := 0; i < pd.count; i++ {
		if pd.segments[i] != r.segments[i].value && !r.segments[i].isVariable {
			return false
		}
	}

	return true
}
