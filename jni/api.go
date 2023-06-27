package jni

import (
	"github.com/alist-org/alist/v3/cmd"
	"github.com/alist-org/alist/v3/cmd/flags"
	"github.com/alist-org/alist/v3/drivers/alist_v3"
	"github.com/alist-org/alist/v3/internal/errs"
	"github.com/alist-org/alist/v3/internal/fs"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/internal/op"
	"github.com/alist-org/alist/v3/pkg/utils"
	"github.com/alist-org/alist/v3/server/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func ExecCommand(jsonArgs string) string {
	var args []string
	err := utils.Json.UnmarshalFromString(jsonArgs, &args)
	if err != nil {
		log.Errorf("failed unmarshal storageConfig: %+v", err)
		return Error("failed unmarshal args json: " + err.Error())
	}

	cmd.RootCmd.SetArgs(args)
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Errorf("execute error: %+v", err)
		return Error("execute error: " + err.Error())
	}
	return OK("")
}

func GetAdmin(dataDir string) string {
	flags.DataDir = dataDir
	cmd.Init()
	admin, err := op.GetAdmin()
	if err != nil {
		log.Errorf("failed get admin user: %+v", err)
		return Error("failed get admin user: " + err.Error())
	} else {
		log.Infof("admin user's info: \nusername: %s\npassword: %s", admin.Username, admin.Password)
		return OK(Account{admin.Username, admin.Password})
	}
}

func ListFile(jsonArgs string) string {
	var req alist_v3.ListReq
	err := utils.Json.UnmarshalFromString(jsonArgs, &req)
	if err != nil {
		log.Errorf("failed unmarshal storageConfig: %+v", err)
		return Error("failed unmarshal args json: " + err.Error())
	}

	Init()

	admin, err := op.GetAdmin()
	if err != nil {
		log.Errorf("failed get admin user: %+v", err)
		return Error("failed get admin user: " + err.Error())
	}

	reqPath, err := admin.JoinPath(req.Path)
	if err != nil {
		log.Errorf("failed get admin user req path: %+v", err)
		return Error("failed get admin user req path: " + err.Error())
	}

	meta, err := op.GetNearestMeta(reqPath)
	if err != nil {
		if !errors.Is(errors.Cause(err), errs.MetaNotFound) {
			return Error("failed get meta: " + err.Error())
		}
	}

	var ctx = Context{}
	ctx.Set("user", admin)
	ctx.Set("meta", meta)

	objs, err := fs.List(&ctx, reqPath, &fs.ListArgs{Refresh: req.Refresh})
	total, objs := pagination(objs, &req.PageReq)
	provider := "unknown"
	storage, err := fs.GetStorage(reqPath, &fs.GetStoragesArgs{})
	if err == nil {
		provider = storage.GetStorage().Driver
	}
	return OK(ListFileData{
		Content:  toObjsResp(objs, reqPath, isEncrypt(meta, reqPath)),
		Total:    int64(total),
		Provider: provider,
	})
}

func pagination(objs []model.Obj, req *model.PageReq) (int, []model.Obj) {
	pageIndex, pageSize := req.Page, req.PerPage
	total := len(objs)
	start := (pageIndex - 1) * pageSize
	if start > total {
		return total, []model.Obj{}
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return total, objs[start:end]
}

func isEncrypt(meta *model.Meta, path string) bool {
	if common.IsStorageSignEnabled(path) {
		return true
	}
	if meta == nil || meta.Password == "" {
		return false
	}
	if !utils.PathEqual(meta.Path, path) && !meta.PSub {
		return false
	}
	return true
}

func toObjsResp(objs []model.Obj, parent string, encrypt bool) []alist_v3.ObjResp {
	var resp []alist_v3.ObjResp
	for _, obj := range objs {
		thumb, _ := model.GetThumb(obj)
		resp = append(resp, alist_v3.ObjResp{
			Name:     obj.GetName(),
			Size:     obj.GetSize(),
			IsDir:    obj.IsDir(),
			Modified: obj.ModTime(),
			Sign:     common.Sign(obj, parent, encrypt),
			Thumb:    thumb,
			Type:     utils.GetObjType(obj.GetName(), obj.IsDir()),
		})
	}
	return resp
}
