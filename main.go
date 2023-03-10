package main

import (
	"os"

	"github.com/chickenzord/goksei/pkg/goksei"
	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Overload(".env")

	username := os.Getenv("GOKSEI_USERNAME")
	password := os.Getenv("GOKSEI_PASSWORD")

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"Account",
		"Symbol",
		"Amount",
		"Closing Price",
		"Current Value",
	})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 3, Align: text.AlignRight},
		{Number: 4, Align: text.AlignRight},
		{Number: 5, Align: text.AlignRight},
	})

	client := goksei.NewClient(username, password)

	equityBalance, err := client.GetShareBalances(goksei.EquityType)
	if err != nil {
		panic(err)
	}

	for _, sb := range equityBalance.Data {
		t.AppendRow([]interface{}{
			sb.Account,
			sb.Symbol(),
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
			sb.Amount,
			humanize.FormatFloat("#,###.##", sb.ClosingPrice),
			humanize.FormatFloat("#,###.##", sb.CurrentValue()),
		})
	}

	t.Render()
}
