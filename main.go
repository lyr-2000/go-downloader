package main

import (
	"errors"
	"flag"
	"fmt"
	"goDownloader/downloader"
	"log"
	"strings"
	"time"
)
//type FileSize int
//var (
//	k FileSize = 1024
//	M FileSize = 1024*1024
//)
//var (
//	k = 1024
//	M = k*1024
//
//)
/*
项目介绍：
 多线程下载工具
 模仿IDM 32线程下载

*/
//var defaultTimeout = time.Second*25
func main() {
	//解析输出参数
	config,err := parseFlag()
	if err!=nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(&config)
	//预处理工作
	config.PreHandle()

	//预备工作

	//创建对象
	if config.Slice== false{
		//非断点下载
		d := downloader.NewDefault()

		//正式下载
		for {
			err = d.Download(config.WorkerCnt,config.Url,
				config.DownLoadPath+fName(config.Url),
				config.CacheSize,
				config.Proxy,
				int64(time.Second)* int64(config.Timeout) ,
			)
			if err == nil {
				break
			}
			log.Printf("下载出现异常 := %v\n",err)
			if !config.Retry {
				break
			}
			log.Printf("下载重试\n")

		}
	}else {
		//使用断点下载的方式进行下载
		d := downloader.NewSliceDownloader(config.TmpPath)

		//正式下载
		for {
			err = d.Download(config.WorkerCnt,config.Url,
				config.DownLoadPath+fName(config.Url),
				config.CacheSize,
				config.Proxy,
				int64(time.Second)* int64(config.Timeout) ,
			)
			if err == nil {
				break
			}
			log.Printf("下载出现异常 := %v\n",err)
			if !config.Retry {
				break
			}
			log.Printf("下载重试\n")

		}
	}
}


//下载配置
type DownloadConfig struct {
	WorkerCnt int //工作线程数
	Url string //下载地址
	DownLoadPath string //下载地址
	TmpPath string // 临时文件目录
	CacheSize int //缓冲区大小

	Proxy string
	Retry bool //失败重试


	Timeout int //超时时间

	Slice bool //分片下载，支持断点下载





}
func (w *DownloadConfig) String() string {
	return fmt.Sprintf("url:=%v," +
		"\nworkCnt:=%v\n, downloadPath:=%v\n" +
		"tmpPath:=%v" +
		"",w.Url,w.WorkerCnt,w.DownLoadPath,w.TmpPath)
}


func (w *DownloadConfig)PreHandle() {

}

// 解析参数
func parseFlag() (config DownloadConfig,err error){
	//解析输入 url
	flag.StringVar(&config.Url,"url","","url不能为空，请设置 -url 指定")
	flag.StringVar(&config.DownLoadPath,"path","./","文件下载路径")
	flag.StringVar(&config.TmpPath,"tmp","./tmp","下载的临时文件")
	flag.IntVar(&config.CacheSize,"buf",1024*100,"缓冲区大小【单位：字节】")

	flag.IntVar(&config.WorkerCnt,"c",8,"下载线程数")
	flag.StringVar(&config.Proxy,"proxy","","设置代理,【本人没用过】：例如:http://127.0.0.1:7890")
	flag.BoolVar(&config.Retry,"retry",true,"失败重试,")
	//默认超时时间为  25s
	flag.IntVar(&config.Timeout,"s",25,"设置超时时间，单位为秒")
	flag.BoolVar(&config.Slice,"bd",true,"断点下载,breakpoint downloading")

	//解析
	flag.Parse()
	if config.Url=="" {
		return config, errors.New("没有指定下载链接,请输入 -path")
	}
	if config.WorkerCnt<=0 {
		return config,errors.New("输入的工作线程数不合法,请合法设置 -workCnt")
	}
	if config.WorkerCnt > 60 {
		log.Printf("warn: workerCnt > 60")
	}

	return config,err
}







//获取文件名字
func fName(path string) string {
	index := strings.LastIndex(path, ".")
	pathSize := len(path)
	if index< pathSize || index>=pathSize {
		idx := strings.LastIndex(path, "/")
		if idx>=0 && idx<pathSize {
			return path[idx+1:]
		}
		return "result"
	}
	return path[index:]

}


