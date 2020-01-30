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

func TestCalcSizeInByte(t *testing.T) {
	tests := []struct {
		name   string
		number float64
		unit   string
		want   int64
	}{
		{
			name:   "Kilobytes to Bytes",
			number: 752,
			unit:   "KB",
			want:   752000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes := CalcSizeInByte(tt.number, tt.unit)

			if bytes != tt.want {
				t.Errorf("Got: %v - want: %v", bytes, tt.want)
			}
		})
	}
}
