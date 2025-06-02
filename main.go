package main

import "fmt"

// USAGE EXAMPLE
func main() {
	portfolio := Portfolio{
		Stocks:          map[Stock]float64{AAPL: 20, META: 6},
		AllocatedStocks: map[Stock]float64{AAPL: 0.6, META: 0.4},
	}

	transaction := portfolio.Rebalance()

	fmt.Println("BUY  ", transaction.Buy)
	fmt.Println("SELL ", transaction.Sell)
}

// IMPLEMENTATION STARTS HERE

type Stock string

const (
	AAPL Stock = "AAPL"
	META Stock = "META"
)

// CurrentPrice returns the current price of a given stock.
// It's been stubbed for the purpose of the task.
func (s Stock) CurrentPrice() float64 {
	// stubbing the current price. Should be an API call instead.
	currentPrice := map[Stock]float64{
		AAPL: 201.28,
		META: 670.05,
	}

	return currentPrice[s]
}

// Portfolio contains the number of shares per stock,
// and the target allocation for each stock.
type Portfolio struct {
	Stocks          map[Stock]float64 // stock -> shares
	AllocatedStocks map[Stock]float64 // stock -> distribution
}

func (p Portfolio) totalValue() float64 {
	var value float64
	for stock, shares := range p.Stocks {
		value += shares * stock.CurrentPrice()
	}

	return value
}

// Transaction contains the number of stocks to sell and or buy, if any.
type Transaction struct {
	Sell map[Stock]float64
	Buy  map[Stock]float64
}

// Rebalance determines which stocks should be sold and which should be bought
// based on a given portfolio allocation indicating the distribution of stocks
// the portfolio is aiming.
func (p Portfolio) Rebalance() Transaction {
	t := Transaction{
		Sell: make(map[Stock]float64),
		Buy:  make(map[Stock]float64),
	}

	portfolioValue := p.totalValue()

	for stock, shares := range p.Stocks {
		currentValue := shares * stock.CurrentPrice()
		targetValue := p.AllocatedStocks[stock] * portfolioValue

		// we have two scenarios:
		//  1. valueDiff > 0 (currentValue > targetValue): we have a surplus amount
		//     of stocks which we can sell in order to meet our target allocation.
		//  2. valueDiff < 0 (currentValue < targetValue): we are under-allocated, so
		//     we need to buy more stock shares to meet our target allocation.
		valueDiff := currentValue - targetValue
		// the amount of shares to buy/sell. It's an absolute value so we need
		// to make it always posive.
		sharesDiff := valueDiff / stock.CurrentPrice()

		if valueDiff > 0 {
			t.Sell[stock] = sharesDiff
		}

		if valueDiff < 0 {
			t.Buy[stock] = -sharesDiff
		}
	}

	return t
}
