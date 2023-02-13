package goksei

import (
	"strings"
)

type PortfolioType string

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

type ShareBalanceResponse struct {
	Total float64        `json:"summaryValue"`
	Data  []ShareBalance `json:"data"`
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

	id           int     `json:"id"`             // not sure what is it used for
	tipe         string  `json:"tipe"`           // not sure what is it used for
	rate         string  `json:"rate"`           // not sure what is it used for
	currentValue float64 `json:"nilaiInvestasi"` // better calculate it client-side from Amount*ClosingPrice
}

func (c *ShareBalance) CurrentValue() float64 {
	return c.Amount * c.ClosingPrice
}

func (c *ShareBalance) Symbol() string {
	return strings.Split(c.FullName, " - ")[0]
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
