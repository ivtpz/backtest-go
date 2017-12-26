package backtest

import (
	"math"
)

// Implementing all orders as price takers
// Future enhancement: allow for market maker orders

// ExecutionHandler is the basic interface for executing orders
type ExecutionHandler interface {
	ExecuteOrder(OrderEvent, DataHandler) (*Fill, error)
}

// Exchange is a basic execution handler implementation
type Exchange struct {
	Symbol         string
	ExchangeFee    float64
	CommissionRate float64
}

// ExecuteOrder executes an order event
func (e *Exchange) ExecuteOrder(order OrderEvent, data DataHandler) (*Fill, error) {
	// fetch latest known data event for the symbol
	latest := data.Latest(order.GetSymbol())
	// simple implementation, creates a direct fill from the order
	// based on the last known data price
	f := &Fill{
		Event:    Event{Time: order.GetTime(), Symbol: order.GetSymbol()},
		Exchange: e.Symbol,
		Qty:      order.GetQty(),
		Price:    latest.LatestPrice(), // last price from data event
	}

	switch order.GetDirection() {
	case "buy":
		f.Direction = "BOT"
	case "sell":
		f.Direction = "SLD"
	}

	f.Commission = e.calculateCommission(float64(f.Qty), f.Price)
	f.ExchangeFee = e.calculateExchangeFee()
	f.Cost = e.calculateCost(f.Commission, f.ExchangeFee)

	return f, nil
}

// calculateComission() calculates the commission for a stock trade
func (e *Exchange) calculateCommission(qty, price float64) float64 {
	// var comMin =
	// var comMax =
	var comRate = e.CommissionRate // 0.0025 // Poloniex market taker fee

	// switch {
	// case (qty * price * comRate) < comMin:
	// 	return comMin
	// case (qty * price * comRate) > comMax:
	// 	return comMax
	// default:
	// Round to 4 decimals
	return math.Floor(qty*price*comRate*10000) / 10000
	// }
}

// calculateExchangeFee() calculates the exchange fee for a stock trade
func (e *Exchange) calculateExchangeFee() float64 {
	return e.ExchangeFee
}

// calculateCost() calculates the total cost for a stock trade
func (e *Exchange) calculateCost(commission, fee float64) float64 {
	return commission + fee
}
