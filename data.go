package goksei

import (
	"embed"
	"encoding/csv"
	"regexp"
)

//go:embed data
var embedFS embed.FS

var mutualFunds map[string]MutualFund

func initializeMutualFunds() {
	f, err := embedFS.Open("data/mutualfunds.csv")
	if err != nil {
		panic(err)
	}

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		panic(err)
	}

	mutualFunds = make(map[string]MutualFund)

	for _, row := range rows {
		mutualFunds[row[0]] = MutualFund{
			Code:              row[0],
			ProductName:       row[1],
			InvestmentManager: row[2],
			FundType:          row[3],
		}
	}
}

type MutualFund struct {
	Code              string
	ProductName       string
	FundType          string
	InvestmentManager string
}

func MutualFundByCode(code string) (mutualFund *MutualFund, ok bool) {
	if len(mutualFunds) == 0 {
		initializeMutualFunds()
	}

	m, ok := mutualFunds[code]
	if !ok {
		return nil, false
	}

	return &m, true
}

func MutualFunds() []MutualFund {
	if len(mutualFunds) == 0 {
		initializeMutualFunds()
	}

	result := []MutualFund{}

	for _, f := range mutualFunds {
		result = append(result, f)
	}

	return result
}

type CustodianBank struct {
	Code string
	Name string
}

var (
	numberSuffix = regexp.MustCompile(`[0-9]+$`)

	custodianBanks map[string]CustodianBank

	staticCustodianBanks = []CustodianBank{
		{
			Code: "JAGO",
			Name: "PT Bank Jago Tbk",
		},
	}
)

func initializeCustodianBanks() {
	f, err := embedFS.Open("data/custodian_banks.csv")
	if err != nil {
		panic(err)
	}

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		panic(err)
	}

	custodianBanks = make(map[string]CustodianBank)

	for _, bank := range staticCustodianBanks {
		custodianBanks[bank.Code] = bank
	}

	for _, row := range rows[1:] {
		custodianBanks[stripNumberSuffix(row[1])] = CustodianBank{
			Code: row[1],
			Name: row[2],
		}
	}
}

func CustodianBanks() []CustodianBank {
	if len(custodianBanks) == 0 {
		initializeCustodianBanks()
	}

	result := []CustodianBank{}

	for _, bank := range custodianBanks {
		result = append(result, bank)
	}

	return result
}

func CustodianBankNameByID(id string) (name string, ok bool) {
	if len(custodianBanks) == 0 {
		initializeCustodianBanks()
	}

	bank, ok := custodianBanks[stripNumberSuffix(id)]
	if !ok {
		return "", false
	}

	return bank.Name, true
}

func stripNumberSuffix(s string) string {
	return numberSuffix.ReplaceAllString(s, "")
}
