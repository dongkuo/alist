package jni

import (
	"github.com/alist-org/alist/v3/cmd"
	"github.com/alist-org/alist/v3/cmd/flags"
	"github.com/alist-org/alist/v3/internal/op"
	"github.com/alist-org/alist/v3/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func ExecCommand(jsonArgs string) string {
	var args []string
	err := utils.Json.UnmarshalFromString(jsonArgs, &args)
	if err != nil {
		log.Errorf("failed unmarshal storageConfig: %+v", err)
		return resErr("failed unmarshal args json: " + err.Error())
	}

	cmd.RootCmd.SetArgs(args)
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Errorf("execute error: %+v", err)
		return resErr("execute error: " + err.Error())
	}
	return resOk("")
}

func GetAdmin(dataDir string) string {
	flags.DataDir = dataDir
	cmd.Init()
	admin, err := op.GetAdmin()
	if err != nil {
		log.Errorf("failed get admin user: %+v", err)
		return resErr("failed get admin user: " + err.Error())
	} else {
		log.Infof("admin user's info: \nusername: %s\npassword: %s", admin.Username, admin.Password)
		return resOk(account{admin.Username, admin.Password})
	}
}
