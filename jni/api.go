package jni

import (
	"context"
	"fmt"
	"github.com/alist-org/alist/v3/cmd/flags"
	_ "github.com/alist-org/alist/v3/drivers"
	"github.com/alist-org/alist/v3/drivers/alist_v3"
	"github.com/alist-org/alist/v3/internal/bootstrap"
	"github.com/alist-org/alist/v3/internal/bootstrap/data"
	"github.com/alist-org/alist/v3/internal/conf"
	"github.com/alist-org/alist/v3/internal/db"
	"github.com/alist-org/alist/v3/internal/errs"
	"github.com/alist-org/alist/v3/internal/fs"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/internal/op"
	"github.com/alist-org/alist/v3/pkg/utils"
	"github.com/alist-org/alist/v3/server"
	"github.com/alist-org/alist/v3/server/common"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var httpSrv *http.Server

func InitConfig(args string) string {
	log.Infof("prepare to init...")
	var initArgs InitArgs
	err := utils.Json.UnmarshalFromString(args, &initArgs)
	if err != nil {
		log.Errorf("failed unmarshal init args: %+v", err)
		return Error("failed unmarshal init args, " + err.Error())
	}

	flags.DataDir = initArgs.DataDir
	flags.Debug = initArgs.Debug
	flags.NoPrefix = initArgs.NoPrefix
	flags.Dev = initArgs.Dev
	flags.ForceBinDir = initArgs.ForceBinDir
	flags.LogStd = initArgs.LogStd

	bootstrap.InitConfig()
	bootstrap.Log()
	log.SetOutput(log.StandardLogger().Out)
	bootstrap.InitDB()
	data.InitData()
	bootstrap.InitIndex()
	bootstrap.InitAria2()
	bootstrap.InitQbittorrent()
	loadStorages()
	log.Infof("init done!")
	println("InitConfig done, DataDir: ", flags.DataDir)
	return OK("init done!")
}

func loadStorages() {
	storages, err := db.GetEnabledStorages()
	if err != nil {
		utils.Log.Fatalf("failed get enabled storages: %+v", err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func(storages []model.Storage) {
		defer wg.Done()
		for i := range storages {
			err := op.LoadStorage(context.Background(), storages[i])
			if err != nil {
				utils.Log.Errorf("failed get enabled storages: %+v", err)
			} else {
				utils.Log.Infof("success load storage: [%s], driver: [%s]",
					storages[i].MountPath, storages[i].Driver)
			}
		}
		conf.StoragesLoaded = true
	}(storages)
	wg.Wait()
}

func StartServer() string {
	if httpSrv != nil {
		return Error("http server has been running")
	}
	log.Infof("start http server...")
	if !flags.Debug && !flags.Dev {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.LoggerWithWriter(log.StandardLogger().Out), gin.RecoveryWithWriter(log.StandardLogger().Out))
	server.Init(r)
	httpBase := fmt.Sprintf("%s:%d", conf.Conf.Address, conf.Conf.Port)

	httpSrv = &http.Server{Addr: httpBase, Handler: r}

	log.Infof("http server start at %s", httpBase)
	go func() {
		err := httpSrv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server failed to start: %s", err.Error())
			httpSrv = nil
		}
	}()

	go func() {
		var quit = make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		shutdown()
	}()

	return OK("http server start at " + httpBase)
}

func StopServer() string {
	shutdown()
	return OK("http server exit")
}

func shutdown() {
	if httpSrv == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	var waitShutdown sync.WaitGroup
	waitShutdown.Add(1)
	go func() {
		defer waitShutdown.Done()
		if err := httpSrv.Shutdown(ctx); err != nil {
			log.Fatal("http server shutdown:", err)
		}
	}()
	waitShutdown.Wait()
	httpSrv = nil

	log.Println("http server exit")
}

func GetAdmin() string {
	admin, err := op.GetAdmin()
	if err != nil {
		log.Errorf("failed get admin user: %+v", err)
		return Error("failed get admin user, " + err.Error())
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
		return Error("failed unmarshal args json, " + err.Error())
	}

	admin, err := op.GetAdmin()
	if err != nil {
		log.Errorf("failed get admin user: %+v", err)
		return Error("failed get admin user, " + err.Error())
	}

	reqPath, err := admin.JoinPath(req.Path)
	if err != nil {
		log.Errorf("failed get admin user req path: %+v", err)
		return Error("failed get admin user req path, " + err.Error())
	}

	meta, err := op.GetNearestMeta(reqPath)
	if err != nil {
		if !errors.Is(errors.Cause(err), errs.MetaNotFound) {
			return Error("failed get meta, " + err.Error())
		}
	}

	var ctx = Context{}
	ctx.Set("user", admin)
	ctx.Set("meta", meta)

	objs, err := fs.List(&ctx, reqPath, &fs.ListArgs{Refresh: req.Refresh})
	if err != nil {
		return Error("failed list file, " + err.Error())
	}

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
