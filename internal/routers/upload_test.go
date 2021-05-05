package routers

import (
	"contrplatform/configs"
	"log"
	"path/filepath"
	"reflect"
	"testing"
)

func Test_getIdFilesName(t *testing.T) {
	path,_ := filepath.Abs(configs.UploadSavePath+"/61637")
	log.Println(path)
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
		{
			name: "exp1",
			args: args{
				 path,
			},
			want: []string{"AcceptsHalo3D.sol"},
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getIdFilesName(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getIdFilesName() = %v, want %v", got, tt.want)
			}
		})
	}
}
