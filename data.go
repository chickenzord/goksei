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
