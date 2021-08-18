package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getTagByName(t *testing.T) {
	type args struct {
		object interface{}
		index  int
		name   Keyword
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name: "1",
			args: args{
				object: struct {
					Name string `test:"123"`
				}{},
				index: 0,
				name:  "test",
			},
			want:  "123",
			want1: true,
		},
		{
			name: "2",
			args: args{
				object: struct {
					Name  string `test:"123"`
					Name2 string `test:"12344ada"`
				}{},
				index: 1,
				name:  "test",
			},
			want:  "12344ada",
			want1: true,
		},
		{
			name: "2",
			args: args{
				object: &struct {
					Name  string `test:"123"`
					Name2 string `test:"12344ada"`
				}{},
				index: 1,
				name:  "test",
			},
			want:  "12344ada",
			want1: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getTagByName(tt.args.object, tt.args.index, tt.args.name)
			if got != tt.want {
				t.Errorf("getTagByName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getTagByName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	assert.Panics(t, func() { getTagByName("123", 1, "123") }, "panic")
}
