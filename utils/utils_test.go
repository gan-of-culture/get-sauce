package utils

import "testing"

func TestGetLastItem(t *testing.T) {
	tests := []struct {
		name string
		list []string
		want string
	}{
		{
			name: "String slice",
			list: []string{"1", "2", "3", "last item"},
			want: "last item",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := GetLastItemString(tt.list)

			if item != tt.want {
				t.Errorf("Got: %v - want: %v", item, tt.want)
			}
		})
	}
}
