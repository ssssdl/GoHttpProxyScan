package plugin

import (
	"../MassageQueue"
	"net/http"
)

type HttpHistory struct{}

//todo 插件文档整理以下规范，以下为规范，需要使用者填写info中的信息
func (p *HttpHistory) GetPluginInfo() map[string]string {
	var info map[string]string
	info = make(map[string]string)

	/*****【填写插件信息】*****/
	info["vulID"] = "0001"                             //漏洞编号
	info["version"] = "1"                              //插件版本，从1开始，递增
	info["author"] = "Archer"                          //插件作者
	info["vulDate"] = "2020-05-05"                     //漏洞公开时间
	info["createDate"] = "2020-05-05"                  //插件编写时间
	info["updateDate"] = "2020-05-05"                  //插件更新时间，默认和编写时间相同
	info["references"] = "https://github.com/ssssdl"   //漏洞来源
	info["name"] = "Http History"                      //插件名称
	info["appPowerLink"] = "https://github.com/ssssdl" //漏洞组件官网
	info["appName"] = "http"                           //漏洞组件名称，如apache
	info["appVersion"] = "1.1"                         //漏洞组件版本
	info["vulType"] = "info"                           //漏洞类型		//todo 文档中整理标准漏洞类型，漏洞类型只能填文档中的
	info["vulDesc"] = "基础插件无漏洞"                        //漏洞描述
	info["samples"] = "http://localhost"               //测试样例，利用该插件测试成功的网站
	info["install_requires"] = ""                      //插件需要提前安装的第三方库，填github地址
	info["desc"] = "打印http历史记录"                        //插件描述

	return info
}

func (p *HttpHistory) GetHttp(Req *http.Request, Resp string) {
	var msg map[string]string
	msg = make(map[string]string)
	msg["URL"] = Req.URL.String() //获取url
	msg["Status"] = Resp[9:12]    //获取响应状态码
	MassageQueue.HttpHistoryQueue.Put(msg)
}
