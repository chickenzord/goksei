package goksei

import (
	"reflect"
	"testing"
)

func TestMutualFundByCode(t *testing.T) {
	type args struct {
		code string
	}

	tests := []struct {
		name           string
		args           args
		wantMutualFund *MutualFund
		wantOk         bool
	}{
		{
			name:   "DH002FICDANPAS00",
			args:   args{code: "DH002FICDANPAS00"},
			wantOk: true,
			wantMutualFund: &MutualFund{
				Code:              "DH002FICDANPAS00",
				ProductName:       "Danamas Pasti",
				FundType:          "fixed_income_fund",
				InvestmentManager: "Sinarmas Asset Management, PT",
			},
		},
		{
			name:           "not_found",
			args:           args{code: "something-not-exists"},
			wantOk:         false,
			wantMutualFund: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMutualFund, gotOk := MutualFundByCode(tt.args.code)
			if !reflect.DeepEqual(gotMutualFund, tt.wantMutualFund) {
				t.Errorf("MutualFundByCode() gotMutualFund = %v, want %v", gotMutualFund, tt.wantMutualFund)
			}
			if gotOk != tt.wantOk {
				t.Errorf("MutualFundByCode() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
