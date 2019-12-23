package util

import (
	"fmt"
)

func end_info_log() {
	fmt.Println("文件对比完成")
}

func create_info_log(str string, isDir bool) {
	s := ""
	if isDir {
		s = "文件夹"
	} else {
		s = "文件"
	}
	fmt.Println(fmt.Sprintf("[+]\t创建处理: %-3s \t[ %s ]", s, str))
}

func delete_info_log(str string, isDir bool) {
	s := ""
	if isDir {
		s = "文件夹"
	} else {
		s = "文件"
	}
	fmt.Println(fmt.Sprintf("[-]\t删除处理: %-3s \t[ %s ]", s, str))
}

func copy_info_log(s, t string) {
	fmt.Println(fmt.Sprintf("[+]\t拷贝处理: [ %s ] 到 [ %s ]", s, t))
}

var errs []string

func err_info_log(s_str, t_str string, s_isDir, t_isDir bool) {
	for _, v := range errs {
		if v == t_str {
			return
		}
	}
	errs = append(errs, t_str)
	var s, t string
	if s_isDir {
		s = "文件夹"
	} else {
		s = "文件"
	}
	if t_isDir {
		t = "文件夹"
	} else {
		t = "文件"
	}

	fmt.Println(fmt.Sprintf("[E]\t处理失败: [-s %s] [-t %s] %s  %s", s, t, s_str, t_str))
}
