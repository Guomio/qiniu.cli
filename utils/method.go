package utils

import (
	"log"
	"path/filepath"
	"strconv"
)

// GetAbsPath 获取绝对路径，获取失败则退出进程
func GetAbsPath(s string) string {
	dir, err := filepath.Abs(s)
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

// GetFileSize 获取文件大小
func GetFileSize(n int64) string {
	table := []string{"B", "KB", "MB", "GB", "GB", "PB", "EB", "ZB", "YB"}
	i := 0
	for n>>10 >= 1 {
		n >>= 10
		i++
	}
	if i > len(table) {
		return "too big!!!"
	}
	return strconv.FormatInt(n, 10) + " " + table[i]
}
