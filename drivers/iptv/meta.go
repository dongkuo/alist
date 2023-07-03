package iptv

import (
	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/op"
)

type Addition struct {
	PlaylistUrl string `json:"playlist_url" help:"The url typically ends with .m3u or .m3u8"  required:"true"`
}

var config = driver.Config{
	Name:              "IPTV",
	LocalSort:         false,
	OnlyLocal:         false,
	OnlyProxy:         false,
	NoCache:           false,
	NoUpload:          false,
	NeedMs:            false,
	DefaultRoot:       "root",
	CheckStatus:       false,
	NoOverwriteUpload: false,
}

func init() {
	op.RegisterDriver(func() driver.Driver {
		return &IPTV{}
	})
}
