package downloader

import (
	"fmt"
	"goDownloader/util"
	"testing"
)

func TestDownload(t *testing.T) {
	d := NewDefault()

	for  {
		err := d.Download(6,"https://github.com/Wox-launcher/Wox/releases/download/v1.4.1115/Wox-Full-Installer.1.4.1115.exe",
			"./wox_008.exe",
			1024<<2,
			"",
			333,
		)
		if err==nil {
			break
		}
		if err !=nil {
			fmt.Printf(err.Error()+"\n")
		}
	}
	//if err !=nil {
	//	fmt.Printf(err.Error())
	//}
}

func TestDownloadZip(t *testing.T) {
	d := NewDefault()
	err := d.Download(6,"https://github.com/kkdai/youtube/archive/refs/tags/v2.7.4.zip",
		"./result.zip",
		1024<<1,
		"",
		333,
	)
	if err !=nil {
		fmt.Printf(err.Error())
	}
}


func TestDownloadZip22(t *testing.T) {
	d := NewDefault()
	err := d.Download(6,"https://github.com/Wox-launcher/Wox/releases/download/v1.4.1196/Wox-Full-Installer.1.4.1196.exe",
		"wox1.exe",
		1024<<2,
		"",
		33,
	)
	if err !=nil {
		fmt.Printf(err.Error())
	}
}



func TestDownload222(t *testing.T) {
	d := NewDefault()
	err := d.Download(
		4,
		"http://doc.lyr-2000.xyz/post/01.%E7%A8%8B%E5%BA%8F%E8%AF%AD%E8%A8%80/05.android%E7%9B%B8%E5%85%B3/01.mvvm%E6%9E%B6%E6%9E%84/",
		"./res.html",
		2048*8,
		"",
		33,
	)
	if err !=nil {
		fmt.Printf(err.Error())
	}
}

func TestMD5(t *testing.T) {
	fmt.Println(util.MD5("22222222"))

}

