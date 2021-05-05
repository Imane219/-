package tester_pool

import (
	"fmt"
	"testing"
)

func TestTesterPool_GetSfuzzOutputs(t1 *testing.T) {
	type fields struct {
		testerMap map[string]*Tester
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Output
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "exp",
			fields: fields{
				testerMap: map[string]*Tester{},
			},
			args: args{id: "123"},
			want: []*Output{
				{
					FileName:              "LuckyDoubler.sol",
					ContractName:          "LuckyDoubler",
					Duration:              "0 days, 0 hrs, 2 min, 0 sec                     ",
					Coverage:              "41%            ",
					Branches:              24,
					Predicates:            4,
					Tracebits:             10,
					GaslessSend:           true,
					ExceptionDisorder:     true,
					Reentrancy:            false,
					TimestampDependency:   false,
					BlockNumberDependency: false,
					DelegateCall:          false,
					FreezingEther:         false,
					IntegerOverFlow:       false,
					IntegerUnderFlow:      false,
				},
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TesterPool{
				testerMap: tt.fields.testerMap,
			}
			got, err := t.GetSfuzzOutputs(tt.args.id)
			if (err != nil) != tt.wantErr {
				t1.Errorf("GetResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("%+v",got[0])
		})
	}
}
