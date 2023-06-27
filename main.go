package main

import "C"
import (
	"github.com/alist-org/alist/v3/jni"
)

//export execCommand
func execCommand(in *C.char) *C.char {
	var args = C.GoString(in)
	var result = jni.ExecCommand(args)
	return C.CString(result)
}

//export getAdmin
func getAdmin(in *C.char) *C.char {
	var dataDir = C.GoString(in)
	var result = jni.GetAdmin(dataDir)
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
	var json = jni.ListFile(`{"page":1,"path":"/","per_page":50,"refresh":false}`)
	println(json)
}
