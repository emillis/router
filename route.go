package veryFastRouter

//===========[STRUCTS]====================================================================================================

//route contains all the needed information to handle a single url path
type route struct {
	originalPattern string
	segments        []segment
	hasVariables    bool
	variablesCount  int
	methods         []string
	handler         HandlerFunc
}

//compare path supplied as *stringArray to the route and returns whether it matches
func (r *route) compare(path []string) (bool, []string) {
	if len(path) != len(r.segments) {
		return false, nil
	}

	variables := make([]string, 0, r.variablesCount*2)

	for i := 0; i < len(path); i++ {
		//If both values match, perfect! Both segments are the same - continue to check the rest.
		//Otherwise, proceed to further checks
		if path[i] == r.segments[i].original {
			continue
		}

		//If the segment that doesn't match is also not path variable - this route doesn't match!
		//However, if this segment doesn't match, but is path variable - add it to the variable array and continue
		if !r.segments[i].isVariable {
			return false, nil
		}

		//Assigning KEY
		variables = append(variables, r.segments[i].key)
		//Assigning VALUE
		variables = append(variables, path[i][1:])
	}

	//If it was path match, return true and an array of variables
	return true, variables
}
