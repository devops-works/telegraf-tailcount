package main

import (
	"reflect"
	"testing"
)

func Test_getArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    *config
		wantErr bool
	}{
		{name: "no args", args: []string{}, want: nil, wantErr: true},
		{name: "wrong number of args", args: []string{"-i", "10", "-p", "1", "-m", "tailcount", "-t", "tag1=value1,tag2=value2"}, want: nil, wantErr: true},
		{name: "wrong tags", args: []string{"-t", "tag1=value1 tag2=value2", "file1"}, want: nil, wantErr: true},
		{name: "good number of args with tags", args: []string{"-i", "10", "-p", "1", "-m", "tailcount", "-t", "tag1=value1,tag2=value2", "file1"},
			want: &config{
				file:         "file1",
				interval:     10,
				peakInterval: 1,
				measurement:  "tailcount",
				tags:         "file=file1,tag1=value1,tag2=value2",
			},
			wantErr: false,
		},
		{name: "good number of args without tags", args: []string{"-i", "10", "-p", "1", "-m", "tailcount", "file1"},
			want: &config{
				file:         "file1",
				interval:     10,
				peakInterval: 1,
				measurement:  "tailcount",
				tags:         "file=file1",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getArgs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("getArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
