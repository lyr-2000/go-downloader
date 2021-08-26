## 介绍


这个是基于 golang 实现的多线程下载工具

使用方法
```go
go build main.go -o download.exe


```

编译出二进制文件之后

`main.exe -h`
查询使用方法

## 支持的功能
1. 断点下载（多线程下载）
2. 设置超时时间
3. 错误重试


### 断点下载实现原理

断点下载是先将分片下载到一个文件里面
用户退出进程后重新进入进程，会检查分片文件是否下载完成
如果没下载完成，就在分片文件末尾追加内容 ，通过  http 请求 的 Ranges 字段 来指定 文件范围

### usage

```
Usage of download.exe:
  -bd
        断点下载,breakpoint downloading (default true)
  -buf int
        缓冲区大小【单位：字节】 (default 4096)
  -c int
        下载线程数 (default 8)
  -path string
        文件下载路径 (default "./")
  -proxy string
        设置代理,【本人没用过】：例如:http://127.0.0.1:7890
  -retry
        失败重试, (default true)
  -s int
        设置超时时间，单位为秒 (default 25)
  -tmp string
        下载的临时文件 (default "./tmp")
  -url string
        url不能为空，请设置 -url 指定

```
