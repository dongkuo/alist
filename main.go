package main

import "C"
import (
	"github.com/alist-org/alist/v3/jni"
	"os"
	"os/signal"
	"syscall"
)

//export initConfig
func initConfig(in *C.char) *C.char {
	var args = C.GoString(in)
	println("initConfig: ", args)
	var result = jni.InitConfig(args)
	return C.CString(result)
}

//export startServer
func startServer() *C.char {
	println("startServer")
	var result = jni.StartServer()
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
	println("getAdmin")
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
	jni.InitConfig(`{"data_dir": "data", "log_std": true}`)
	jni.GetAdmin()
	//var files = jni.ListFile(`{"page":1,"path":"/","per_page":50,"refresh":false}`)
	//println(files)
	jni.StartServer()
	var quit = make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	//cmd.Execute()
}
