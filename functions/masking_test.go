package functions

import (
	"reflect"
	"testing"
)

func TestMaskingDataFromBytes(t *testing.T) {
	type args struct {
		from                []byte
		sensitivePathFields []string
	}
	tests := []struct {
		name   string
		args   args
		wantTo []byte
	}{
		{
			name: "success masking first layer",
			args: args{
				from:                []byte("{\"msg\":\"hello\",\"sensitive\":\"hide\"}"),
				sensitivePathFields: []string{"sensitive"},
			},
			wantTo: []byte("{\"msg\":\"hello\",\"sensitive\":\"-MASKED-\"}"),
		},
		{
			name: "success masking first layer with nested child",
			args: args{
				from:                []byte("{\"msg\":\"hello\",\"sensitive\":{\"sensitive-1\":\"hide-1\",\"sensitive-2\":\"hide-2\"}}"),
				sensitivePathFields: []string{"sensitive"},
			},
			wantTo: []byte("{\"msg\":\"hello\",\"sensitive\":\"-MASKED-\"}"),
		},
		{
			name: "success masking second layer",
			args: args{
				from:                []byte("{\"msg\":\"hello\",\"sensitive\":{\"sensitive-1\":\"hide-1\",\"sensitive-2\":\"hide-2\"}}"),
				sensitivePathFields: []string{"sensitive.sensitive-1"},
			},
			wantTo: []byte("{\"msg\":\"hello\",\"sensitive\":{\"sensitive-1\":\"-MASKED-\",\"sensitive-2\":\"hide-2\"}}"),
		},
		{
			name: "success masking combination first layer & second layer",
			args: args{
				from:                []byte("{\"msg\":\"hello\",\"sensitive\":{\"sensitive-1\":\"hide-1\",\"sensitive-2\":\"hide-2\"}}"),
				sensitivePathFields: []string{"sensitive.sensitive-1", "msg"},
			},
			wantTo: []byte("{\"msg\":\"-MASKED-\",\"sensitive\":{\"sensitive-1\":\"-MASKED-\",\"sensitive-2\":\"hide-2\"}}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTo := MaskingDataFromBytes(tt.args.from, tt.args.sensitivePathFields); !reflect.DeepEqual(gotTo, tt.wantTo) {
				t.Errorf("MaskingDataFromBytes() = %v, want %v", gotTo, tt.wantTo)
			}
		})
	}
}
