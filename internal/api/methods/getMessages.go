package methods

import (
	"elide/internal/helpers"
	"elide/pkg/elide"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/gotgproto/functions"
	"github.com/anonyindian/gotgproto/storage"
	"github.com/gotd/td/tg"
)

func getInputMessageIds(p *elide.GetMessagesBody) []tg.InputMessageClass {
	data := make([]tg.InputMessageClass, len(p.MessageIds))
	for i, v := range p.MessageIds {
		data[i] = &tg.InputMessageID{ID: v}
	}
	return data
}

func getMessages(ctx *ext.Context, chatId int64, messageIds []tg.InputMessageClass) ([]tg.MessageClass, error) {
	if chatId == 0 {
		return functions.GetChatMessages(ctx.Context, ctx.Client, messageIds)
	}
	peer := storage.GetPeerById(chatId)
	if peer.ID == 0 {
		return nil, ext.ErrPeerNotFound
	}
	switch storage.EntityType(peer.Type) {
	case storage.TypeChannel:
		return functions.GetChannelMessages(ctx.Context, ctx.Client, &tg.InputChannel{
			ChannelID:  peer.ID,
			AccessHash: peer.AccessHash,
		}, messageIds)
	default:
		return functions.GetChatMessages(ctx.Context, ctx.Client, messageIds)
	}
}

func GetMessages(ctx *ext.Context, body json.RawMessage) (any, error) {
	var p elide.GetMessagesBody
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, fmt.Errorf("failed to read body for getMessagesBody: %s", err.Error())
	}
	if p.ChatId != 0 {
		p.ChatId = helpers.PatchChatIdFromBotApi(p.ChatId)
	}
	if p.MessageIds == nil || len(p.MessageIds) == 0 {
		return nil, errors.New("message ids not provided")
	}
	return getMessages(ctx, p.ChatId, getInputMessageIds(&p))
}
