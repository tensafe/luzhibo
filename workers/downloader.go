package workers

import (
	"io"
	"net/http"
	"os"
	"path"
	nhttp"github.com/Baozisoftware/golibraries/http"
	"github.com/Baozisoftware/luzhibo/api/getters"
)

//下载器

type downloader struct {
	url      string
	filePath string
	cb       WorkCompletedCallBack
	run      bool
	ch       chan bool
	ch2      chan bool
	ch3      chan bool
	fm       bool
	client   *nhttp.HttpClient
}

func newDownloader(url, filepath string, callbcak WorkCompletedCallBack) *downloader {
	if url != "" && filepath != "" {
		r := &downloader{}
		r.url = url
		r.filePath = filepath
		r.cb = callbcak
		r.client = nhttp.NewHttpClient()
		return r
	}
	return nil
}

//Start 实现接口
func (i *downloader) Start() {
	if i.run {
		return
	}
	i.run = true
	i.ch = make(chan bool, 0)
	i.ch2 = make(chan bool, 1)
	i.ch3 = make(chan bool, 1)
	if i.fm {
		go i.ffmpeg(i.url, i.filePath)
	} else {
		go i.http(i.url, i.filePath)
	}
}

//Stop 实现接口
func (i *downloader) Stop() {
	if i.run {
		i.ch2 <- true
		i.run = false
		i.ch3 <- true
		<-i.ch
		close(i.ch)
		close(i.ch2)
		close(i.ch3)
	}
}

//Restart 实现接口
func (i *downloader) Restart() (Worker, error) {
	if i.run {
		i.Stop()
	}
	i.Start()
	return i, nil
}

//GetTaskInfo 实现接口
func (i *downloader) GetTaskInfo(g bool) (int64, bool, int64, string, *getters.LiveInfo) {
	return 0, i.run, 0, i.filePath, nil
}

func (i *downloader) http(url, filepath string) {
	ec := int64(0) //正常停止
	defer func() {
		if !i.run {
			i.ch <- true
		}
		if !i.run {
			ec = 1 //主动停止
		}
		i.run = false
		if i.cb != nil {
			i.cb(ec)
		}
	}()
	client, resp, err := httpGetResp(url)
	client.SetReadBodyTimeout(300)
	if err != nil || resp.StatusCode != 200 {
		ec = 2 //请求时错误
		return
	}
	defer resp.Body.Close()
	f, err := createFile(filepath)
	if err != nil {
		ec = 3 //创建文件错误
		return
	}
	defer f.Close()
	go func() {
		for i.run {
			data, err := client.ReadBodyWithTimeout(resp)
			if err != nil {
				if err == io.EOF {
					_, err = f.Write(data)
					if err != nil {
						ec = 5 //写入文件错误
					} else {
						break
					}
				} else {
					ec = 4 //下载数据错误
				}
			} else {
				_, err = f.Write(data)
				if err != nil {
					ec = 5 //写入文件错误
				}
			}
			if ec > 0 {
				break
			}
		}
		if i.run {
			i.ch3 <- true
		}
	}()
	<-i.ch3
}

func httpGetResp(url string) (client *nhttp.HttpClient, resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err == nil {
		client = nhttp.NewHttpClient()
		client.SetProxy(Proxy)
		client.SetResponseHeaderTimeout(30)
		resp, err = client.Do(req)
	}
	return
}

func (i *downloader) ffmpeg(url, filepath string) {
	ec := int64(0) //正常停止
	defer func() {
		if !i.run {
			i.ch <- true
		}
		if !i.run {
			ec = 1 //主动停止
		}
		i.run = false
		if i.cb != nil {
			i.cb(ec)
		}
	}()

	err := os.MkdirAll(path.Dir(filepath), os.ModePerm)
	if err != nil {
		ec = 3
		return
	}
	cmd := NewFFmpeg(i.url, i.filePath)
	go func() {
		if err := cmd.Start(); err != nil {
			ec = 6 //ffmpeg启动失败
			i.ch2 <- true
		}
		cmd.Wait()
		if i.run {
			i.ch2 <- true
		}
	}()
	<-i.ch2
	if cmd.Process != nil {
		cmd.Process.Kill()
	}
}

func (i *downloader) UseFFmpeg() {
	i.fm = true
}
