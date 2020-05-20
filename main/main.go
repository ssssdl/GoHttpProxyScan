package main

import (
	"../web"
	"log"
)

func main() {
	//todo：总结配置在main函数中
	//todo 写几个网页用于演示
	/***** 【日志信息配置】 *****/
	log.SetPrefix("【main】")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)

	/***** 【创建消息队列】 *****/

	/***** 【启动代理程序】 *****/
	//go http_proxy_simple()
	//go https_Decrypt()	//简易https代理

	/***** 【启动WEB服务】 *****/
	web.Server(":8080")

}
