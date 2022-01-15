package container

import (
	"reflect"
	"strings"
	"testing"
)

func Test_configuration_ReadCommandArgs(t *testing.T) {
	type fields struct {
		args        []string
		commandArgs map[string]string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "ReadCommandArgs",
			fields: fields{
				args: strings.Split("- -x -? -D -Xkey1 -Xkey2=val2 -Dkey3 -Dkey4=val4 -- -Dkey5=val5", " "),
				commandArgs: map[string]string{
					"key3": "",
					"key4": "val4",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &configuration{commandArgs: make(map[string]string, 16)}
			c.ReadCommandArgs(tt.fields.args)
			if reflect.DeepEqual(c.commandArgs, tt.fields.commandArgs) {
				return
			}
			t.Errorf("Read Command Args: expected: %v, but: %v", tt.fields.commandArgs, c.commandArgs)
		})
	}
}
