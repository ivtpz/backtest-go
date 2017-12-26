package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ivtpz/test-order-service"
)

func main() {
	test := backtest.New()

	symbols := []string{"USDT-ETH"}
	test.SetSymbols(symbols)

	data := backtest.Data{}
	data.Load("poloniex", "USDT-ETH", "12/10/2017 03:00:00 PM", "12/11/2017 03:00:00 PM")
	test.SetData(&data)

	portfolio := backtest.Portfolio{}
	portfolio.SetInitialCash(1000)
	test.SetPortfolio(&portfolio)

	strategy := backtest.Strategy{}
	test.SetStrategy(&strategy)

	exchange := backtest.Exchange{Symbol: "poloniex", ExchangeFee: 0, CommissionRate: 0.0025}
	test.SetExchange(&exchange)

	statistic := backtest.Statistic{}
	test.SetStatistic(&statistic)

	test.Run()

	statistic.PrintResult()
	fmt.Println(portfolio.Value())
	fmt.Println(statistic.SharpRatio(0.01))
	fmt.Println(statistic.SortinoRatio(0.01))

	http.HandleFunc("/", statistic.GraphResult)
	log.Fatal(http.ListenAndServe(":8088", nil))
}
