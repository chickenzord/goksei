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
		wantCustodianBank string
		wantOk            bool
	}{
		{
			name: "jago",
			args: args{
				id: "JAGO1",
			},
			wantCustodianBank: "PT Bank Jago Tbk",
			wantOk:            true,
		},
		{
			name: "jago",
			args: args{
				id: "JAGO2",
			},
			wantCustodianBank: "PT Bank Jago Tbk",
			wantOk:            true,
		},
		{
			name: "bri",
			args: args{
				id: "BRI01",
			},
			wantCustodianBank: "Bank Rakyat Indonesia (Persero), PT",
			wantOk:            true,
		},
		{
			name: "bri",
			args: args{
				id: "BRI02",
			},
			wantCustodianBank: "Bank Rakyat Indonesia (Persero), PT",
			wantOk:            true,
		},
		{
			name: "monopoli",
			args: args{
				id: "MONO1",
			},
			wantCustodianBank: "",
			wantOk:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCustodianBank, gotOk := CustodianBankNameByID(tt.args.id)
			if !reflect.DeepEqual(gotCustodianBank, tt.wantCustodianBank) {
				t.Errorf("CustodianBankByID() gotCustodianBank = %v, want %v", gotCustodianBank, tt.wantCustodianBank)
			}

			if gotOk != tt.wantOk {
				t.Errorf("CustodianBankByID() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_stripNumberSuffix(t *testing.T) {
	type args struct {
		s string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{"NAME01"},
			want: "NAME",
		},
		{
			args: args{"02NAME01"},
			want: "02NAME",
		},
		{
			args: args{"NA1ME"},
			want: "NA1ME",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stripNumberSuffix(tt.args.s); got != tt.want {
				t.Errorf("stripNumberSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}
