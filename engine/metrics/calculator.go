package metrics

import (
	"finance-dashboard/engine/models"
	"math"
)

func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return math.Round((sum/float64(len(values)))*100) / 100
}

func percentChange(from, to float64) float64 {
	if from == 0 {
		return 0
	}
	change := ((to - from) / from) * 100
	return math.Round(change*100) / 100
}

func lastNCloses(history []models.Candle, n int) []float64 {
	if n <= 0 || len(history) == 0 {
		return nil
	}
	start := len(history) - n
	if start < 0 {
		start = 0
	}
	closes := make([]float64, 0, len(history)-start)
	for _, c := range history[start:] {
		closes = append(closes, c.Close)
	}
	return closes
}

func lastNPrices(history []models.CryptoCandle, n int) []float64 {
	if n <= 0 || len(history) == 0 {
		return nil
	}
	start := len(history) - n
	if start < 0 {
		start = 0
	}
	prices := make([]float64, 0, len(history)-start)
	for _, c := range history[start:] {
		prices = append(prices, c.Price)
	}
	return prices
}

func CalculateStockMetrics(stock models.Stock) models.Metrics {
	history := stock.History

	ma7 := average(lastNCloses(history, 7))
	ma14 := average(lastNCloses(history, 14))
	closes30 := lastNCloses(history, 30)
	ma30 := average(closes30)

	current := stock.CurrentPrice
	change7d, change30d := 0.0, 0.0

	if len(history) >= 7 {
		change7d = percentChange(history[len(history)-7].Close, current)
	}
	if len(closes30) >= 30 {
		change30d = percentChange(closes30[0], current)
	}

	trend := "neutral"
	if ma7 > ma30 {
		trend = "up"
	} else if ma7 < ma30 {
		trend = "down"
	}

	return models.Metrics{
		Symbol:    stock.Symbol,
		MA7:       ma7,
		MA14:      ma14,
		MA30:      ma30,
		Change7d:  change7d,
		Change30d: change30d,
		Trend:     trend,
	}
}

func CalculateCryptoMetrics(crypto models.Crypto) models.Metrics {
	prices30 := lastNPrices(crypto.History, 30)

	ma7 := average(lastNPrices(crypto.History, 7))
	ma14 := average(lastNPrices(crypto.History, 14))
	ma30 := average(prices30)

	current := crypto.CurrentPrice
	change7d, change30d := 0.0, 0.0

	if len(prices30) >= 7 {
		change7d = percentChange(prices30[len(prices30)-7], current)
	}
	if len(prices30) >= 30 {
		change30d = percentChange(prices30[0], current)
	}

	trend := "neutral"
	if ma7 > ma30 {
		trend = "up"
	} else if ma7 < ma30 {
		trend = "down"
	}

	return models.Metrics{
		Symbol:    crypto.ID,
		MA7:       ma7,
		MA14:      ma14,
		MA30:      ma30,
		Change7d:  change7d,
		Change30d: change30d,
		Trend:     trend,
	}
}
