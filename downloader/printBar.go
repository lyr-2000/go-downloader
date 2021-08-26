package downloader

import (
	"fmt"
	"github.com/redmask-hb/GoSimplePrint/goPrint"
)

func BarInstance(workerCnt int ) *goPrint.Bar{
	bar := goPrint.NewBar(workerCnt)
	bar.SetNotice("下载进度：")
	bar.SetGraph(">>")
	bar.SetNoticeColor(goPrint.FontColor.Purple)
	bar.SetGraphColor(goPrint.FontColor.Purple)

	//bar.SetColor(goPrint.BarColor{
	//	Notice: goPrint.FontColor.Purple,
	//	Graph: goPrint.FontColor.Aqua,
	//})

	return bar
}

func getWorkingBar(byteSize int64,bufSize int64,id int) *goPrint.Bar {
	cnt:=byteSize/bufSize

	bar := goPrint.NewBar(int(cnt))
	bar.SetNotice(fmt.Sprintf("任务[%v]下载进度",id))
	bar.SetGraph("☆")

	return bar
}

