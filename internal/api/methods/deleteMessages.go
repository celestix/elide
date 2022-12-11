package methods

import (
	"elide/internal/helpers"
	"elide/pkg/elide"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/gotgproto/storage"
	"github.com/gotd/td/tg"
)

func deleteMessages(ctx *ext.Context, chatId int64, messageIds []int, revoke bool) error {
	if chatId == 0 {
		_, err := ctx.Client.MessagesDeleteMessages(ctx.Context, &tg.MessagesDeleteMessagesRequest{
			ID:     messageIds,
			Revoke: revoke,
		})
		return err
	}
	peer := storage.GetPeerById(chatId)
	if peer.ID == 0 {
		return ext.ErrPeerNotFound
	}
	_, err := ctx.Client.ChannelsDeleteMessages(ctx.Context, &tg.ChannelsDeleteMessagesRequest{
		Channel: &tg.InputChannel{
			ChannelID:  peer.ID,
			AccessHash: peer.AccessHash,
		},
		ID: messageIds,
	})
	return err
}

func DeleteMessages(ctx *ext.Context, body json.RawMessage) (any, error) {
	var p elide.DeleteMessagesBody
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, fmt.Errorf("failed to read body for deleteMessages: %s", err.Error())
	}
	if p.ChatId != 0 {
		p.ChatId = helpers.PatchChatIdFromBotApi(p.ChatId)
	}
	if p.MessageIds == nil || len(p.MessageIds) == 0 {
		return nil, errors.New("message ids not provided")
	}
	return nil, deleteMessages(ctx, p.ChatId, p.MessageIds, p.Revoke)
}
