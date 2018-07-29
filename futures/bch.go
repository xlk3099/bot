package futures

type bch struct {
	symbol       string
	contractType string
	date         string
	ticker       Ticker
}

// Ticker 交易对
type Ticker struct {
	High       float64 `json:"high"`
	Vol        int     `json:"vol"` // 24 小时成交量
	Last       float64 `json:"last"`
	Low        float64 `json:"low"`
	ContractID string  `json:"contract_id"`
	Buy        float64 `json:"buy"`  // 买一价
	Sell       float64 `json:"sell"` // 卖一价
	UnitAmount int     `json:"unit_amount"`
}
