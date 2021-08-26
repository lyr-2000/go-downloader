package downloader

import (
	"log"
	"net/http"
	"net/url"
	"time"
)

func NewGetRequest(path string) (*http.Response,error) {
	return http.Get(path)
}

func NewGetRequestObject(path string) (*http.Request,error) {
	req,err := http.NewRequest("GET",path,nil)

	return req,err
}
func defaultClient(timeout int64) *http.Client {
	return &http.Client{
		Timeout: time.Duration(timeout),
	}
}


func NewRequestClient(proxyUri string,timeout int64) *http.Client {
	//如果没有设置代理，就用默认的
	if proxyUri=="" {
		return defaultClient(timeout)
	}

	uri, err := url.Parse(proxyUri)
	if err!=nil {
		log.Printf("无法使用代理 ~~\n")
		//使用默认方式
		return defaultClient(timeout)
	}
	//使用代理下载
	client := http.Client{
		Transport: &http.Transport{
			// 设置代理
			Proxy: http.ProxyURL(uri),

		},
		Timeout: time.Duration(timeout),
	}
	return &client
}