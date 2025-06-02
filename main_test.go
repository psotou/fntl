package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPortfolio_Rebalance(t *testing.T) {
	tests := []struct {
		name      string
		portfolio Portfolio
		want      Transaction
	}{
		{
			name: "Sell META buy AAPL",
			// currentValueAAPL = 20 * 201.28 = 4025.6
			// currentValueMETA =  6 * 670.05 = 4020.3
			// totalValue = currentValueAAPL + currentValueMETA = 8045.9
			//
			// targetValueAAPL = 0.6 * totalValue = 4827.54
			// valueDiffAAPL = currentValueAAPL - targetValueAAPL = -801.94
			// Since valueDiffAAPL < 0, we need to buy more AAPL to meet the requirement.
			// sharesToBuyAAPL  = |valueDiffAAPL| / 201.28 = 3.9837
			// sharesToSellMETA = |valueDiffMETA| / 670.05 = 1.1968
			portfolio: Portfolio{
				Stocks:          map[Stock]float64{AAPL: 20, META: 6},
				AllocatedStocks: map[Stock]float64{AAPL: 0.6, META: 0.4},
			},
			want: Transaction{
				Sell: map[Stock]float64{META: 1.1968},
				Buy:  map[Stock]float64{AAPL: 3.9837},
			},
		},
		{
			name: "Portfolio is balanced",
			// Taking the first test case, we have:
			// sharesAAPL + sharesToBuyAAPL = 20 + 3.9837 = 23.9837
			// sharesMETA + sharesToBuyMETA =  6 + 1.1968 =  4.8032
			portfolio: Portfolio{
				Stocks:          map[Stock]float64{AAPL: 23.9837, META: 4.8032},
				AllocatedStocks: map[Stock]float64{AAPL: 0.6, META: 0.4},
			},
			want: Transaction{
				Sell: make(map[Stock]float64),
				Buy:  make(map[Stock]float64),
			},
		},
		{
			name: "Sell all AAPL, buy META",
			// totalValue = 8045.9
			// currentValueAAPL = 4025.6
			// targetValueAAPL = 0.0
			// valueDiffAAPL = 4025.6, which is > 0, so we sell.
			// sharesToSellAAPL = 20
			// sharesToBuyMETA  = 4025.6 / 670.05 = 6.0079
			portfolio: Portfolio{
				Stocks:          map[Stock]float64{AAPL: 20, META: 6},
				AllocatedStocks: map[Stock]float64{AAPL: 0.0, META: 1.0},
			},
			want: Transaction{
				Sell: map[Stock]float64{AAPL: 20},
				Buy:  map[Stock]float64{META: 6.0079},
			},
		},
		{
			name: "Sell all META, buy AAPL",
			// totalValue = 8045.9
			// currentValueMETA = 4020.3
			// targetValueMETA = 0.0
			// valueDiffMETA = 4020.3, which is > 0, so we sell.
			// sharesToSellMETA  = 6
			// sharesToBuyAAPL = 4020.3 / 201.28 = 19.9737
			portfolio: Portfolio{
				Stocks:          map[Stock]float64{AAPL: 20, META: 6},
				AllocatedStocks: map[Stock]float64{AAPL: 1.0, META: 0.0},
			},
			want: Transaction{
				Sell: map[Stock]float64{META: 6},
				Buy:  map[Stock]float64{AAPL: 19.9737},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.portfolio.Rebalance()

			delta := 0.001

			// Sell AAPL
			if !assert.InDelta(t, tt.want.Sell[AAPL], got.Sell[AAPL], delta) {
				t.Errorf("%s':\nGot: %+v\nWant:%+v", tt.name, got, tt.want)
			}

			// Sell META
			if !assert.InDelta(t, tt.want.Sell[META], got.Sell[META], delta) {
				t.Errorf("%s':\nGot: %+v\nWant:%+v", tt.name, got, tt.want)
			}

			// Buy AAPL
			if !assert.InDelta(t, tt.want.Buy[AAPL], got.Buy[AAPL], delta) {
				t.Errorf("%s':\nGot: %+v\nWant:%+v", tt.name, got, tt.want)
			}

			// Buy META
			if !assert.InDelta(t, tt.want.Buy[META], got.Buy[META], delta) {
				t.Errorf("%s':\nGot: %+v\nWant:%+v", tt.name, got, tt.want)
			}
		})
	}
}
