package excel

import (
	"reflect"
	"testing"
)

func Test_findHeaders(t *testing.T) {
	type typedRow struct {
		a int    `excel:"a"`
		b string `excel:"b"`
	}

	type args struct {
		row TypedRow
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"test1", args{typedRow{}}, []string{"a", "b"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findHeaders(tt.args.row)
			if (err != nil) != tt.wantErr {
				t.Errorf("findHeaders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findHeaders() got = %v, want %v", got, tt.want)
			}
		})
	}
}
