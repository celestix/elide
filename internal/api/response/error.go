package response

import "encoding/json"

func InitError(err error) []byte {
	if err == nil {
		return Error("Unknown")
	}
	return Error(err.Error())
}

func Error(err string) []byte {
	b, _ := json.Marshal(Response{
		Ok:    false,
		Error: err,
	})
	return b
}
