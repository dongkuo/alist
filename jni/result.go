package jni

import (
	"fmt"
	"github.com/alist-org/alist/v3/drivers/alist_v3"
	"github.com/alist-org/alist/v3/pkg/utils"
)

const okCode = 200
const okMessage = "success"
const commonErrCode = 1
const jsonMarshalErr = `{"code":400,"err":"json marshal error: %s"}`

type Result struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ListFileData struct {
	Content  []alist_v3.ObjResp `json:"content"`
	Total    int64              `json:"total"`
	Provider string             `json:"provider"`
}

func OK(data any) string {
	return buildResult(okCode, data, okMessage)
}

func Error(message string) string {
	return ErrorWithCode(commonErrCode, message)
}

func ErrorWithCode(code int, message string) string {
	return buildResult(code, nil, message)
}

func buildResult(code int, data any, message string) string {
	res := Result{Code: code, Data: data, Message: message}
	json, err := utils.Json.MarshalToString(res)
	if err != nil {
		return fmt.Sprintf(jsonMarshalErr, err.Error())
	}
	return json
}
