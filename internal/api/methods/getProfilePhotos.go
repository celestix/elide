package methods

import (
	"elide/pkg/elide"
	"encoding/json"
	"errors"
	"fmt"

	"elide/internal/helpers"
	photos "elide/internal/helpers/photos"
	telegraph "elide/internal/helpers/telegraph"

	"github.com/anonyindian/gotgproto/ext"
	"github.com/gotd/td/tg"
	"go.uber.org/multierr"
)

func GetProfilePhotos(ctx *ext.Context, body json.RawMessage) (any, error) {
	var p elide.GetProfilePhotosBody
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, fmt.Errorf("failed to read body for getProfilePhotos: %s", err.Error())
	}
	if p.UserID == nil {
		return nil, errors.New("user_id not provided")
	}
	var (
		rawPhotos     []tg.PhotoClass
		err           error
		userId        int64
		requestParams = &tg.PhotosGetUserPhotosRequest{
			Offset: p.Offset,
			MaxID:  p.MaxID,
			Limit:  p.Limit,
		}
	)
	switch uid := p.UserID.(type) {
	case int:
		rawPhotos, err = ctx.GetUserProfilePhotos(int64(uid), requestParams)
		if err != nil {
			return nil, err
		}
		userId = int64(uid)
	case int64:
		rawPhotos, err = ctx.GetUserProfilePhotos(uid, requestParams)
		if err != nil {
			return nil, err
		}
		userId = uid
	case string:
		peer, err := helpers.GetPeerByUsername(ctx, uid)
		if err != nil {
			return nil, err
		}
		requestParams.UserID = &tg.InputUser{
			UserID:     peer.ID,
			AccessHash: peer.AccessHash,
		}
		rawRawPhotos, err := ctx.Client.PhotosGetUserPhotos(ctx, requestParams)
		if err != nil {
			return nil, err
		}
		rawPhotos = rawRawPhotos.GetPhotos()
		userId = peer.ID
	default:
		return nil, errors.New("unsupported type used for user_id")
	}
	pPaths, photoErrs := photos.Download(ctx, userId, rawPhotos)
	if pPaths == nil {
		if photoErrs == nil {
			return nil, errors.New("no photos found")
		}
		return nil, photoErrs
	}
	urls, tgraphErrs := telegraph.UploadFiles(pPaths)
	return urls, multierr.Append(photoErrs, tgraphErrs)
}
