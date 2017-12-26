package backtest

import (
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

// PortfolioHandler is the combined interface building block for a portfolio.
// Events - fill, order
type PortfolioHandler interface {
	OnSignaler
	OnFiller
	// Investor - Intended to check if portfolio holds a given asset
	Updater
	Casher
	Valuer
	Reseter
}

// OnSignaler as an intercafe for the OnSignal method
type OnSignaler interface {
	OnSignal(SignalEvent, DataHandler) (*Order, error)
}

// OnFiller as an intercafe for the OnFill method
type OnFiller interface {
	OnFill(FillEvent, DataHandler) (*Fill, error)
}

// Casher handles basic portolio info
type Casher interface {
	SetInitialCash(float64)
	InitialCash() float64
	SetCash(float64)
	Cash() float64
}

// Valuer returns the values of the portfolio
type Valuer interface {
	Value() float64
	ViewHoldings()
}

// Updater handles the updating of the portfolio on data events
type Updater interface {
	Update(DataEventHandler)
}

// Portfolio represent a simple portfolio struct.
type Portfolio struct {
	initialCash  float64
	cash         float64
	holdings     map[string]position
	transactions []FillEvent
	// sizeManager  SizeHandler
	// riskManager  RiskHandler
}

// // SetSizeManager sets the size manager to be used with the portfolio
// func (p *Portfolio) SetSizeManager(size SizeHandler) {
// 	p.sizeManager = size
// }

// // SetRiskManager sets the risk manager to be used with the portfolio
// func (p *Portfolio) SetRiskManager(risk RiskHandler) {
// 	p.riskManager = risk
// }

// Reset the portfolio into a clean state with set initial cash.
func (p *Portfolio) Reset() {
	p.cash = 0
	// p.holdings = nil
	p.transactions = nil
}

// OnSignal handles an incomming signal event
func (p *Portfolio) OnSignal(signal SignalEvent, data DataHandler) (*Order, error) {
	orderType := "MKT" // default Market - may want to expand to other types later
	var limit float64

	if signal.GetDirection() == "" {
		return &Order{}, errors.New("No direction")
	}

	currQty := p.holdings[signal.GetSymbol()].qty
	currCash := p.Cash()
	currPrice := data.Latest(signal.GetSymbol()).LatestPrice()

	if signal.GetDirection() == "sell" && currQty <= 0.2 {
		return &Order{}, errors.New("No holdings to sell")
	}

	if signal.GetDirection() == "buy" && currCash <= 0.2*currPrice {
		return &Order{}, errors.New("Not enough cash to buy")
	}

	initialOrder := &Order{
		Event: Event{
			Time:   signal.GetTime(),
			Symbol: signal.GetSymbol(),
		},
		Direction: signal.GetDirection(),
		// Qty should be set by PositionSizer
		Qty:       0.2,
		OrderType: orderType,
		Limit:     limit,
	}

	// Last price for asset
	// latest := data.Latest(signal.GetSymbol())

	// sizedOrder, err := p.sizeManager.SizeOrder(initialOrder, latest, p)
	// if err != nil {
	// }

	// order, err := p.riskManager.EvaluateOrder(sizedOrder, latest, p.holdings)
	// if err != nil {
	// }

	return initialOrder, nil
}

// OnFill handles an incomming fill event
func (p *Portfolio) OnFill(fill FillEvent, data DataHandler) (*Fill, error) {
	// Check for nil map, else initialise the map
	if p.holdings == nil {
		p.holdings = make(map[string]position)
	}

	// check if portfolio has already a holding of the symbol from this fill
	if pos, ok := p.holdings[fill.GetSymbol()]; ok {
		// update existing Position
		pos.Update(fill)
		p.holdings[fill.GetSymbol()] = pos
	} else {
		// create new position
		pos := position{}
		pos.Create(fill)
		p.holdings[fill.GetSymbol()] = pos
	}

	// update cash
	if fill.GetDirection() == "BOT" {
		p.cash = p.cash - fill.NetValue()
	} else {
		// direction is "SLD"
		p.cash = p.cash + fill.NetValue()
	}

	// add fill to transactions
	p.transactions = append(p.transactions, fill)

	f := fill.(*Fill)
	return f, nil
}

// IsInvested checks if the portfolio has an open position on the given symbol
func (p Portfolio) IsInvested(symbol string) (pos position, ok bool) {
	pos, ok = p.holdings[symbol]
	if ok && (pos.qty != 0) {
		return pos, true
	}
	return pos, false
}

// // IsLong checks if the portfolio has an open long position on the given symbol
// func (p Portfolio) IsLong(symbol string) (pos position, ok bool) {
// 	pos, ok = p.holdings[symbol]
// 	if ok && (pos.qty > 0) {
// 		return pos, true
// 	}
// 	return pos, false
// }

// // IsShort checks if the portfolio has an open short position on the given symbol
// func (p Portfolio) IsShort(symbol string) (pos position, ok bool) {
// 	pos, ok = p.holdings[symbol]
// 	if ok && (pos.qty < 0) {
// 		return pos, true
// 	}
// 	return pos, false
// }

// Update updates the holding on a data event
func (p *Portfolio) Update(d DataEventHandler) {
	if pos, ok := p.IsInvested(d.GetSymbol()); ok {
		pos.UpdateValue(d)
		p.holdings[d.GetSymbol()] = pos
	}
}

// SetInitialCash sets the initial cash value of the portfolio
func (p *Portfolio) SetInitialCash(initial float64) {
	p.initialCash = initial
}

// InitialCash returns the initial cash value of the portfolio
func (p Portfolio) InitialCash() float64 {
	return p.initialCash
}

// SetCash sets the current cash value of the portfolio
func (p *Portfolio) SetCash(cash float64) {
	p.cash = cash
}

// Cash returns the current cash value of the portfolio
func (p Portfolio) Cash() float64 {
	return p.cash
}

// Value return the current total value of the portfolio
func (p Portfolio) Value() float64 {
	holdingValue := decimal.NewFromFloat(0)
	for _, pos := range p.holdings {
		marketValue := decimal.NewFromFloat(pos.marketValue)
		holdingValue = holdingValue.Add(marketValue)
	}

	cash := decimal.NewFromFloat(p.cash)
	value, _ := cash.Add(holdingValue).Round(4).Float64()
	return value
}

func (p Portfolio) ViewHoldings() {
	fmt.Println(p.holdings)
}
