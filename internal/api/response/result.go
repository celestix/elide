package response

import (
	"encoding/json"
)

type Response struct {
	Ok     bool   `json:"ok"`
	Error  string `json:"error,omitempty"`
	Result any    `json:"result,omitempty"`
}

func Result(res any) []byte {
	b, _ := json.Marshal(Response{
		Ok:     true,
		Result: res,
	})
	return b
}
