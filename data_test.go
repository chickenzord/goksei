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

func TestCustodianBankByID(t *testing.T) {
	type args struct {
		id string
	}

	tests := []struct {
		name              string
		args              args
		wantCustodianBank *CustodianBank
		wantOk            bool
	}{
		{
			name: "jago",
			args: args{
				id: "JAGO1",
			},
			wantCustodianBank: &CustodianBank{
				ID:   "JAGO1",
				Name: "PT Bank Jago Tbk",
			},
			wantOk: true,
		},
		{
			name: "bri",
			args: args{
				id: "BRI01",
			},
			wantCustodianBank: &CustodianBank{
				ID:   "BRI01",
				Name: "Bank Rakyat Indonesia (Persero), PT",
			},
			wantOk: true,
		},
		{
			name: "monopoli",
			args: args{
				id: "MONO1",
			},
			wantCustodianBank: nil,
			wantOk:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCustodianBank, gotOk := CustodianBankByID(tt.args.id)
			if !reflect.DeepEqual(gotCustodianBank, tt.wantCustodianBank) {
				t.Errorf("CustodianBankByID() gotCustodianBank = %v, want %v", gotCustodianBank, tt.wantCustodianBank)
			}
			if gotOk != tt.wantOk {
				t.Errorf("CustodianBankByID() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
