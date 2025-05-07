package data

import (
	"math"

	"blockbuster/api/constants"
)

var MemberTypes = map[string]int{
	constants.MEMBER_TYPE_BASIC:   5,
	constants.MEMBER_TYPE_ADVANCE: 8,
	constants.MEMBER_TYPE_PREMIUM: math.MaxInt64,
}
