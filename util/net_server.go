package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	_http_dir_info os.FileInfo
)

func OpenHttpListen() {
	if !Http_Open {
		return
	}
	if Source_path == "" {
		exit("-s 参数为空,未指定源路径")
	}
	_http_dir_info = fileStat(Source_path) //获取源路径信息
	//开启监听服务
	http.HandleFunc("/list", _handle_filesList)
	http.HandleFunc("/discovery", _handle_discovery)
	http.HandleFunc("/content", _handle_get_file_data)
	http.HandleFunc("/dirlist", _handle_dir_list)
	http.HandleFunc("/filemd5", _handle_file_md5)
	http.HandleFunc("/isExist", _handle_file_exist)
	http.HandleFunc("/fileinfo", _handle_get_file_info)
	go func() {
		time.Sleep(200 * time.Millisecond)
		fmt.Println("[INFO] 文件同步服务已启动,源路径为 : ", Source_path)
	}()
	err := http.ListenAndServe(":"+_httpProd, nil)
	if err != nil {
		exit("ListenAndServe error: " + err.Error())
	}
}

func _handle_discovery(w http.ResponseWriter, r *http.Request) {
	r.Body.Close()
	w.Write([]byte(_discovery_text))
}

func _handle_filesList(w http.ResponseWriter, r *http.Request) {
	r.Body.Close()
	bytes, err := json.Marshal(osFileInfo_to_dirStruct(_http_dir_info))
	if err != nil {
		fmt.Println("[ERROR] json marshal fail : ", err.Error())
		return
	}
	w.Write(bytes)
}

func _handle_dir_list(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		fmt.Println("[ERROR] Read All fail : ", err.Error())
		return
	}
	dirstu := new(request_dir)
	err = json.Unmarshal(bytes, dirstu)
	if err != nil {
		fmt.Println("[ERROR] Json Unmarshal fail : ", err.Error())
		return
	}
	if dirstu.Dir != "" && _http_dir_info.IsDir() {
		dirstu.Dir = pathSeparator(Source_path, dirstu.Dir)
	} else {
		dirstu.Dir = Source_path
	}
	// fmt.Println("=========", dirstu.Dir)
	if !fileIsExist(dirstu.Dir) {
		fmt.Println("[ERROR] file not exist")
		return
	}
	hds := make([]*_http_dir_struct, 0)
	fi := fileStat(dirstu.Dir)
	if fi.IsDir() {
		fis, err := ioutil.ReadDir(dirstu.Dir)
		if err != nil {
			fmt.Println("[ERROR] Read Dir fail : ", err.Error())
			return
		}
		for _, v := range fis {
			hds = append(hds, osFileInfo_to_dirStruct(v))
		}
	} else {
		hds = append(hds, osFileInfo_to_dirStruct(fi))
	}
	bytes, err = json.Marshal(hds)
	if err != nil {
		fmt.Println("[ERROR] Json Marshal fail : ", err.Error())
		return
	}
	w.Write(bytes)
}

func _handle_get_file_data(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		fmt.Println("[ERROR] Read All fail : ", err.Error())
		return
	}
	dirstu := new(request_dir)
	err = json.Unmarshal(bytes, dirstu)
	if err != nil {
		fmt.Println("[ERROR] Json Unmarshal fail : ", err.Error())
		return
	}
	if dirstu.Dir != "" && _http_dir_info.IsDir() {
		dirstu.Dir = pathSeparator(Source_path, dirstu.Dir)
	} else {
		dirstu.Dir = Source_path
	}
	if Log_print {
		fmt.Println(r.Host, " 获取 ", dirstu.Dir) ////////////////////////
	}
	if !fileIsExist(dirstu.Dir) {
		fmt.Println("[ERROR] file not exist : " + dirstu.Dir)
		return
	}
	bytes, err = ioutil.ReadFile(dirstu.Dir)
	if err != nil {
		fmt.Println("[ERROR] Read file fail : ", err.Error())
		return
	}

	w.Write(bytes)
}

func _handle_file_md5(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		fmt.Println("[ERROR] Read All fail : ", err.Error())
		return
	}
	dirstu := new(request_dir)
	err = json.Unmarshal(bytes, dirstu)
	if err != nil {
		fmt.Println("[ERROR] Json Unmarshal fail : ", err.Error())
		return
	}
	if dirstu.Dir != "" && _http_dir_info.IsDir() {
		dirstu.Dir = pathSeparator(Source_path, dirstu.Dir)
	} else {
		dirstu.Dir = Source_path
	}
	m, err := file_md5(dirstu.Dir)
	if err != nil {
		fmt.Println("[ERROR] File Md5 fail : ", err.Error())
		return
	}
	bytes, err = json.Marshal(request_file_md5{Md5: m})
	if err != nil {
		fmt.Println("[ERROR] Json Marshal fail : ", err.Error())
		return
	}
	w.Write(bytes)
}

func _handle_file_exist(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		fmt.Println("[ERROR] Read All fail : ", err.Error())
		return
	}
	dirstu := new(request_dir)
	err = json.Unmarshal(bytes, dirstu)
	if err != nil {
		fmt.Println("[ERROR] Json Unmarshal fail : ", err.Error())
		return
	}
	if dirstu.Dir != "" && _http_dir_info.IsDir() {
		dirstu.Dir = pathSeparator(Source_path, dirstu.Dir)
	} else {
		dirstu.Dir = Source_path
	}

	bytes, err = json.Marshal(request_file_isExist{Exist: fileIsExist(dirstu.Dir)})
	if err != nil {
		fmt.Println("[ERROR] Json Marshal fail : ", err.Error())
		return
	}
	w.Write(bytes)
}

func _handle_get_file_info(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		fmt.Println("[ERROR] Read All fail : ", err.Error())
		return
	}
	dirstu := new(request_dir)
	err = json.Unmarshal(bytes, dirstu)
	if err != nil {
		fmt.Println("[ERROR] Json Unmarshal fail : ", err.Error())
		return
	}
	if dirstu.Dir != "" && _http_dir_info.IsDir() {
		dirstu.Dir = pathSeparator(Source_path, dirstu.Dir)
	} else {
		dirstu.Dir = Source_path
	}
	if !fileIsExist(dirstu.Dir) {
		fmt.Println("[ERROR] file not exist")
		return
	}
	bytes, err = json.Marshal(osFileInfo_to_dirStruct(fileStat(dirstu.Dir)))
	if err != nil {
		fmt.Println("[ERROR] Json Marshal fail : ", err.Error())
		return
	}
	w.Write(bytes)
}
