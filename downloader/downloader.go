package downloader

import (
	"fmt"
	"github.com/redmask-hb/GoSimplePrint/goPrint"
	"io"
	"log"
	"os"
)
//下载对象
type Downloader interface {
	  Download(	WorkerCnt int  , Url string , DownLoadPath string , bufSize int ,proxyUrl string,timeout int64 ) error
	  //DownloadWithProxy(	WorkerCnt int  , Url string , DownLoadPath string ,proxy string, bufSize int ) error


}
type DownloaderImpl struct {
}

//工厂模式返回对象
func  NewDefault() Downloader {
	return &DownloaderImpl{}
}


//func NewDownLoader(w *DownLoadCallback) Downloader{
//	return &DownloaderImpl{}
//}
//下载
func (w *DownloaderImpl) Download(	WorkerCnt int  ,	//协程数量
								Url string ,			//网页链接
								DownLoadPath string ,	//写入的文件路径
								//TmpPath string ,
								bufSize int ,		//单个协程分配的缓冲区大小
								proxyUrl string, // 下载的代理

								timeout int64,//超时时间
) error{
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
	log.Printf("文件大小 := %v M,工作线程数量 := %v,分块大小 := %v M\n",task.Size/1024/1024,task.WorkCnt,task.Size/1024/1024/task.WorkCnt)

	task.Begin()
	return nil
}




//文件切片
type FileSlice struct {
	Start int64
	End int64


}
func (w *FileSlice) String() string {
	return fmt.Sprintf("range[%v,%v]",w.Start,w.End)
}

type FileTask struct {
	Url string //下载
	ProxyUrl string //下载的代理url
	Size int64  //文件大小
	File *os.File //写入的文件引用
	FileSlice []FileSlice //切片任务
	//Speeds int64 //下载速度
	WorkCnt int64 //下载线程数量
	DownloadSize int64 //已经下载的大小

	BufferSize int  //每个任务的缓冲区大小

	//MD5 string //文件签名 =  MD5(文件名字+ 修改时间)
	TmpPath string
	Timeout int64
	Bar *goPrint.Bar //下载进度条



}

//开始下载
func (task *FileTask) Begin() {
	//go func() {
	if task.Size <=1024 {
		//太小了
		task.FileSlice = append(task.FileSlice,FileSlice{0,task.Size})
		log.Printf("小文件，使用单线程下载")
		task.WorkCnt = 1


	}else {
		if task.WorkCnt <=0 {
			//防止出现问题
			task.WorkCnt = 1
		}
		sliceSize :=   task.Size /task.WorkCnt
		//sliceSize :=   task.Size /task.WorkCnt
		if sliceSize==0 {
			//文件太小了
			for sliceSize==0 && task.WorkCnt>1 {
				task.WorkCnt--
				sliceSize = task.Size/task.WorkCnt
			}
		}

		var sliceBegin int64 = 0
		for i:=0;i<int(task.WorkCnt);i++ {
			//遍历
			var end = int64(i + 1) * sliceSize
			task.FileSlice = append(task.FileSlice, FileSlice{sliceBegin,end})
			sliceBegin = end + 1
		}
		//剩余部分给最后的线程
		task.FileSlice[task.WorkCnt-1].End = task.Size


	}
	//}()
	err := task.download()
	if err!=nil {
		log.Fatalf("下载异常 _task.download status := %v\n",err)
	}
}

//下载到临时文件里面
func (task *FileTask) BeginDownloadInTmp() {
	//go func() {
	if task.Size <=1024 {
		//太小了
		task.FileSlice = append(task.FileSlice,FileSlice{0,task.Size})
		log.Printf("小文件，使用单线程下载")
		task.WorkCnt = 1


	}else {
		if task.WorkCnt <=0 {
			//防止出现问题
			task.WorkCnt = 1
		}
		sliceSize :=   task.Size /task.WorkCnt
		if sliceSize==0 {
			//文件太小了
			for sliceSize==0 && task.WorkCnt>1 {
				task.WorkCnt--
				sliceSize = task.Size/task.WorkCnt
			}
		}

		var sliceBegin int64 = 0
		for i:=0;i<int(task.WorkCnt);i++ {
			//遍历
			var end = int64(i + 1) * sliceSize
			task.FileSlice = append(task.FileSlice, FileSlice{sliceBegin,end})
			sliceBegin = end + 1
		}
		//剩余部分给最后的线程
		task.FileSlice[task.WorkCnt-1].End = task.Size


	}
	//}()
	err := task.downloadToTmp()
	if err!=nil {
		log.Fatalf("\n下载异常 _task.download status := %v\n",err)
	}
}



func (t *FileTask) download() error {
	workerChannel := make(chan int,t.WorkCnt)
	if t.Bar==nil {
		t.Bar = BarInstance(int(t.WorkCnt))
	}
	for i:= range t.FileSlice {
		go func(id int) {
			for {
				//尝试下载
				sliceErr := t.DownloadBlock(id)
				if sliceErr!=nil {
					//t.DownloadBlock(id)
					log.Printf("下载重试 :=%v ,异常信息 :=%v",id,sliceErr)

					continue
				}
				log.Printf("任务:%d 下载完成\n",id)
				break
			}
			workerChannel <- id
		}(i)
		log.Printf("准备任务:=%v ，%v",i,&t.FileSlice[i])
	}
	var res int
	for i:=0;i<int(t.WorkCnt);i++ {

		//等待通道结果
		id := <- workerChannel
		t.DownloadSize += t.FileSlice[id].End - t.FileSlice[id].Start + 1
		res++
		log.Printf("任务下载完成 ,完成数 %d\n",res)
		t.Bar.PrintBar(res)
	}
	log.Printf("总下载量为 :=%v", t.DownloadSize)

	return nil

}



func (t *FileTask) downloadToTmp() error {
	workerChannel := make(chan int,t.WorkCnt)
	if t.Bar==nil {
		t.Bar = BarInstance(int(t.WorkCnt))
	}
	for i:= range t.FileSlice {
		go func(id int) {
			for {
				//尝试下载
				sliceErr := t.DownloadToTmpPath(id)
				if sliceErr!=nil {
					//t.DownloadBlock(id)
					log.Printf("下载重试 :=%v ,异常信息 :=%v\n",id,sliceErr)

					continue
				}
				log.Printf("\n任务:%d 下载完成\n",id)
				break
			}
			workerChannel <- id
		}(i)
		log.Printf("准备任务:=%v ，%v\n",i,&t.FileSlice[i])
	}
	var res int
	for i:=0;i<int(t.WorkCnt);i++ {

		//等待通道结果
		id := <- workerChannel
		t.DownloadSize += t.FileSlice[id].End - t.FileSlice[id].Start + 1
		res++
		log.Printf("\n任务下载完成 ,完成数 %d\n",res)
		t.Bar.PrintBar(res)
	}
	//开始合并文件
	log.Printf("\n开始合并文件 ----  -----\n")
	var targetFileOffset int64 = 0
	var buf =make([]byte,t.BufferSize)
	for id,_ := range t.FileSlice {
		originName := t.File.Name()
		tmpName_ := fileName(originName)
		//indexName := IndexName(tmpName_, int(t.WorkCnt),id)
		indexName := TmpFilePath(t.TmpPath,tmpName_, int(t.WorkCnt),id)
		//t.File.WriteAt()
		sourceFile,serr := os.Open(indexName)
		if serr!=nil {
			log.Fatalf("tmp file open error := %v",serr)
		}
		defer sourceFile.Close()

		for  {
			n,err := sourceFile.Read(buf)
			if err == io.EOF {
				break
			}
			if err!=nil {
				log.Fatalf("writing error occured := %d\n",err)
				return err
			}
			if n>0 {
				t.File.WriteAt(buf[:n],targetFileOffset)
				targetFileOffset = targetFileOffset+int64(n)
			}
		}
	}
	log.Printf("总下载量为 :=%v", t.DownloadSize)

	return nil

}



func (t *FileTask) DownloadBlock(id int) error {
	req,err := NewGetRequestObject(t.Url)
	if err !=nil {
		return err
	}
	ref := &t.FileSlice[id]
	l := ref.Start
	r := ref.End
	if r!=-1 && t.WorkCnt>1 {
		//如果 r 是 -1 ，就不用设置 range 分片，否则的话 就是 range分片下载
		req.Header.Set("Range",
			fmt.Sprintf("bytes=%v-%v",l,r))
	}
	resp ,err := NewRequestClient(t.ProxyUrl,t.Timeout).Do(req)
	//获得响应
	if err!=nil {
		return err
	}
	defer resp.Body.Close()
	var buf = make([]byte,t.BufferSize)
	//writeOFFset
	writeOffset := l
	for {
		//循环读取任务
		n,derr := resp.Body.Read(buf)
		if n>0 {
			//注意： n=0 不一定就是下载完成了
			writedSize,err := t.File.WriteAt(buf[0:n], writeOffset)
			if err!=nil {
				log.Fatalf("write to file error %v\n",err)
				os.Exit(1)

			}
			writeOffset += int64(writedSize)

		}
		//log.Printf("n:= %d",n)
		//log.Printf("err := %v",derr)
		if derr==io.EOF {
			 //t.File.
			 return nil
		}
		if derr!=nil {
			return derr
		}

		//if n<t.BufferSize {
		//	//比如 倒一桶水，没装满，说明没水了
		//	break
		//}

	}
	//------------ endl ----------------
	return nil

}







// 下载到临时目录里面，后面再合并

func (t *FileTask) DownloadToTmpPath(id int) error {
	req,err := NewGetRequestObject(t.Url)
	if err !=nil {
		return err
	}
	originName := t.File.Name()
	tmpName_ := fileName(originName)
	//indexName := IndexName(tmpName_, int(t.WorkCnt),id)
	indexName := TmpFilePath(t.TmpPath,tmpName_, int(t.WorkCnt),id)
	ref := &t.FileSlice[id]
	l := ref.Start
	r := ref.End

	var fileContentOffset int64 = 0
	if r!=-1 && t.WorkCnt>1 {
		tmpStat, tmpErr := os.Stat(indexName)

		if tmpErr ==nil && tmpStat!=nil && tmpStat.Size() == r-l+1 {
			log.Printf("\n检索到已下载的分片 id:=%v\n",id)
			//已经下载完成了
			return nil
		}
		if tmpErr==nil  && tmpStat!=nil {
			fileContentOffset = tmpStat.Size()
			log.Printf("任务:%v 使用断点下载",id)
		}
		//如果 r 是 -1 ，就不用设置 range 分片，否则的话 就是 range分片下载
		req.Header.Set("Range",
			fmt.Sprintf("bytes=%v-%v",l+fileContentOffset,r))

	}

	resp ,err := NewRequestClient(t.ProxyUrl,t.Timeout).Do(req)
	var tmpFile *os.File

	tmpFile,err = os.OpenFile(indexName,os.O_WRONLY|os.O_CREATE,0777)
	if err!=nil {
		log.Fatalf("append to file err := %v",err)
		return err
	}
	//tmpFile.Seek()
	defer tmpFile.Close()


	//获得响应
	if err!=nil {
		return err
	}
	defer resp.Body.Close()
	var buf = make([]byte,t.BufferSize)
	//writeOFFset , 文件写入的位置， 从文件末尾写入
	var writeOffset int64 = fileContentOffset

	for {

		//循环读取任务
		n,derr := resp.Body.Read(buf)
		if n>0 {
			//注意： n=0 不一定就是下载完成了
			writedSize,err := tmpFile.WriteAt(buf[0:n], writeOffset)
			if err!=nil {
				log.Fatalf("write to file error %v\n",err)
				os.Exit(1)

			}
			writeOffset += int64(writedSize)


		}

		if derr==io.EOF {
			//t.File.
			return nil
		}
		if derr!=nil {
			return derr
		}

	}
	//------------ endl ----------------
	return nil

}





