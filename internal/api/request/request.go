package request

import "encoding/json"

type Request struct {
	Method string          `json:"method"`
	Data   json.RawMessage `json:"data"`
}

func Parse(b []byte) (*Request, error) {
	var r Request
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
