package main

import "C"
import (
	"github.com/alist-org/alist/v3/jni"
)

//export initConfig
func initConfig(in *C.char) *C.char {
	var args = C.GoString(in)
	println("initConfig: ", args)
	var result = jni.InitConfig(args)
	return C.CString(result)
}

//export startServer
func startServer(in *C.char) *C.char {
	var args = C.GoString(in)
	println("startServer")
	var result = jni.StartServer(args)
	return C.CString(result)
}

//export stopServer
func stopServer() *C.char {
	println("stopServer")
	var result = jni.StopServer()
	return C.CString(result)
}

//export getAdmin
func getAdmin() *C.char {
	var result = jni.GetAdmin()
	return C.CString(result)
}

//export listFile
func listFile(in *C.char) *C.char {
	var args = C.GoString(in)
	println("listFile: ", args)
	var result = jni.ListFile(args)
	return C.CString(result)
}

func main() {
	jni.InitConfig(`{"dataDir": "data", "logStd": true}`)
}
