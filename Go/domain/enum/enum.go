package enum

type SortOrder string

const (
	SortOrder_asc  = SortOrder("asc")
	SortOrder_desc = SortOrder("desc")
)

var validSortOrder = []SortOrder{SortOrder_asc, SortOrder_desc}

func (v SortOrder) String() string {
	return string(v)
}

func (v SortOrder) IsValid() bool {
	for _, item := range validSortOrder {
		if item == v {
			return true
		}
	}
	return false
}
