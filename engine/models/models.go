package models

type Candle struct {
	Date   string  `json:"date"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int64   `json:"volume"`
}

type Stock struct {
	Symbol       string   `json:"symbol"`
	Name         string   `json:"name"`
	CurrentPrice float64  `json:"current_price"`
	Currency     string   `json:"currency"`
	MarketCap    float64  `json:"market_cap"`
	High52w      float64  `json:"52w_high"`
	Low52w       float64  `json:"52w_low"`
	History      []Candle `json:"history"`
	UpdatedAt    string   `json:"updated_at"`
	Error        string   `json:"error,omitempty"`
}

type Crypto struct {
	ID           string         `json:"id"`
	CurrentPrice float64        `json:"current_price"`
	MarketCap    float64        `json:"market_cap"`
	Volume24h    float64        `json:"volume_24h"`
	Change24h    float64        `json:"change_24h"`
	History      []CryptoCandle `json:"history"`
	UpdatedAt    string         `json:"updated_at"`
	Error        string         `json:"error,omitempty"`
}

type CryptoCandle struct {
	Date  string  `json:"date"`
	Price float64 `json:"price"`
}

type Metrics struct {
	Symbol    string  `json:"symbol"`
	MA7       float64 `json:"ma_7"`
	MA14      float64 `json:"ma_14"`
	MA30      float64 `json:"ma_30"`
	Change7d  float64 `json:"change_7d_pct"`
	Change30d float64 `json:"change_30d_pct"`
	Trend     string  `json:"trend"`
}

type Alert struct {
	Symbol    string  `json:"symbol"`
	Type      string  `json:"type"`
	Message   string  `json:"message"`
	Threshold float64 `json:"threshold_pct"`
	Triggered float64 `json:"triggered_at"`
	Time      string  `json:"time"`
}
