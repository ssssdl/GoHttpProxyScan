package plugin

import (
	"net/http"
	"strings"
)

//F5负载均衡泄露内网IP地址
type BIGipServerpool struct{}

//todo 插件文档整理以下规范，以下为规范，需要使用者填写info中的信息
func (p *BIGipServerpool) GetPluginInfo() map[string]string {
	var info map[string]string
	info = make(map[string]string)

	/*****【填写插件信息】*****/
	info["vulID"] = "0004"                             //漏洞编号
	info["version"] = "1"                              //插件版本，从1开始，递增
	info["author"] = "Archer"                          //插件作者
	info["vulDate"] = "2020-05-05"                     //漏洞公开时间
	info["createDate"] = "2020-05-20"                  //插件编写时间
	info["updateDate"] = "2020-05-20"                  //插件更新时间，默认和编写时间相同
	info["references"] = "https://github.com/ssssdl"   //漏洞来源
	info["name"] = "BIGipServerpool Scan"              //插件名称
	info["appPowerLink"] = "https://github.com/ssssdl" //漏洞组件官网
	info["appName"] = "http"                           //漏洞组件名称，如apache
	info["appVersion"] = "1.1"                         //漏洞组件版本
	info["vulType"] = "Info"                           //漏洞类型		//todo 文档中整理标准漏洞类型，漏洞类型只能填文档中的
	info["vulDesc"] = ""                               //漏洞描述
	info["samples"] = "http://localhost"               //测试样例，利用该插件测试成功的网站
	info["install_requires"] = ""                      //插件需要提前安装的第三方库，填github地址
	info["desc"] = "用于扫描点击劫持的插件，基于BIGipServerpool"     //插件描述

	return info
}

func (p *BIGipServerpool) GetHttp(Req *http.Request, Resp string) {
	/**
	思路：
	扫描相关参数
	***/
	if !strings.ContainsAny(Resp, "BIGipServerpool") {
		Result := Req.URL.String() + " 请求中含有BIGipServerpool，可能存在泄露内网IP泄露"
		PutScanResults("WARNING", Result, "BIGipServerpool Scan")
	}
}
