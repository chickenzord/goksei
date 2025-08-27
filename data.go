package goksei

import (
	"embed"
	"encoding/csv"
	"regexp"
	"sort"
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

// MutualFund contains information about a mutual fund product
// including its code, name, type, and investment manager.
type MutualFund struct {
	Code              string
	ProductName       string
	FundType          string
	InvestmentManager string
}

// MutualFundByCode looks up a mutual fund by its code.
// Returns the mutual fund information and true if found, nil and false otherwise.
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

// MutualFunds returns all mutual fund data sorted by code.
// The data is loaded from embedded CSV files containing OJK and KSEI mutual fund information.
func MutualFunds() []MutualFund {
	if len(mutualFunds) == 0 {
		initializeMutualFunds()
	}

	result := []MutualFund{}

	for _, f := range mutualFunds {
		result = append(result, f)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Code < result[j].Code
	})

	return result
}

// CustodianBank contains information about a custodian bank
// including its code and full name.
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

// CustodianBanks returns all custodian bank data sorted by code.
// The data is loaded from embedded CSV files from the KSEI website.
func CustodianBanks() []CustodianBank {
	if len(custodianBanks) == 0 {
		initializeCustodianBanks()
	}

	result := []CustodianBank{}

	for _, bank := range custodianBanks {
		result = append(result, bank)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Code < result[j].Code
	})

	return result
}

// CustodianBankByCode looks up a custodian bank by its code.
// The function strips numeric suffixes from codes for matching.
// Returns the bank information and true if found, nil and false otherwise.
func CustodianBankByCode(code string) (custodianBank *CustodianBank, ok bool) {
	if len(custodianBanks) == 0 {
		initializeCustodianBanks()
	}

	bank, ok := custodianBanks[stripNumberSuffix(code)]
	if !ok {
		return nil, false
	}

	return &bank, true
}

// CustodianBankNameByCode looks up a custodian bank name by its code.
// The function strips numeric suffixes from codes for matching.
// Returns the bank name and true if found, empty string and false otherwise.
func CustodianBankNameByCode(code string) (name string, ok bool) {
	if len(custodianBanks) == 0 {
		initializeCustodianBanks()
	}

	bank, ok := custodianBanks[stripNumberSuffix(code)]
	if !ok {
		return "", false
	}

	return bank.Name, true
}

// CustodianBankNameByID returns bank name by ID
//
// Deprecated: use CustodianBankNameByCode instead
func CustodianBankNameByID(id string) (name string, ok bool) {
	return CustodianBankNameByCode(id)
}

func stripNumberSuffix(s string) string {
	return numberSuffix.ReplaceAllString(s, "")
}
