package downloader

import "github.com/redmask-hb/GoSimplePrint/goPrint"

func BarInstance(workerCnt int ) *goPrint.Bar{
	bar := goPrint.NewBar(workerCnt)
	bar.SetNotice("下载进度：")
	bar.SetGraph(">>")
	bar.SetNoticeColor(goPrint.FontColor.Purple)
	return bar
}



