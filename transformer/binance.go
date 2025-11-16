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

type BinanceTransformer struct {
	allowed map[string]struct{}
}

var binanceSymbolMap = loadBinanceSymbolMap()

func loadBinanceSymbolMap() map[string]string {
	m := make(map[string]string)
	raw := os.Getenv("SYMBOL_MAP")
	if raw == "" {
		return m
	}
	for _, pair := range strings.Split(raw, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		parts := strings.Split(pair, ":")
		if len(parts) != 2 {
			continue
		}
		key := strings.ToUpper(strings.TrimSpace(parts[0]))
		val := strings.TrimSpace(parts[1])
		if key == "" || val == "" {
			continue
		}
		m[key] = val
	}
	return m
}

func NewBinanceTransformer(symbols []string) *BinanceTransformer {
	allowed := make(map[string]struct{})
	for _, s := range symbols {
		allowed[strings.ToUpper(strings.TrimSpace(s))] = struct{}{}
	}
	return &BinanceTransformer{allowed: allowed}
}

func ptr[T any](v T) *T { return &v }

func makeMsgID(source, symbol string, tsEventNs int64, tradeID int64) string {
	return fmt.Sprintf("%s|%s|%d|%d", source, symbol, tsEventNs, tradeID)
}

func (bt *BinanceTransformer) Transform(payload []byte) (*CanonicalEvent, error) {
	var t binanceTrade
	err := json.Unmarshal(payload, &t)
	if err != nil || t.Symbol == "" {
		var wrapper struct {
			Stream string       `json:"stream"`
			Data   binanceTrade `json:"data"`
		}
		if err2 := json.Unmarshal(payload, &wrapper); err2 != nil {
			if err != nil {
				return nil, err
			}
			return nil, err2
		}
		t = wrapper.Data
	}

	symbol := strings.ToUpper(t.Symbol)
	if _, ok := bt.allowed[symbol]; !ok {
		return nil, nil
	}

	if mapped, ok := binanceSymbolMap[symbol]; ok {
		symbol = mapped
	}

	source := os.Getenv("SOURCE")
	assetClass := os.Getenv("ASSET_CLASS")
	eventType := os.Getenv("EVENT_TYPE")
	version := os.Getenv("VERSION")

	tsEventNs := t.TradeTime * int64(time.Millisecond)
	tsIngestNs := time.Now().UnixNano()
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
