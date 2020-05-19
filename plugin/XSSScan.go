package plugin

import (
	"../MassageQueue"
	"log"
	"net/http"
	"net/url"
	"strings"
)

//todo 制作获取插件信息接口

type XSSScan struct{}

//todo 插件文档整理以下规范，以下为规范，需要使用者填写info中的信息
func (p *XSSScan) GetPluginInfo() map[string]string {
	var info map[string]string
	info = make(map[string]string)

	/*****【填写插件信息】*****/
	info["vulID"] = "0002"                             //漏洞编号
	info["version"] = "1"                              //插件版本，从1开始，递增
	info["author"] = "Archer"                          //插件作者
	info["vulDate"] = "2020-05-05"                     //漏洞公开时间
	info["createDate"] = "2020-05-05"                  //插件编写时间
	info["updateDate"] = "2020-05-05"                  //插件更新时间，默认和编写时间相同
	info["references"] = "https://github.com/ssssdl"   //漏洞来源
	info["name"] = "XSS Scanner"                       //插件名称
	info["appPowerLink"] = "https://github.com/ssssdl" //漏洞组件官网
	info["appName"] = "http"                           //漏洞组件名称，如apache
	info["appVersion"] = "1.1"                         //漏洞组件版本
	info["vulType"] = "info"                           //漏洞类型		//todo 文档中整理标准漏洞类型，漏洞类型只能填文档中的
	info["vulDesc"] = ""                               //漏洞描述
	info["samples"] = "http://localhost"               //测试样例，利用该插件测试成功的网站
	info["install_requires"] = ""                      //插件需要提前安装的第三方库，填github地址
	info["desc"] = "简单扫描反射型XSS"                        //插件描述

	return info
}

func (p *XSSScan) GetHttp(Req *http.Request, Resp string) {
	var msg map[string]string
	msg = make(map[string]string)
	/**
	思路：
	获取url中参数内容,到响应中匹配，如果匹配到了发送警告消息，标明是那个参数，如果含有<>"'等字符发送严重消息
	***/
	queryForm, err := url.ParseQuery(Req.URL.RawQuery)
	if err == nil {
		for param := range queryForm {
			if len(queryForm[param]) > 0 {
				log.Println(queryForm[param][0])
				//匹配响应中是否含有请求内容
				if strings.Contains(Resp, queryForm[param][0]) {
					if strings.ContainsAny(queryForm[param][0], "<") {
						/*******【整理消息】********/
						msg["Level"] = "ERROR"
						msg["Content"] = "XSS Scan:" + Req.URL.String() + " 请求中含有XSS漏洞，出现漏洞的参数为" + param

						/*****【将消息发送到推送队列】******/
						MassageQueue.MsgQueue.Put(msg)
					} else {
						msg["Level"] = "WARNING"
						msg["Content"] = "XSS Scan:" + Req.URL.String() + " 请求中可能含有XSS漏洞，出现漏洞的参数为" + param

						/*****【将消息发送到推送队列】******/
						MassageQueue.MsgQueue.Put(msg)
					}
				}
			}
		}
	}
}
