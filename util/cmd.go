package util

import (
	"flag"
	"fmt"
	"os"
)

var (
	Source_path string //源路径
	Target_path string //目的路径
	All_file    bool   //创建、拷贝、删除
	Create_file bool   //创建不存在文件
	Copy_file   bool   //拷贝不相同文件
	Delete_file bool   //删除未对比文件
	Log_print   bool   //日志打印
	//--------------------
	Http_Open        bool   //开启Http监听端口
	Http_server_addr string //http远端服务地址
	Http_file_md5    bool   //校验下载文件MD5值
	Scan_lan_server  bool   //扫描本机网卡局域网内文件同步服务
	//--------------------
)

func CmdParams() {
	flag.StringVar(&Source_path, "s", "", "[***] 必填 : 源路径（根据此项路径对比 t 参数路径内文件/文件夹内容）--source")
	flag.StringVar(&Target_path, "t", "", "[***] 必填 : 目的路径(被对比路径) --target")
	flag.BoolVar(&Log_print, "list", false, "打印日志信息 --list")
	flag.BoolVar(&All_file, "all", false, "文件对比时执行全部可执行操作（相当于 \"-c\" \"-u\" \"-d\" 参数同时使用）--all")
	flag.BoolVar(&Create_file, "c", false, "创建不存在文件 --create")
	flag.BoolVar(&Copy_file, "u", false, "更新同名文件内容 --update")
	flag.BoolVar(&Delete_file, "d", false, "删除不存在文件 --delete")
	//----------------------------------------------------------------------------------
	flag.BoolVar(&Http_Open, "server", false, "开启Http监听 --server")
	flag.StringVar(&Http_server_addr, "server-addr", "", "远端http服务IP地址 --server-addr")
	flag.BoolVar(&Http_file_md5, "md5", false, "基于同步文件时,将会校验文件MD5，但会增加同步时间 --md5")
	flag.BoolVar(&Scan_lan_server, "scan", false, "扫描本机网卡局域网内文件同步服务 --scan")
	//----------------------------------------------------------------------------------
	flag.Parse()
	if All_file {
		Create_file, Copy_file, Delete_file = true, true, true
	}
}

func exit(str string) {
	// red := string([]byte{27, 91, 57, 49, 109})
	// reset := string([]byte{27, 91, 48, 109})
	// fmt.Print(red, str, reset)
	fmt.Print(str)
	os.Exit(2)
}
