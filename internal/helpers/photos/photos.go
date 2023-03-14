package helpers

import (
	"elide/internal/api/config"
	"fmt"
	"os"
	"strconv"

	"github.com/anonyindian/gotgproto/ext"
	"github.com/gotd/td/telegram/downloader"
	"github.com/gotd/td/tg"
	"go.uber.org/multierr"
)

func Download(ctx *ext.Context, userId int64, p []tg.PhotoClass) ([]string, error) {
	d := downloader.NewDownloader()
	var err error
	var pPaths []string
	for _, pc := range p {
		photo, ok := pc.AsNotEmpty()
		if !ok {
			continue
		}
		var loc tg.InputFileLocationClass
		if photo.VideoSizes != nil {
			tType := ""
			for _, vsz := range photo.VideoSizes {
				vs, ok := vsz.(*tg.VideoSize)
				if !ok {
					continue
				}
				tType = vs.Type
			}
			loc = &tg.InputPhotoFileLocation{
				ID:            photo.ID,
				AccessHash:    photo.AccessHash,
				FileReference: photo.FileReference,
				ThumbSize:     tType,
			}
		} else {
			loc = &tg.InputPhotoFileLocation{
				ID:            photo.ID,
				AccessHash:    photo.AccessHash,
				FileReference: photo.FileReference,
				ThumbSize:     photo.Sizes[len(photo.Sizes)-1].GetType(),
			}
		}
		userDir := fmt.Sprintf("%s/downloads/photos/%s", config.WorkingDir, strconv.FormatInt(userId, 10))
		pPath := fmt.Sprintf("%s/%d_%d", userDir, photo.ID, photo.AccessHash)

		if _, err111 := os.Stat(pPath); os.IsNotExist(err111) {
			config.Debugf("[Elide-DEBUG][Photos][Downloader]: file %d-%d not found in cache\n", photo.ID, photo.AccessHash)
			err1 := downloadFile(ctx, d, loc, userDir, pPath)
			if err != nil {
				err = multierr.Append(err, err1)
				continue
			}
		} else {
			config.Debugf("[Elide-DEBUG][Photos][Downloader]: file %d-%d found in cache\n", photo.ID, photo.AccessHash)
		}
		pPaths = append(pPaths, pPath)
	}
	return pPaths, err
}

func downloadFile(ctx *ext.Context, d *downloader.Downloader, loc tg.InputFileLocationClass, userDir, pPath string) error {
	db := d.Download(ctx.Client, loc)
	_, fErr := os.Stat(userDir)
	if fErr != nil {
		os.Mkdir(userDir, os.ModeDir)
	}
	_, err1 := db.ToPath(ctx, pPath)
	return err1
}
