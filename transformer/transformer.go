package transformer

import (
	"encoding/json"
	"strconv"
)

type CanonicalTick struct {
	Source    string  `json:"source"`
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Quantity  float64 `json:"quantity"`
	Timestamp int64   `json:"timestamp"`
}

type Transformer interface {
	Transform(payload []byte) (*CanonicalTick, error)
}

type binanceTick struct {
	Symbol    string `json:"s"`
	Price     string `json:"p"`
	Quantity  string `json:"q"`
	Timestamp int64  `json:"T"`
}

type BinanceTransformer struct {}

func NewBinanceTransformer() *BinanceTransformer {
	return &BinanceTransformer{}
}

func (bt *BinanceTransformer) Transform(payload []byte) (*CanonicalTick, error) {
	var bTick binanceTick
	if err := json.Unmarshal(payload, &bTick); err != nil {
		return nil, err
	}
	price, err := strconv.ParseFloat(bTick.Price, 64)
	if err != nil {
		return nil, err
	}
	quantity, err := strconv.ParseFloat(bTick.Quantity, 64)
	if err != nil {
		return nil, err
	}
	canonical := &CanonicalTick{
		Source:    "binance",
		Symbol:    bTick.Symbol,
		Price:     price,
		Quantity:  quantity,
		Timestamp: bTick.Timestamp,
	}
	return canonical, nil
}