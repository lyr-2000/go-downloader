package downloader

import (
	"fmt"
	"strings"
)

func fileName(path string) string {
	index := strings.LastIndex(path, ".")
	pathSize := len(path)
	//if index==0 {
	//	// 形如 ./wox.exe
	//	// -> 改为  wox.exe
	//	strings.LastIndex("")
	//}
	if index< pathSize || index>=pathSize {
		idx := strings.LastIndex(path, "/")
		if idx>=0 && idx<pathSize {
			return path[idx+1:]
		}
		return "res"
	}
	return path[index:]
}
//索引文件名字，用于后面写入
func IndexName(fName string,workCnt int,id int) string {
	return fmt.Sprintf("index_%v_workCnt_%v_%v",fName,workCnt,id)
}

func TmpFilePath(tmpPath string ,fName string ,workCnt int, id int ) string {
	return tmpPath+"/"+ IndexName(fName,workCnt,id)
}
//func TmpPath