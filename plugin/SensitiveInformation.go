package plugin

import (
	"net/http"
	"regexp"
)

//todo 测试此插件
//用于扫描敏感信息的插件
//主要扫描的敏感信息包括身份证号  电话号和邮箱，还有内网IP

type SensitiveInformation struct{}

//todo 插件文档整理以下规范，以下为规范，需要使用者填写info中的信息
func (p *SensitiveInformation) GetPluginInfo() map[string]string {
	var info map[string]string
	info = make(map[string]string)

	/*****【填写插件信息】*****/
	info["vulID"] = "0003"                             //漏洞编号
	info["version"] = "1"                              //插件版本，从1开始，递增
	info["author"] = "Archer"                          //插件作者
	info["vulDate"] = "2020-05-05"                     //漏洞公开时间
	info["createDate"] = "2020-05-05"                  //插件编写时间
	info["updateDate"] = "2020-05-05"                  //插件更新时间，默认和编写时间相同
	info["references"] = "https://github.com/ssssdl"   //漏洞来源
	info["name"] = "Sensitive information Scanner"     //插件名称
	info["appPowerLink"] = "https://github.com/ssssdl" //漏洞组件官网
	info["appName"] = "http"                           //漏洞组件名称，如apache
	info["appVersion"] = "1.1"                         //漏洞组件版本
	info["vulType"] = "info"                           //漏洞类型		//todo 文档中整理标准漏洞类型，漏洞类型只能填文档中的
	info["vulDesc"] = ""                               //漏洞描述
	info["samples"] = "http://localhost"               //测试样例，利用该插件测试成功的网站
	info["install_requires"] = ""                      //插件需要提前安装的第三方库，填github地址
	info["desc"] = "扫描网页中的敏感信息，包括身份证号、电话、邮箱、内网IP"      //插件描述

	return info
}

func (p *SensitiveInformation) GetHttp(Req *http.Request, Resp string) {
	patterns := [...]string{
		"[\\d]+\\.[\\d]+\\.[\\d]+\\.[\\d]+", //匹配IP地址
		"\\d{11}",                           //匹配电话号码
		"\\d{18}",                           //匹配身份证号
		`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`, //匹配邮箱
	}
	var finder *regexp.Regexp
	var info string
	for i := range patterns {
		finder = regexp.MustCompile(patterns[i])
		info = finder.FindString(Resp)
		if info != "" {
			Result := Req.URL.String() + " 请求中含有敏感信息" + info
			PutScanResults("INFO", Result, "Sensitive information Scanner")
		}
	}
}
