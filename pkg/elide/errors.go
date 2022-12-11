package elide

type TelegramError struct {
	rpc string
}

func (t *TelegramError) Error() string {
	return t.rpc
}
