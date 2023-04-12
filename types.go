package goksei

import (
	"math"
	"strings"
)

type PortfolioType string

func (t PortfolioType) Name() string {
	switch t {
	case EquityType:
		return "equity"
	case MutualFundType:
		return "mutual_fund"
	case CashType:
		return "cash"
	case BondType:
		return "bond"
	case OtherType:
		return "other"
	}

	return "unknown"
}

var (
	EquityType     PortfolioType = "EKUITAS"
	MutualFundType PortfolioType = "REKSADANA"
	CashType       PortfolioType = "KAS"
	BondType       PortfolioType = "OBLIGASI"
	OtherType      PortfolioType = "LAINNYA"
)

type PortfolioSummaryResponse struct {
	Total   float64                   `json:"summaryValue"`
	Details []PortfolioSummaryDetails `json:"summaryResponse"`
}

type PortfolioSummaryDetails struct {
	Type    string  `json:"type"`
	Amount  float64 `json:"summaryAmount"`
	Percent float64 `json:"percent"`
}

type CashBalance struct {
	ID            int     `json:"id"`
	AccountNumber string  `json:"rekening"`
	BankID        string  `json:"bank"`
	Currency      string  `json:"currCode"`
	Balance       float64 `json:"saldo"`
	BalanceIDR    float64 `json:"saldoIdr"`
	Status        int     `json:"status"`
}

func (c *CashBalance) CurrentBalance() float64 {
	return math.Max(c.Balance, c.BalanceIDR)
}

type CashBalanceResponse struct {
	Data []CashBalance `json:"data"`
}

type ShareBalanceResponse struct {
	Total float64        `json:"summaryValue"`
	Data  []ShareBalance `json:"data"`
}

func (r *ShareBalanceResponse) RemoveInvalidData() {
	// ref: https://stackoverflow.com/a/20551116
	i := 0

	for _, b := range r.Data {
		if b.Valid() {
			r.Data[i] = b
			i++
		}
	}

	r.Data = r.Data[:i]
}

type ShareBalance struct {
	Account      string  `json:"rekening"`   // Security account number. Example: "XL001CANE000000"
	FullName     string  `json:"efek"`       // Name of the asset. Example: "GOTO - GOTO GOJEK TOKOPEDIA Tbk"
	Participant  string  `json:"partisipan"` // Security or Asset Management name. Example: "MAHAKARYA ARTHA SEKURITAS, PT "
	BalanceType  string  `json:"tipeSaldo"`  // Example: "available"
	Currency     string  `json:"curr"`       // Example: "IDR"
	Amount       float64 `json:"jumlah"`     // units owned
	ClosingPrice float64 `json:"harga"`      // last closing price

	// hidden unknown/unused fields

	id           int     `json:"id"`             //nolint // not sure what is it used for
	tipe         string  `json:"tipe"`           //nolint // not sure what is it used for
	rate         string  `json:"rate"`           //nolint // not sure what is it used for
	currentValue float64 `json:"nilaiInvestasi"` //nolint // better calculate it client-side from Amount*ClosingPrice
}

func (c *ShareBalance) Valid() bool {
	return c.Account != "" && c.FullName != ""
}

func (c *ShareBalance) CurrentValue() float64 {
	return c.Amount * c.ClosingPrice
}

func (c *ShareBalance) Symbol() string {
	return strings.Split(c.FullName, " - ")[0]
}

func (c *ShareBalance) Name() string {
	return strings.Split(c.FullName, " - ")[1]
}

type LoginRequest struct {
	ID       string `json:"id"`
	AppType  string `json:"appType"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Validation string `json:"validation"`
}

type GlobalIdentityResponse struct {
	Code       string
	Status     string
	Identities []GlobalIdentity
}

type GlobalIdentity struct {
	LoginID  string `json:"idLogin"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	FullName string `json:"fullName"`

	InvestorID   string `json:"investorId"`
	InvestorName string `json:"sidName"`
	CitizenID    string `json:"nikId"`
	PassportID   string `json:"passportId"`
	TaxID        string `json:"npwp"`   // Indonesian Tax number (NPWP)
	CardID       string `json:"cardId"` // KSEI card ID
}
