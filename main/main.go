package main

import (
	"../web"
	"log"
)

func main() {
	//todo：总结配置在main函数中
	/***** 【日志信息配置】 *****/
	log.SetPrefix("【main】")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	/***** 【启动WEB服务】 *****/

	//传入启动地址
	web.Server(":8080")

}
