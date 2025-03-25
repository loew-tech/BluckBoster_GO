package data

import "math"

var MemberTypes = map[string]int{
	MEMBER_TYPE_BASIC:   5,
	MEMBER_TYPE_ADVANCE: 8,
	MEMBER_TYPE_PREMIUM: math.MaxInt64,
}
