package ldapapi

import "testing"

func TestGetFirstMissingInt(t *testing.T) {
	cases := []struct {
		name    string
		nums    []int
		minimum int
		want    int
	}{
		{name: "legacy first gap", nums: []int{1003, 1000, 1001}, want: 1002},
		{name: "legacy contiguous", nums: []int{1002, 1000, 1001}, want: 1003},
		{name: "legacy single ID", nums: []int{1000}, want: 1001},
		{name: "legacy duplicates", nums: []int{1000, 1000, 1001, 1003}, want: 1002},
		{name: "legacy zero ID", nums: []int{0, 1, 3}, want: 2},
		{name: "empty without minimum", want: 1},
		{name: "empty slice without minimum", nums: []int{}, want: 1},
		{name: "empty with minimum", minimum: 1000, want: 1000},
		{name: "minimum below existing IDs", nums: []int{1002, 1003}, minimum: 1000, want: 1000},
		{name: "minimum is occupied", nums: []int{1000, 1001, 1003}, minimum: 1000, want: 1002},
		{name: "ignore IDs below minimum", nums: []int{100, 101, 1000, 1002}, minimum: 1000, want: 1001},
		{name: "minimum above existing IDs", nums: []int{100, 101}, minimum: 1000, want: 1000},
		{name: "duplicates around minimum", nums: []int{999, 1000, 1000, 1001}, minimum: 1000, want: 1002},
		{name: "negative minimum preserves legacy range", nums: []int{100, 102}, minimum: -1, want: 101},
		{name: "empty with negative minimum", minimum: -1, want: 1},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := getFirstMissingInt(c.nums, c.minimum); got != c.want {
				t.Errorf("getFirstMissingInt(%v, %d) = %d, want %d", c.nums, c.minimum, got, c.want)
			}
		})
	}
}
