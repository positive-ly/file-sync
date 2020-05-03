package util

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"

	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	_server_file_info *_http_dir_struct
)

func _connect_http_discovery() {
	resp, err := http.Get(fmt.Sprintf("http://%s:%s/discovery", Http_server_addr, _httpProd))
	if err != nil {
		exit("[ERROR] request discovery fail : " + err.Error())
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		exit("[ERROR] read response body fail : " + err.Error())
	}
	if string(bytes) != _discovery_text {
		exit("[ERROR] Attestation fail")
	}
	if Log_print {
		fmt.Println(fmt.Sprintf("[INFO] %s 同步服务校验正常,准备连接同步服务...", Http_server_addr))
	}
}

func _connect_http_filelist() {
	resq, err := http.Get(fmt.Sprintf("http://%s:%s/list", Http_server_addr, _httpProd))
	if err != nil {
		exit("[ERROR]  request list fail : " + err.Error())
	}
	defer resq.Body.Close()
	_server_file_info = new(_http_dir_struct)
	err = json.NewDecoder(resq.Body).Decode(_server_file_info)
	if err != nil {
		exit("[ERROR]  Json Decoder fail : " + err.Error())
	}
	if Log_print {
		fmt.Println(fmt.Sprintf("[INFO] 连接 %s 同步服务成功，正在同步文件列表...", Http_server_addr))
	}
}

func ConnectTarServer() {
	if Http_server_addr == "" {
		return
	}
	if Target_path == "" {
		exit("[ERROR] -t 参数为空,未指定目的路径(被对比路径)")
	}
	_connect_http_discovery() //服务校验
	_connect_http_filelist()  //获取文件列表
	l_ds := osFileInfo_to_dirStruct(fileStat(Target_path))
	if _server_file_info.IsDir != l_ds.IsDir {
		exit("[ERROR] 两端文件类型不匹配，对比失败")
	}
	if _server_file_info.IsDir && l_ds.IsDir {
		_net_delete_copmare("", Target_path)
		_net_dir_comparison("", Target_path)
	} else {
		if _net_mod_time_and_size_compare(_server_file_info, l_ds) {
			creOrUpdFile(_server_file_info.FileName, Target_path, _server_file_info, _UPDATE_FILE)
		}
	}

}

func _net_delete_copmare(s, l string) {
	t_list, err := ioutil.ReadDir(l)
	if err != nil {
		exit(err.Error())
	}
	for _, v := range t_list {
		f_file := pathSeparator(s, v.Name())
		if http_file_isExist(f_file) {
			s_file := _http_file_info(f_file)
			if v.IsDir() == s_file.IsDir {
				if v.IsDir() && s_file.IsDir {
					_net_delete_copmare(f_file, pathSeparator(l, v.Name()))
				}
			} else {
				err_info_log(f_file, pathSeparator(l, v.Name()), s_file.IsDir, v.IsDir())
				continue
			}
		} else { //删除
			if Log_print {
				delete_info_log(v.Name(), l, v.IsDir())
			}
			if Delete_file {
				err := os.RemoveAll(pathSeparator(l, v.Name()))
				if err != nil {
					exit(err.Error())
				}
			}
		}
	}
}

func _net_dir_comparison(s, l string) {
	body := &request_dir{Dir: s}
	bodys, err := json.Marshal(body)
	if err != nil {
		exit("[ERROR]  Json Marshal fail : " + err.Error())
	}
	resp, err := http.Post(fmt.Sprintf("http://%s:%s/dirlist", Http_server_addr, _httpProd), "application/json;charset=UTF-8", bytes.NewReader(bodys))
	if err != nil {
		exit("[ERROR]  request dirlist fail : " + err.Error())
	}
	sds := make([]*_http_dir_struct, 0)
	err = json.NewDecoder(resp.Body).Decode(&sds)
	if err != nil {
		resp.Body.Close()
		exit("[ERROR]  Json Decode fail : " + err.Error())
	}
	resp.Body.Close()
	for _, v := range sds {
		l_file := pathSeparator(l, v.FileName)
		if fileIsExist(l_file) {
			t_file := osFileInfo_to_dirStruct(fileStat(l_file))
			if v.IsDir == t_file.IsDir {
				if v.IsDir && t_file.IsDir {
					_net_dir_comparison(pathSeparator(s, v.FileName), l_file)
				} else {
					if _net_mod_time_and_size_compare(v, t_file) {
						creOrUpdFile(pathSeparator(s, v.FileName), l_file, v, _UPDATE_FILE) //将源文件内容复制到目标文件
					}
				}
			} else {
				err_info_log(pathSeparator(s, v.FileName), l_file, v.IsDir, t_file.IsDir)
				continue
			}
		} else {
			if Log_print && v.IsDir {
				create_info_log(v.FileName, l, v.IsDir)
			}
			if v.IsDir {
				if Create_file {
					err := os.MkdirAll(l_file, os.ModePerm)
					if err != nil {
						exit(err.Error())
					}
					_net_dir_comparison(pathSeparator(s, v.FileName), l_file)
				}
			} else {
				creOrUpdFile(pathSeparator(s, v.FileName), l_file, v, _CREATE_FILE)
			}
		}
	}
}

func creOrUpdFile(s, l string, s_ds *_http_dir_struct, f_type file_type) {
	if Log_print {
		if f_type == _CREATE_FILE {
			create_info_log(s_ds.FileName, l, s_ds.IsDir)
		}
		if f_type == _UPDATE_FILE {
			copy_info_log(s, l)
		}
	}
	if f_type == _CREATE_FILE && Create_file {
		_net_source_copy_tatget(s, l, s_ds)
	}
	if f_type == _UPDATE_FILE && Copy_file {
		_net_source_copy_tatget(s, l, s_ds)
	}
}

func _net_source_copy_tatget(s, l string, s_ds *_http_dir_struct) {
fileDowon:
	// s_tm := time.Now()
	body := &request_dir{Dir: s}
	bodys, err := json.Marshal(body)
	if err != nil {
		exit("[ERROR]  Json Marshal fail : " + err.Error())
	}
	resp, err := http.Post(fmt.Sprintf("http://%s:%s/content", Http_server_addr, _httpProd), "application/json;charset=UTF-8", bytes.NewReader(bodys))
	if err != nil {
		exit("[ERROR]  request dirlist fail : " + err.Error())
	}
	l_f, err := os.Create(l)
	if err != nil {
		exit("[ERROR] OS Create fail : " + err.Error())
	}
	bufReader := bufio.NewReader(resp.Body)
	bufs := make([]byte, 4096)
	for {
		n, err := bufReader.Read(bufs)
		if err != nil && err != io.EOF {
			resp.Body.Close()
			l_f.Close()
			exit("[ERROR] bufio Read fail : " + err.Error())
		}
		if n == 0 {
			break
		}
		_, err1 := l_f.Write(bufs[:n])
		if err1 != nil && err1 != io.EOF {
			resp.Body.Close()
			l_f.Close()
			exit("[ERROR] File Write fail : " + err1.Error())
		}
		// if err == io.EOF { //最后一次读取，有值，但err也是EOF。所以需要额外判断！！！不判断就多循环一次
		// 	break
		// }
	}
	resp.Body.Close()
	l_f.Close()
	/*
		// bs, err := ioutil.ReadAll(resp.Body)
		// if err != nil {
		// 	resp.Body.Close()
		// 	exit("[ERROR] Read All fail : " + err.Error())
		// }
		// resp.Body.Close()
		// l_f, err := os.Create(l)
		// if err != nil {
		// 	exit("[ERROR] OS Create fail : " + err.Error())
		// }
		// _, err = l_f.Write(bs)
		// if err != nil {
		// 	l_f.Close()
		// 	exit("[ERROR] File Write fail : " + err.Error())
		// }
		// l_f.Close()
	*/
	tm := time.Unix(s_ds.ModTime, 0)
	os.Chtimes(l, tm, tm)
	if Http_file_md5 {
		//校验文件MD5值
		// if Log_print {
		// 	fmt.Println("[INFO] 校验文件MD5值,文件路径 : ", l)
		// }
		ok, err := local_server_file_commMd5(s, l)
		if err != nil {
			exit("[ERROR] comparison File MD5 fail : " + err.Error())
		}
		if !ok {
			// fmt.Println("[ERROR] 校验文件MD5值失败,文件路径 : ", l)
			goto fileDowon
		}
		// fmt.Println("[SUCCESS] 校验文件MD5值成功,文件路径 : ", l)
	}
	// fmt.Println("*************:", time.Since(s_tm).String())
}

func _net_file_md5(s string) (*request_file_md5, error) {
	body := &request_dir{Dir: s}
	bodys, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s:%s/filemd5", Http_server_addr, _httpProd), "application/json;charset=UTF-8", bytes.NewReader(bodys))
	if err != nil {
		return nil, err
	}
	m := new(request_file_md5)
	err = json.NewDecoder(resp.Body).Decode(m)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return m, nil
}

func local_server_file_commMd5(s, l string) (bool, error) {
	s_md5, err := _net_file_md5(s)
	if err != nil {
		return false, err
	}
	l_md5, err := file_md5(l)
	if err != nil {
		return false, err
	}
	// fmt.Println("-----------------")
	// fmt.Println(s_md5.Md5)
	// fmt.Println("-----------------")
	// fmt.Println(l_md5)
	// fmt.Println("-----------------")
	// fmt.Println(s_md5.Md5 == l_md5)
	// fmt.Println("-----------------")
	return s_md5.Md5 == l_md5, nil
}

func _net_file_isExist(s string) (*request_file_isExist, error) {
	body := &request_dir{Dir: s}
	bodys, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s:%s/isExist", Http_server_addr, _httpProd), "application/json;charset=UTF-8", bytes.NewReader(bodys))
	if err != nil {
		return nil, err
	}
	m := new(request_file_isExist)
	err = json.NewDecoder(resp.Body).Decode(m)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return m, nil
}

func http_file_isExist(s string) bool {
	exist, err := _net_file_isExist(s)
	if err != nil {
		exit("[ERROR] _net_file_isExist fail : " + err.Error())
	}
	return exist.Exist
}

func _http_file_info(s string) *_http_dir_struct {
	body := &request_dir{Dir: s}
	bodys, err := json.Marshal(body)
	if err != nil {
		exit("[ERROR]  Json Marshal fail : " + err.Error())
	}
	resp, err := http.Post(fmt.Sprintf("http://%s:%s/fileinfo", Http_server_addr, _httpProd), "application/json;charset=UTF-8", bytes.NewReader(bodys))
	if err != nil {
		exit("[ERROR]  request dirlist fail : " + err.Error())
	}
	sds := new(_http_dir_struct)
	err = json.NewDecoder(resp.Body).Decode(sds)
	if err != nil {
		resp.Body.Close()
		exit("[ERROR]  Json Decode fail : " + err.Error())
	}
	resp.Body.Close()
	return sds
}

func _net_mod_time_and_size_compare(s, l *_http_dir_struct) bool {
	return !(s.ModTime == l.ModTime && s.Size == l.Size)
}
