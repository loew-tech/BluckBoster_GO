package data

import (
	"blockbuster/api/constants"

	"math"
)

var MemberTypes = map[string]int{
	constants.MEMBER_TYPE_BASIC:   5,
	constants.MEMBER_TYPE_ADVANCE: 8,
	constants.MEMBER_TYPE_PREMIUM: math.MaxInt64,
}
