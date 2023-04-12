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

var (
	numberSuffix = regexp.MustCompile(`[0-9]+$`)

	custodianBankNames map[string]string

	staticCustodianBankNames = map[string]string{
		"JAGO": "PT Bank Jago Tbk",
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

	custodianBankNames = make(map[string]string)

	for id, name := range staticCustodianBankNames {
		custodianBankNames[stripNumberSuffix(id)] = name
	}

	for _, row := range rows[1:] {
		custodianBankNames[stripNumberSuffix(row[1])] = row[2]
	}
}

func CustodianBankNameByID(id string) (name string, ok bool) {
	if len(custodianBankNames) == 0 {
		initializeCustodianBanks()
	}

	m, ok := custodianBankNames[stripNumberSuffix(id)]
	if !ok {
		return "", false
	}

	return m, true
}

func stripNumberSuffix(s string) string {
	return numberSuffix.ReplaceAllString(s, "")
}
