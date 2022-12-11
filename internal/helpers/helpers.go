package helpers

import (
	"fmt"
	"strconv"
	"strings"
)

func PatchChatIdFromBotApi(chatId int64) int64 {
	chat := strings.TrimPrefix(fmt.Sprint(chatId), "-100")
	chatId, _ = strconv.ParseInt(chat, 10, 64)
	return chatId
}

func PatchChatIdToBotApi(chatId int64) int64 {
	if chat := fmt.Sprint(chatId); !strings.HasPrefix(chat, "-100") {
		chat = "-100" + chat
		chatId, _ = strconv.ParseInt(chat, 10, 64)
	}
	return chatId
}
