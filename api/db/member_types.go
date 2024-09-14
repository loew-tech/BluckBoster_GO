package db

import "math"

var MemberTypes = map[string]int{
	MemberTypeBasic:   5,
	MemberTypeAdvance: 8,
	MemberTypePremium: math.MaxInt64,
}
