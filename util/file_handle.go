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
//参数: 源路径，目标路径
func delete_copmare(s, t string, layNum ...int) {
	l_n := 1
	if len(layNum) > 0 {
		l_n = layNum[0]
	}
	//获取被对比文件夹路径
	t_list, err := ioutil.ReadDir(t)
	if err != nil {
		exit(err.Error())
	}
	for _, v := range t_list { //目标路径文件遍历
		f_file := pathSeparator(s, v.Name()) //源路径文件地址拼接
		if fileIsExist(f_file) {             //如果源路径相同文件/文件夹存在
			s_file := fileStat(f_file)       //获取源路径同名文件信息
			if v.IsDir() == s_file.IsDir() { //判断文件夹是否相同类型
				if v.IsDir() && s_file.IsDir() { //如果都是文件夹，
					delete_copmare(f_file, pathSeparator(t, v.Name()), l_n+1) //递归搜索
				}
			} else { //类型不同报错提示
				err_info_log(f_file, pathSeparator(t, v.Name()), s_file.IsDir(), v.IsDir())
				continue
			}
		} else { //删除目标路径文件
			if Layer_num > l_n { //如果当前层级小于等于预设层级，则跳过删除操作
				continue
			}
			if Log_print {
				delete_info_log(v.Name(), t, v.IsDir())
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
func increment_compare(s, t string, layNum ...int) {
	l_n := 1
	if len(layNum) > 0 {
		l_n = layNum[0]
	}
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
					increment_compare(pathSeparator(s, v.Name()), f_file, l_n+1)
				} else {
					if Layer_num > l_n { //如果当前层级小于等于预设层级，则跳过
						continue
					}
					if mod_time_and_size_compare(v, t_file) {
						if Log_print {
							copy_info_log(v.Name(), t)
						}
						if Copy_file {
							create_and_copy_file(pathSeparator(s, v.Name()), f_file, v, _UPDATE_FILE) //将源文件内容复制到目标文件
						}
					}
				}
			} else {
				err_info_log(pathSeparator(s, v.Name()), f_file, v.IsDir(), t_file.IsDir())
				continue
			}
		} else {
			if Layer_num > l_n { //如果当前层级小于等于预设层级，则跳过
				continue
			}
			if v.IsDir() {
				if Log_print {
					create_info_log(v.Name(), t, v.IsDir())
				}
				if Create_file {
					err := os.MkdirAll(f_file, os.ModePerm)
					if err != nil {
						exit(err.Error())
					}
					increment_compare(pathSeparator(s, v.Name()), f_file, l_n+1)
				}
			} else {
				if Log_print {
					create_info_log(v.Name(), t, v.IsDir())
				}
				if Copy_file {
					create_and_copy_file(pathSeparator(s, v.Name()), f_file, v, _CREATE_FILE) //将源文件内容复制到目标文件
				}
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
