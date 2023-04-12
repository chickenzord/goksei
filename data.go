package goksei

import (
	"embed"
	"encoding/csv"
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

type CustodianBank struct {
	ID   string
	Name string
}

var (
	custodianBanks map[string]CustodianBank

	staticCustodianBanks = map[string]CustodianBank{
		"JAGO1": {
			ID: "JAGO1", Name: "PT Bank Jago Tbk",
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

	for id, bank := range staticCustodianBanks {
		custodianBanks[id] = bank
	}

	for _, row := range rows[1:] {
		custodianBanks[row[1]] = CustodianBank{
			ID:   row[1],
			Name: row[2],
		}
	}
}

func CustodianBankByID(id string) (custodianBank *CustodianBank, ok bool) {
	if len(custodianBanks) == 0 {
		initializeCustodianBanks()
	}

	m, ok := custodianBanks[id]
	if !ok {
		return nil, false
	}

	return &m, true
}
