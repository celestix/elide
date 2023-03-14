package methods

import (
	photos "elide/internal/helpers/photos"
	telegraph "elide/internal/helpers/telegraph"
	"elide/pkg/elide"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/gotgproto/generic"
	"github.com/gotd/td/tg"
	"go.uber.org/multierr"
)

func GetChatPhoto(ctx *ext.Context, body json.RawMessage) (any, error) {
	var p elide.GetChatPhotoBody
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, fmt.Errorf("failed to read body for getChatPhoto: %s", err.Error())
	}
	if p.ChatId == nil {
		return nil, errors.New("chat_id not provided")
	}
	var (
		chat   tg.ChatFullClass
		err    error
		chatId int64
	)
	switch cid := p.ChatId.(type) {
	case int:
		chat, err = generic.GetChat(ctx, cid)
	case int64:
		chat, err = generic.GetChat(ctx, cid)
	case string:
		chat, err = generic.GetChat(ctx, cid)
	}
	if err != nil {
		return nil, err
	}
	var photo tg.PhotoClass
	switch v := chat.(type) {
	case *tg.ChannelFull:
		photo = v.ChatPhoto
	case *tg.ChatFull:
		photo = v.ChatPhoto
	}
	var errPhotoNotFound = errors.New("photo not found")
	if photo == nil {
		return nil, errPhotoNotFound
	}
	pPaths, photoErrs := photos.Download(ctx, chatId, []tg.PhotoClass{photo})
	if pPaths == nil {
		if photoErrs == nil {
			return nil, errPhotoNotFound
		}
		return nil, photoErrs
	}
	urls, tgraphErrs := telegraph.UploadFiles(pPaths)
	if urls != nil {
		return urls[0], photoErrs
	}
	return nil, multierr.Append(photoErrs, tgraphErrs)
}
