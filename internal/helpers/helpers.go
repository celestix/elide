package helpers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/gotgproto/storage"
)

type Peer struct {
	ID         int64
	AccessHash int64
}

func GetPeerByUsername(ctx *ext.Context, username string) (*Peer, error) {
	peer := storage.GetPeerByUsername(username)
	if peer.ID != 0 {
		return &Peer{peer.ID, peer.AccessHash}, nil
	}
	chat, err := ctx.ResolveUsername(username)
	if err != nil {
		return nil, err
	}
	return &Peer{chat.GetID(), chat.GetAccessHash()}, nil
}

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
