package tester_pool

import "testing"

func TestTesterPool_GeneID(t1 *testing.T) {
	type fields struct {
		testerMap map[string]*Tester
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
		{	"exp1",
			fields{
				make(map[string]*Tester),
			},
			"",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TesterPool{
				testerMap: tt.fields.testerMap,
			}
			if got := t.GeneID(); got != tt.want {
				t1.Errorf("GeneID() = %v, want %v", got, tt.want)
			}
		})
	}
}
