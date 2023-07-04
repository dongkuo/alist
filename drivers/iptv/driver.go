package iptv

import (
	"context"
	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/errs"
	"github.com/alist-org/alist/v3/internal/model"
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
	d.request(d.PlaylistUrl)
	return nil, errs.NotImplement
}

func (d *IPTV) Link(ctx context.Context, file model.Obj, args model.LinkArgs) (*model.Link, error) {
	// TODO return link of file, required
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

func (d *IPTV) request(uri string) ([]byte, error) {
	return nil, nil
}

var _ driver.Driver = (*IPTV)(nil)
