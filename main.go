package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// a.1 实现读取文件handler
	fileHandler := http.FileServer(http.Dir("./video"))

	// a.2 注册handler
	http.Handle("/video/", http.StripPrefix("/video/", fileHandler))

	// 注册上传文件的handler
	http.HandleFunc("/api/upload", uploadHandler)

	http.HandleFunc("/api/list", getFileListHandler)

	http.HandleFunc("/sayHello", sayHello)
	http.ListenAndServe(":8090", nil)
}

// 1.业务逻辑
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 1.限制客户端上传视频文件的大小
	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)
	err := r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 2.获取上传的文件
	file, fileHeader, err := r.FormFile("uploadFile")
	defer file.Close()

	// 3.检查文件类型
	ret := strings.HasSuffix(fileHeader.Filename, ".mp4")
	if ret == false {
		http.Error(w, "not mp4", http.StatusInternalServerError)
		return
	}

	// 4.获取随机名称
	md5Byte := md5.Sum([]byte(fileHeader.Filename + time.Now().String()))
	md5Str := fmt.Sprintf("%x", md5Byte)
	newFileName := md5Str + ".mp4"

	// 5.写入文件
	dst, err := os.Create("./video/" + newFileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

// 获取视频文件列表
func getFileListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	files, _ := filepath.Glob("video/*")
	var ret []string
	for _, file := range files {
		ret = append(ret, "http://"+r.Host+"/video/"+filepath.Base(file))
	}
	retJSON, _ := json.Marshal(ret)
	w.Write(retJSON)
	return
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello worl"))
}
