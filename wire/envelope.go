package wire

type Envelope struct {
	Type string      `json:"type"`
	Body interface{} `json:"body,omitempty"`
}
