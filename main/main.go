package main

import (
	"../proxy"
	"../web"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// 简易http代理
//func http_proxy_simple(MsgQueue *MassageQueue.MassageQueue){
func http_proxy_simple() {
	proxy := proxy.New()
	server := &http.Server{
		Addr:         ":8081",
		Handler:      proxy,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// https解密 需要导入根证书 mitm-proxy.crt
// 实现证书缓存接口
type Cache struct {
	m sync.Map
}

func (c *Cache) Set(host string, cert *tls.Certificate) {
	c.m.Store(host, cert)
}
func (c *Cache) Get(host string) *tls.Certificate {
	v, ok := c.m.Load(host)
	if !ok {
		return nil
	}

	return v.(*tls.Certificate)
}
func https_Decrypt() {
	proxy := proxy.New(proxy.WithDecryptHTTPS(&Cache{}))
	server := &http.Server{
		Addr:         ":8081",
		Handler:      proxy,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func main() {
	//todo：多线程启动sse和proxy

	/***** 【日志信息配置】 *****/
	log.SetPrefix("【main】")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)

	/***** 【创建消息队列】 *****/

	/***** 【启动代理程序】 *****/
	//go http_proxy_simple(MsgQueue)
	//https_Decrypt(MsgQueue)	//简易https代理

	/***** 【Put测试】 *****/
	go func() {
		time.Sleep(5000 * time.Millisecond) //不只是给浏览器加载网页的时间  还是一个异步的时间  没有这个就加载不到消息 和bug无关
		//for i:=0;i<100 ;i++  {
		//	MassageQueue.MsgQueue.Put(strconv.Itoa(i))
		//}
		//time.Sleep(2000*time.Millisecond)
		//for i:=0;i<100 ;i++  {
		//	MassageQueue.MsgQueue.Put(strconv.Itoa(i))
		//}
		http_proxy_simple()
	}()
	/***** 【启动WEB服务】 *****/
	//设置sse刷新频率
	//time.Sleep(10 * time.Second)
	REFRESH := 1 //刷新间隔,防止占用服务器资源
	go web.Server(":8080", REFRESH)

	//防止线程结束
	var str string
	fmt.Scan(&str)
}
