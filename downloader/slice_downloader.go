package downloader

import (
	"fmt"
	"log"
	"os"
)

//这个 downloader 会写到 临时文件里面， 下载结束之后 再聚合到新的文件里面


////下载对象
//type Downloader interface {
//	Download(	WorkerCnt int  , Url string , DownLoadPath string , bufSize int ,proxyUrl string,timeout int64 ) error
//	//DownloadWithProxy(	WorkerCnt int  , Url string , DownLoadPath string ,proxy string, bufSize int ) error
//
//
//}
type SliceDownloaderImpl struct {
	TmpPath string //临时索引文件目录
}


func NewSliceDownloader(tmpPath string) Downloader {
	//
	w := SliceDownloaderImpl{
		TmpPath: tmpPath,
	}
	_,err := os.Stat(tmpPath)
	if os.IsNotExist(err) {
		os.Mkdir(tmpPath,os.ModeDir)
		//os.Chmod(tmpPath,os.mod)
	}

	return &w;
}
func (w *SliceDownloaderImpl) Download(	WorkerCnt int  ,	//协程数量
	Url string ,			//网页链接
	DownLoadPath string ,	//写入的文件路径
	bufSize int ,		//单个协程分配的缓冲区大小
	proxyUrl string, // 下载的代理

	timeout int64,//超时时间

) error {
	fmt.Println(DownLoadPath)
	f, err := os.Create(DownLoadPath)
	defer f.Close()
	if err!=nil {
		log.Fatalf(err.Error())
	}
	//defer f.Close()
	task := FileTask{
		File: f,
		Url: Url,
		BufferSize: bufSize,
		WorkCnt: int64(WorkerCnt),
		ProxyUrl: proxyUrl,
		TmpPath: w.TmpPath,
	}
	var localSize int64
	stat, err := os.Stat(DownLoadPath)
	if err!=nil {

		return err

	}
	localSize =  stat.Size()
	//获取文件大小
	if task.Size <= 0 {
		// 获取文件信息
		req, err := NewGetRequestObject(task.Url)
		resp,err := NewRequestClient(proxyUrl,task.Timeout).Do(req)
		if err != nil {
			//log.Printf("请求结束 %v",err)
			return  err
		}


		task.Size = int64(resp.ContentLength)
		if task.Size == localSize {
			log.Printf("已经下载完成 ,文件大小 %v ~~ \n",task.Size)
			_ = resp.Body.Close()
			return nil
		}
		defer resp.Body.Close()
	}
	log.Printf("\n文件大小 := %v M,工作线程数量 := %v,分块大小 := %v M\n",task.Size/1024/1024,task.WorkCnt,task.Size/1024/1024/task.WorkCnt)
	//下载
	task.BeginDownloadInTmp()
	return nil
}











