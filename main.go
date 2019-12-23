package main

import (
	// "file-sync/test"
	"file-sync/util"
)

func main() {
	//获取命令行参数
	util.CmdParams()
	//对比文件
	util.Contrast(util.Source_path, util.Target_path)
	//同步服务服务端
	util.OpenHttpListen()
	//同步服务接收端
	util.ConnectTarServer()
	//扫描局域网内同步服务服务端
	util.LocalIpAddCIDR()
	//-------------------
	// test.LocalIpAddress()
}
