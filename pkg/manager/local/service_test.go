package local

import (
	"reflect"
	"testing"
)

func TestService_walkPath(t *testing.T) {
	tests := []struct {
		name string
		args string
		want []string
	}{
		{
			args: "this.is.a.path",
			want: []string{"this", "is", "a", "path"},
		},
		{
			args: "this\\.is.a.path",
			want: []string{"this\\.is", "a", "path"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{}
			var sts []string
			for str := range s.walkPath(tt.args) {
				sts = append(sts, str)
			}
			if !reflect.DeepEqual(sts, tt.want) {
				t.Errorf("%s != %s", sts, tt.want)
			}
		})
	}
}
