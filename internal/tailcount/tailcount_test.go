package tailcount

import (
	"testing"
)

func Test_min(t *testing.T) {
	type args struct {
		s []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "min1", args: args{s: []int{1, 2, 3, 4, 5}}, want: 1},
		{name: "min2", args: args{s: []int{5, 4, 3, 2, 1}}, want: 1},
		{name: "min1", args: args{s: []int{-1, -2, -3, -4, -5}}, want: -5},
		{name: "min2", args: args{s: []int{-5, -4, -3, -2, -1}}, want: -5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := min(tt.args.s); got != tt.want {
				t.Errorf("min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sum(t *testing.T) {
	type args struct {
		s []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "sum1", args: args{s: []int{1, 2, 3, 4}}, want: 10},
		{name: "sum2", args: args{s: []int{1, 2, 3, 4, 5}}, want: 15},
		{name: "sum3", args: args{s: []int{-5, 1, 2, 3, 4, 5}}, want: 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sum(tt.args.s); got != tt.want {
				t.Errorf("sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_max(t *testing.T) {
	type args struct {
		s []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "max1", args: args{s: []int{1, 2, 3, 4, 5}}, want: 5},
		{name: "max2", args: args{s: []int{-1, -2, -3, -4, -5}}, want: -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := max(tt.args.s); got != tt.want {
				t.Errorf("max() = %v, want %v", got, tt.want)
			}
		})
	}
}
