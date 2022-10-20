package veryFastRouter

import (
	"errors"
	"fmt"
	"strings"
)

//===========[STRUCTS]====================================================================================================

//route contains all the needed information to handle a single url path
type route struct {
	originalPattern string
	segments        []segment
	hasVariables    bool
	hasMatchAll     bool
	variablesCount  int
	methods         []string
	handler         HandlerFunc
}

//compare path supplied as *stringArray to the route and returns whether it matches
func (r *route) compare(path []string) (bool, []string) {
	if (len(path) < len(r.segments)) || (len(path) != len(r.segments) && !r.hasMatchAll) {
		return false, nil
	}

	variables := make([]string, 0, r.variablesCount*2)

	segLen, pathLen := len(r.segments)-1, len(path)-1
	for ; segLen >= 0; segLen, pathLen = segLen-1, pathLen-1 {
		//If both values match, perfect! Both segments are the same - continue to check the rest.
		//Otherwise, proceed to further checks
		if path[pathLen] == r.segments[segLen].original {
			continue
		}

		//If the segment that doesn't match is also not path variable - this route doesn't match!
		//However, if this segment doesn't match, but is path variable - add it to the variable array and continue
		if r.segments[segLen].isVariable {
			//Assigning KEY and VALUE
			variables = append(variables, r.segments[segLen].key, path[pathLen][1:])
			continue
		}

		if r.segments[segLen].isMatchAll {
			variables = append(variables, r.segments[segLen].key, strings.Join(path[pathLen:], ""))
			break
		}

		return false, nil
	}

	//If it was path match, return true and an array of variables
	return true, variables
}

//compareRoutes compares two routes and returns boolean based on weather the two are the same
func (r *route) compareRoutes(r2 *route) error {
	//Matching static patterns
	if r.originalPattern == r2.originalPattern {
		return errors.New(fmt.Sprintf("pattern \"%s\" has already been added to the HttpRouter", r2.originalPattern))
	}

	//Matching patterns with variables
	if r.hasVariables && r2.hasVariables && len(r.segments) == len(r2.segments) {
		for i, s1 := range r.segments {
			s2 := r2.segments[i]

			if s1.isVariable || s2.isVariable {
				continue
			}

			if s1.original == s2.original {
				continue
			}

			return nil
		}

		return errors.New(fmt.Sprintf("variable pattern \"%s\" conflicts with existing pattern \"%s\" in HttpRouter", r2.originalPattern, r.originalPattern))
	}

	//Matching "Match All" patterns
	if r.hasMatchAll && r2.hasMatchAll {
		length := len(r.segments)
		if len(r2.segments) > length {
			length = len(r2.segments)
		}

		//Looping in reverse
		for i, j := len(r.segments)-1, len(r2.segments)-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
			s1 := r.segments[i]
			s2 := r2.segments[j]

			if !s1.isMatchAll && !s2.isMatchAll && s1.original != s2.original {
				return nil
			}

			if s1.isMatchAll || s2.isMatchAll {
				return errors.New(fmt.Sprintf("conflicting \"Match All\" pattern detected between existing pattern \"%s\" and one being added \"%s\" (both patterns would match the same url path)", r.originalPattern, r2.originalPattern))
			}
		}

		return nil
	}

	return nil
}

//===========[FUNCTIONALITY]====================================================================================================

//newRoute returns pointer to a new route created from path supplied
func newRoute(path string) (*route, error) {
	path, err := fullPathCheck(path)
	if err != nil {
		return nil, err
	}

	r := route{
		originalPattern: path,
		segments:        splitIntoSegments(path),
	}

	for i, segment := range r.segments {
		if segment.isVariable {
			r.hasVariables = true
			r.variablesCount++
		}

		if segment.isMatchAll {
			if i != 0 {
				return nil, errors.New(fmt.Sprintf("\"Match All\" segment of the pattern can only be at the end of the pattern. Detected in segment %d out of %d of pattern \"%s\"", len(r.segments)-i, len(r.segments), r.originalPattern))
			}

			r.hasMatchAll = true
		}
	}

	return &r, nil
}
