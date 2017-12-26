package backtest

import (
	"math/rand"
	"time"
)

// StrategyHandler is a basic strategy interface
type StrategyHandler interface {
	CalculateSignal(DataEventHandler, DataHandler, PortfolioHandler) (SignalEvent, error)
}

type Strategy struct{}

func randInt() int {
	rand.Seed(time.Now().UnixNano())
	num := rand.Float32()
	if num < 0.2 {
		return 1
	} else if num < 0.4 {
		return 2
	}
	return 0
}

func (s *Strategy) CalculateSignal(de DataEventHandler, d DataHandler, p PortfolioHandler) (SignalEvent, error) {
	event := Event{Time: de.GetTime(), Symbol: de.GetSymbol()}
	signal := Signal{Event: event}
	switch randInt() {
	case 1:
		signal.SetDirection("buy")
		break
	case 2:
		signal.SetDirection("sell")
		break
	}
	return &signal, nil
}
