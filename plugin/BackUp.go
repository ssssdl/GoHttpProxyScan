package plugin

import (
	"bufio"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

//数据包备份插件
type BackUp struct{}

//todo 插件文档整理以下规范，以下为规范，需要使用者填写info中的信息
func (p *BackUp) GetPluginInfo() map[string]string {
	var info map[string]string
	info = make(map[string]string)

	/*****【填写插件信息】*****/
	info["vulID"] = "0002"                             //漏洞编号
	info["version"] = "1"                              //插件版本，从1开始，递增
	info["author"] = "Archer"                          //插件作者
	info["vulDate"] = "2020-05-05"                     //漏洞公开时间
	info["createDate"] = "2020-05-20"                  //插件编写时间
	info["updateDate"] = "2020-05-20"                  //插件更新时间，默认和编写时间相同
	info["references"] = "https://github.com/ssssdl"   //漏洞来源
	info["name"] = "Http Back Up"                      //插件名称
	info["appPowerLink"] = "https://github.com/ssssdl" //漏洞组件官网
	info["appName"] = "http"                           //漏洞组件名称，如apache
	info["appVersion"] = "1.1"                         //漏洞组件版本
	info["vulType"] = "Info"                           //漏洞类型		//todo 文档中整理标准漏洞类型，漏洞类型只能填文档中的
	info["vulDesc"] = ""                               //漏洞描述
	info["samples"] = "http://localhost"               //测试样例，利用该插件测试成功的网站
	info["install_requires"] = ""                      //插件需要提前安装的第三方库，填github地址
	info["desc"] = "http流量导出备份"                        //插件描述

	return info
}

func (p *BackUp) GetHttp(Req *http.Request, Resp string) {
	//思路
	//在根目录下创建备份文件夹
	//文件为按小时备份，文件名为时间精确到小时，后缀.log
	//数据包的分割符==================【时间戳】=========================
	//好 开写
	//测试写入文件
	//t := time.Now()
	//log.Println(t.Format("2020052019"))
	logFileName := time.Now().Format("2006010215") + ".log"
	logDir := "./HttpBackUp/"
	dumpReq, _ := httputil.DumpRequest(Req, true)
	context := "\n==================[" + time.Now().Format("2006-01-02 15:04:05") + "]=======================\n" +
		string(dumpReq) + "\n" +
		"=======================================================\n" +
		Resp
	file, err := os.OpenFile(logDir+logFileName, os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("日志文件的打开异常")
		log.Fatal(err)
	}
	defer file.Close()
	buffer := bufio.NewWriter(file)

	bw, err := buffer.WriteString(context)
	if err != nil {
		log.Fatal(err)
	}
	buffer.Flush()
	log.Printf("Bytes written: %d\n", bw)
}
