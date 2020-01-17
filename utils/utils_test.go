package utils

import "testing"

func TestGetLastItem(t *testing.T) {
	tests := []struct {
		name string
		list interface{}
		want interface{}
	}{
		{
			name: "String slice",
			list: []string{"1", "2", "3", "last item"},
			want: "last item",
		}, {
			name: "Int slice",
			list: []int{1, 2, 3},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := GetLastItem(tt.list)

			if item == tt.want {
				t.Errorf("Got: %v - want: %v", item, tt.want)
			}
		})
	}
}
