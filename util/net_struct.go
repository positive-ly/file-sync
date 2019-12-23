package util

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

const (
	_httpProd       = "28888"
	_discovery_text = "file-sync"
)

type _http_dir_struct struct {
	FileName string `json:"file_name"`
	IsDir    bool   `json:"is_dir"`
	ModTime  int64  `json:"mod_time"`
	Size     int64  `json:"size"`
}

type request_dir struct {
	Dir string `json:"dir"`
}

type request_file_md5 struct {
	Md5 string `json:"md5"`
}

type request_file_isExist struct {
	Exist bool `json:"exist"`
}

func osFileInfo_to_dirStruct(fi os.FileInfo) *_http_dir_struct {
	return &_http_dir_struct{
		FileName: fi.Name(), IsDir: fi.IsDir(), ModTime: fi.ModTime().Unix(), Size: fi.Size(),
	}
}

func file_md5(f string) (string, error) {
	file, err := os.Open(f)
	if err != nil {
		return "", err
	}
	md5hash := md5.New()
	if _, err := io.Copy(md5hash, file); err != nil {
		file.Close()
		return "", err
	}
	file.Close()
	return fmt.Sprintf("%x", md5hash.Sum(nil)), nil
}
