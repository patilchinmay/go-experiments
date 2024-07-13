package main

import (
	"reflect"
	"testing"
)

func Test_generateInts(t *testing.T) {
	tests := []struct {
		name string
		want []int
	}{
		{
			name: "should return 0 to 9",
			want: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultCh := generateInts()

			got := []int{}
			for v := range resultCh {
				got = append(got, v)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateInts() = %v, want %v", got, tt.want)
			}
		})
	}
}
