package plugin

import (
	"net/http"
	"strings"
)

//扫描点击劫持的插件
type ClickHijacking struct{}

//todo 插件文档整理以下规范，以下为规范，需要使用者填写info中的信息
func (p *ClickHijacking) GetPluginInfo() map[string]string {
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
	info["name"] = "Click Hijacking Scan"              //插件名称
	info["appPowerLink"] = "https://github.com/ssssdl" //漏洞组件官网
	info["appName"] = "http"                           //漏洞组件名称，如apache
	info["appVersion"] = "1.1"                         //漏洞组件版本
	info["vulType"] = "Click Hijacking"                //漏洞类型		//todo 文档中整理标准漏洞类型，漏洞类型只能填文档中的
	info["vulDesc"] = ""                               //漏洞描述
	info["samples"] = "http://localhost"               //测试样例，利用该插件测试成功的网站
	info["install_requires"] = ""                      //插件需要提前安装的第三方库，填github地址
	info["desc"] = "用于扫描点击劫持的插件，基于X-Frame-Options"     //插件描述

	return info
}

func (p *ClickHijacking) GetHttp(Req *http.Request, Resp string) {
	/**
	思路：
	如果响应中没有含有X-Frame-Options，断定页面可能存在点击劫持
	***/
	if !strings.ContainsAny(Resp, "X-Frame-Options") {
		Result := Req.URL.String() + " 请求中含有点击劫持漏洞"
		PutScanResults("ERROR", Result, "Click Hijacking Scan")
	}
}
