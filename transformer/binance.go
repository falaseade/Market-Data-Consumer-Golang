package transformer

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type binanceTrade struct {
	Symbol    string `json:"s"`
	Price     string `json:"p"`
	Quantity  string `json:"q"`
	TradeID   int64  `json:"t"`
	TradeTime int64  `json:"T"`
	IsBuyerMM bool   `json:"m"`
}

type BinanceTransformer struct {}

func NewBinanceTransformer() *BinanceTransformer {
	return &BinanceTransformer{}
}

func ptr[T any](v T) *T { return &v }

func makeMsgID(source, symbol string, tsEventNs int64, tradeID int64) string {
	return fmt.Sprintf("%s|%s|%d|%d", source, symbol, tsEventNs, tradeID)
}

func (bt *BinanceTransformer) Transform(payload []byte) (*CanonicalEvent, error) {
	var t binanceTrade
	if err := json.Unmarshal(payload, &t); err != nil {
		return nil, err
	}

	source := os.Getenv("SOURCE")
	assetClass := os.Getenv("ASSET_CLASS")
	eventType := os.Getenv("EVENT_TYPE")
	version:= os.Getenv("VERSION")

	tsEventNs := t.TradeTime * int64(time.Millisecond)
	tsIngestNs := time.Now().UnixNano()
	symbol := strings.ToUpper(t.Symbol)
	msgID := makeMsgID(source, symbol, tsEventNs, t.TradeID)

	evt := &CanonicalEvent{
		Version:       version,
		EventType:     eventType,
		AssetClass:    assetClass,
		Source:        source,
		Symbol:        symbol,
		TsEventNanos:  tsEventNs,
		TsIngestNanos: tsIngestNs,
		PriceStr:      ptr(t.Price),
		SizeStr:       ptr(t.Quantity),
		MsgID:         msgID,
		Raw:           payload,
	}
	return evt, nil
}