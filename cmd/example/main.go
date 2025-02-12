package main

import (
	"os"

	"github.com/chickenzord/goksei"
	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Overload(".env")

	username := os.Getenv("GOKSEI_USERNAME")
	password := os.Getenv("GOKSEI_PASSWORD")
	plainPassword := os.Getenv("GOKSEI_PLAIN_PASSWORD") == "true"

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"Account",
		"Symbol",
		"Name",
		"Amount",
		"Closing Price",
		"Current Value",
	})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 4, Align: text.AlignRight},
		{Number: 5, Align: text.AlignRight},
		{Number: 6, Align: text.AlignRight},
	})

	authStore, err := goksei.NewFileAuthStore(".goksei-auth")
	if err != nil {
		panic(err)
	}

	client := goksei.NewClient(goksei.ClientOpts{
		Username:      username,
		Password:      password,
		PlainPassword: plainPassword,
		AuthStore:     authStore,
	})

	equityBalance, err := client.GetShareBalances(goksei.EquityType)
	if err != nil {
		panic(err)
	}

	for _, sb := range equityBalance.Data {
		t.AppendRow([]interface{}{
			sb.Account,
			sb.Symbol(),
			sb.Name(),
			sb.Amount,
			humanize.FormatFloat("#,###.##", sb.ClosingPrice),
			humanize.FormatFloat("#,###.##", sb.CurrentValue()),
		})
	}

	mutualFundBalance, err := client.GetShareBalances(goksei.MutualFundType)
	if err != nil {
		panic(err)
	}

	for _, sb := range mutualFundBalance.Data {
		t.AppendRow([]interface{}{
			sb.Account,
			sb.Symbol(),
			sb.Name(),
			sb.Amount,
			humanize.FormatFloat("#,###.##", sb.ClosingPrice),
			humanize.FormatFloat("#,###.##", sb.CurrentValue()),
		})
	}

	t.Render()
}
