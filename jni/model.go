package jni

import (
	"fmt"
	"github.com/alist-org/alist/v3/pkg/utils"
)

const okCode = 0
const commonErrCode = 1

const jsonMarshalErr = `{"code":1,"err":"json marshal error: %s"}`

type result struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Err  string `json:"err"`
}

type account struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func resOk(data any) string {
	return res(okCode, data, "")
}

func resErr(err string) string {
	return resErrWithCode(commonErrCode, err)
}

func resErrWithCode(code int, err string) string {
	return res(code, nil, err)
}

func res(code int, data any, err string) string {
	res := result{Code: code, Data: data, Err: err}
	json, mErr := utils.Json.MarshalToString(res)
	if mErr != nil {
		return fmt.Sprintf(jsonMarshalErr, mErr.Error())
	}
	return json
}
