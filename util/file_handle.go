package util

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type file_type int

const (
	_CREATE_FILE file_type = iota
	_UPDATE_FILE
)

func Contrast(s, t string) {
	if s == "" || t == "" {
		return
		exit("请使用 \"-s\" \"-t\" 参数指定文件/文件夹")
	}
	s_fi, t_fi := fileStat(s), fileStat(t)
	if s_fi.IsDir() != t_fi.IsDir() {
		exit("对比类型不匹配")
	}
	if s_fi.IsDir() && t_fi.IsDir() { //都是文件夹
		delete_copmare(s, t)
		increment_compare(s, t)
	} else { //都是文件
		if mod_time_and_size_compare(s_fi, t_fi) {
			create_and_copy_file(s, t, s_fi, _UPDATE_FILE) //将源文件内容复制到目标文件
		}
	}
	// end_info_log()
}

//删量对比操作
func delete_copmare(s, t string) {
	t_list, err := ioutil.ReadDir(t)
	if err != nil {
		exit(err.Error())
	}
	for _, v := range t_list {
		f_file := pathSeparator(s, v.Name())
		if fileIsExist(f_file) {
			s_file := fileStat(f_file)
			if v.IsDir() == s_file.IsDir() {
				if v.IsDir() && s_file.IsDir() {
					delete_copmare(f_file, pathSeparator(t, v.Name()))
				}
			} else {
				err_info_log(f_file, pathSeparator(t, v.Name()), s_file.IsDir(), v.IsDir())
				continue
			}
		} else { //删除
			if Log_print {
				delete_info_log(f_file, v.IsDir())
			}
			if Delete_file {
				err := os.RemoveAll(pathSeparator(t, v.Name()))
				if err != nil {
					exit(err.Error())
				}
			}
		}
	}
}

//增量对比操作
func increment_compare(s, t string) {
	s_list, err := ioutil.ReadDir(s)
	if err != nil {
		exit(err.Error())
	}
	for _, v := range s_list {
		f_file := pathSeparator(t, v.Name())
		if fileIsExist(f_file) {
			t_file := fileStat(f_file)
			if v.IsDir() == t_file.IsDir() {
				if v.IsDir() && t_file.IsDir() {
					increment_compare(pathSeparator(s, v.Name()), f_file)
				} else {
					if mod_time_and_size_compare(v, t_file) {
						create_and_copy_file(pathSeparator(s, v.Name()), f_file, v, _UPDATE_FILE) //将源文件内容复制到目标文件
					}
				}
			} else {
				err_info_log(pathSeparator(s, v.Name()), f_file, v.IsDir(), t_file.IsDir())
				continue
			}
		} else {
			if Log_print {
				create_info_log(f_file, v.IsDir())
			}
			if v.IsDir() {
				if Create_file {
					err := os.MkdirAll(f_file, os.ModePerm)
					if err != nil {
						exit(err.Error())
					}
					increment_compare(pathSeparator(s, v.Name()), f_file)
				}
			} else {
				create_and_copy_file(pathSeparator(s, v.Name()), f_file, v, _CREATE_FILE)
			}
		}
	}
}

//创建或拷贝文件内容
func create_and_copy_file(s, t string, s_fi os.FileInfo, f_type file_type) {
	if f_type == _CREATE_FILE && Create_file {
		source_copy_tatget(s, t, s_fi)
	}
	if f_type == _UPDATE_FILE && Copy_file {
		source_copy_tatget(s, t, s_fi)
	}
}

//源文件拷贝到目标文件相同位置
func source_copy_tatget(s, t string, s_fi os.FileInfo) {
	if Log_print {
		copy_info_log(s, t)
	}
	t_f, err := os.Create(t)
	if err != nil {
		exit(err.Error())
	}
	s_f, err := os.Open(s)
	if err != nil {
		t_f.Close()
		exit(err.Error())
	}
	buffs := make([]byte, 4096)
	for {
		n, err := s_f.Read(buffs)
		if err != nil && err != io.EOF {
			t_f.Close()
			s_f.Close()
			exit(err.Error())
		}
		if n == 0 || err == io.EOF {
			break
		}
		_, err = t_f.Write(buffs[:n])
		if err != nil {
			t_f.Close()
			s_f.Close()
			exit(err.Error())
		}
	}
	t_f.Close()
	s_f.Close()
	os.Chtimes(t, s_fi.ModTime(), s_fi.ModTime())
}

//文件修改时间与文件大小对比
func mod_time_and_size_compare(s_fi, t_fi os.FileInfo) bool {
	return !(s_fi.ModTime().Equal(t_fi.ModTime()) && s_fi.Size() == t_fi.Size())
}

//获取文件状态信息
func fileStat(path string) os.FileInfo {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			exit(path + " -- 文件不存在")
		}
		exit(err.Error())
	}
	return fi
}

//判断文件/文件夹是否存在
func fileIsExist(f string) bool {
	_, err := os.Stat(f)
	return !os.IsNotExist(err)
}

//路径拼接
func pathSeparator(path, f_name string) string {
	if path == "" {
		return f_name
	}
	if string(filepath.Separator) == "/" {
		f_name = strings.Replace(f_name, "\\", "/", -1)
	}
	if string(filepath.Separator) == "\\" {
		f_name = strings.Replace(f_name, "/", "\\", -1)
	}
	return (path + string(filepath.Separator) + f_name)
}
