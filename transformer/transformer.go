package transformer

type CanonicalEvent struct {
	Version       string  `json:"version"`
	EventType     string  `json:"event_type"`
	AssetClass    string  `json:"asset_class"`
	Source        string  `json:"source"`
	Symbol        string  `json:"symbol"`
	TsEventNanos  int64   `json:"ts_event_nanos"`
	TsIngestNanos int64   `json:"ts_ingest_nanos"`
	PriceStr      *string `json:"price,omitempty"`
	SizeStr       *string `json:"size,omitempty"`
	MsgID         string  `json:"msg_id"`
	Raw           []byte  `json:"raw,omitempty"`
}

type Transformer interface {
	Transform(payload []byte) (*CanonicalEvent, error)
}