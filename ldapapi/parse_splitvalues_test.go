package ldapapi

import (
	"fmt"
	"slices"
	"testing"
)

func TestSplitAttributeValues(t *testing.T) {
	cases := []struct {
		in     string
		values []string
		want   bool
	}{
		{in: "", values: []string{}, want: true},
		{in: "", values: []string{" "}, want: false},
		{in: ",", values: []string{}, want: true},
		{in: ", ,", values: []string{}, want: true},
		{in: " , , Val1,Val 2 , Val3", values: []string{"Val1", "Val 2", "Val3"}, want: true},
		{in: ", Val1 ,Val 2 , Val3, ,", values: []string{"Val1", "Val 2", "Val3"}, want: true},
		{in: ", Val1, ,Val 2 , Val3, ,", values: []string{"Val1", "Val 2", "Val3"}, want: true},
		{in: "Val1,,,,,Val 2 , Val3, Val   4 , ,", values: []string{"Val1", "Val 2", "Val3", "Val   4"}, want: true},
	}

	for id, c := range cases {
		t.Run(
			fmt.Sprintf("Test %v", id),
			func(t *testing.T) {
				got := SplitAttributeValues(c.in)
				if slices.Equal(c.values, got) != c.want {
					t.Errorf("Splitting %q : expected '%v' for '%+q'=='%+q' (input == output)", c.in, c.want, c.values, got)
					return
				}
			})
	}
}
