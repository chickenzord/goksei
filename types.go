package goksei

import (
	"math"
	"strings"
)

// PortfolioType represents the type of portfolio asset (equity, mutual fund, cash, bond, or other).
type PortfolioType string

// Name returns a lowercase English name for the portfolio type.
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

// Predefined portfolio types used by the KSEI API.
var (
	// EquityType represents stock/equity portfolios.
	EquityType PortfolioType = "EKUITAS"

	// MutualFundType represents mutual fund portfolios.
	MutualFundType PortfolioType = "REKSADANA"

	// CashType represents cash balances.
	CashType PortfolioType = "KAS"

	// BondType represents bond portfolios.
	BondType PortfolioType = "OBLIGASI"

	// OtherType represents other types of portfolios.
	OtherType PortfolioType = "LAINNYA"
)

// PortfolioSummaryResponse represents the response from the portfolio summary API endpoint.
type PortfolioSummaryResponse struct {
	Total   float64                   `json:"summaryValue"`
	Details []PortfolioSummaryDetails `json:"summaryResponse"`
}

// PortfolioSummaryDetails contains detailed information about a specific portfolio type.
type PortfolioSummaryDetails struct {
	Type    string  `json:"type"`
	Amount  float64 `json:"summaryAmount"`
	Percent float64 `json:"percent"`
}

// CashBalance represents a cash balance in a specific account and currency.
type CashBalance struct {
	ID            int     `json:"id"`
	AccountNumber string  `json:"rekening"`
	BankID        string  `json:"bank"`
	Currency      string  `json:"currCode"`
	Balance       float64 `json:"saldo"`
	BalanceIDR    float64 `json:"saldoIdr"`
	Status        int     `json:"status"`
}

// CurrentBalance returns the current balance, choosing the maximum between Balance and BalanceIDR.
func (c *CashBalance) CurrentBalance() float64 {
	return math.Max(c.Balance, c.BalanceIDR)
}

// CashBalanceResponse represents the response from the cash balance API endpoint.
type CashBalanceResponse struct {
	Data []CashBalance `json:"data"`
}

// ShareBalanceResponse represents the response from the share balance API endpoint.
type ShareBalanceResponse struct {
	Total float64        `json:"summaryValue"`
	Data  []ShareBalance `json:"data"`
}

// RemoveInvalidData removes invalid entries from the share balance data in-place.
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

// ShareBalance represents a balance of shares/securities in a specific account.
type ShareBalance struct {
	Account      string  `json:"rekening"`   // Security account number. Example: "XL001CANE000000"
	FullName     string  `json:"efek"`       // Name of the asset. Example: "GOTO - GOTO GOJEK TOKOPEDIA Tbk"
	Participant  string  `json:"partisipan"` // Security or Asset Management name. Example: "MAHAKARYA ARTHA SEKURITAS, PT "
	BalanceType  string  `json:"tipeSaldo"`  // Example: "available"
	Currency     string  `json:"curr"`       // Example: "IDR"
	Amount       float64 `json:"jumlah"`     // units owned
	ClosingPrice float64 `json:"harga"`      // last closing price

	// hidden unknown/unused fields

	id             int     //nolint // not sure what is it used for
	tipe           string  //nolint // not sure what is it used for
	rate           string  //nolint // not sure what is it used for
	nilaiInvestasi float64 //nolint // better calculate it client-side from Amount*ClosingPrice
}

// Valid returns true if the share balance has required fields (Account and FullName).
func (c *ShareBalance) Valid() bool {
	return c.Account != "" && c.FullName != ""
}

// CurrentValue calculates the current market value by multiplying Amount by ClosingPrice.
func (c *ShareBalance) CurrentValue() float64 {
	return c.Amount * c.ClosingPrice
}

// Symbol extracts the security symbol from the FullName field (part before " - ").
func (c *ShareBalance) Symbol() string {
	return strings.Split(c.FullName, " - ")[0]
}

// Name extracts the security name from the FullName field (part after " - ").
func (c *ShareBalance) Name() string {
	return strings.Split(c.FullName, " - ")[1]
}

// LoginRequest represents the request payload for the login API endpoint.
type LoginRequest struct {
	ID       string `json:"id"`
	AppType  string `json:"appType"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the response from the login API endpoint.
type LoginResponse struct {
	Validation string `json:"validation"`
}

// GlobalIdentityResponse represents the response from the global identity API endpoint.
type GlobalIdentityResponse struct {
	Code       string
	Status     string
	Identities []GlobalIdentity
}

// GlobalIdentity contains detailed identity information for a user account.
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
