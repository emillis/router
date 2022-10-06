package veryFastRouter

//===========[STRUCTS]====================================================================================================

//segment is a metadata for a single path segment
type segment struct {
	key        string //Stripped "/:" at the beginning the original
	original   string //Path segment as it's been passed in
	isVariable bool   //Is this segment a variable
	isMatchAll bool   //If pattern has /*segment (match all), this will be set to true
}

//===========[FUNCTIONALITY]====================================================================================================

//newSegment returns a new segment based on the string supplied
func newSegment(seg string) segment {
	s := segment{
		key:        seg[1:],
		original:   seg,
		isVariable: len(seg) > 1 && seg[1] == 58,
		isMatchAll: len(seg) > 1 && seg[1] == 42,
	}

	if s.isVariable || s.isMatchAll {
		s.key = seg[2:]
	}

	return s
}
