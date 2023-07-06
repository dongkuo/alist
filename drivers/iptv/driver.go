package iptv

import (
	"context"
	"fmt"
	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/errs"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/pkg/utils"
	"github.com/jamesnetherton/m3u"
	"strings"
	"time"
)

type IPTV struct {
	model.Storage
	Addition
}

func (d *IPTV) Config() driver.Config {
	return config
}

func (d *IPTV) GetAddition() driver.Additional {
	return &d.Addition
}

func (d *IPTV) Init(ctx context.Context) error {
	return nil
}

func (d *IPTV) Drop(ctx context.Context) error {
	return nil
}

func (d *IPTV) List(ctx context.Context, dir model.Obj, args model.ListArgs) ([]model.Obj, error) {
	playlist, err := m3u.Parse(d.PlaylistUrl)
	if err != nil {
		return nil, err
	}
	return utils.SliceConvert(playlist.Tracks, func(track m3u.Track) (model.Obj, error) {
		return trackToObj(track), nil
	})
}

func trackToObj(track m3u.Track) model.Obj {
	var id, name, thumb = track.Name, track.Name, ""
	for _, tag := range track.Tags {
		if tag.Name == "tvg-name" {
			name = tag.Value
		} else if tag.Name == "tvg-id" {
			id = tag.Value
		} else if tag.Name == "tvg-logo" {
			thumb = tag.Value
		}
	}

	if !strings.HasSuffix(name, ".m3u8") {
		name += ".m3u8"
	}

	return &model.ObjThumbURL{
		Object: model.Object{
			ID:       id,
			Name:     name,
			Size:     0,
			Modified: time.Now(),
			IsFolder: false,
		},
		Thumbnail: model.Thumbnail{Thumbnail: thumb},
		Url:       model.Url{Url: track.URI},
	}
}

func (d *IPTV) Link(ctx context.Context, file model.Obj, args model.LinkArgs) (*model.Link, error) {
	if url, ok := file.(model.URL); ok {
		return &model.Link{URL: url.URL()}, nil
	}
	return nil, errs.NotImplement
}

func (d *IPTV) MakeDir(ctx context.Context, parentDir model.Obj, dirName string) error {
	// TODO create folder, optional
	return errs.NotImplement
}

func (d *IPTV) Move(ctx context.Context, srcObj, dstDir model.Obj) error {
	// TODO move obj, optional
	return errs.NotImplement
}

func (d *IPTV) Rename(ctx context.Context, srcObj model.Obj, newName string) error {
	// TODO rename obj, optional
	return errs.NotImplement
}

func (d *IPTV) Copy(ctx context.Context, srcObj, dstDir model.Obj) error {
	// TODO copy obj, optional
	return errs.NotImplement
}

func (d *IPTV) Remove(ctx context.Context, obj model.Obj) error {
	// TODO remove obj, optional
	return errs.NotImplement
}

func (d *IPTV) Put(ctx context.Context, dstDir model.Obj, stream model.FileStreamer, up driver.UpdateProgress) error {
	// TODO upload file, optional
	return errs.NotImplement
}

func (d *IPTV) request(url string) ([]byte, error) {
	playlist, err := m3u.Parse(url)
	if err != nil {
		return nil, err
	}
	for _, track := range playlist.Tracks {
		fmt.Printf("track: name=%s\nuri=%s\ntags=%s\nlen=%d\n\n", track.Name, track.URI, track.Tags, track.Length)
	}
	return nil, nil
}

var _ driver.Driver = (*IPTV)(nil)
