// Code generated by "stringer -type=Rank"; DO NOT EDIT.

package deck

import "strconv"

const (
	_Rank_name_0 = "Joker"
	_Rank_name_1 = "AceTwoThreeFourFiveSixSevenEightNineTenJackQueenKing"
)

var (
	_Rank_index_1 = [...]uint8{0, 3, 6, 11, 15, 19, 22, 27, 32, 36, 39, 43, 48, 52}
)

func (i Rank) String() string {
	switch {
	case i == -1:
		return _Rank_name_0
	case 1 <= i && i <= 13:
		i -= 1
		return _Rank_name_1[_Rank_index_1[i]:_Rank_index_1[i+1]]
	default:
		return "Rank(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
