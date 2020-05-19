# 插件开发文档及编写规范

## 关于插件实现
> 插件功能实现思路是利用通过实现接口方法，在loader中统一实例化并执行实现

## 插件接口
插件中必须要实现的方法
```
func (p *plugin) GetHttp(Req *http.Request, Resp string)        //加载插件时默认执行的函数
func (p *SensitiveInformation) GetPluginInfo() map[string]string//返回插件的信息
```
插件向前端显示结果，发送消息
```
//可以直接调用消息队列
var msg map[string]string
msg = make(map[string]string)
msg["Level"] = "消息级别"
msg["Content"] = "消息内容"
MassageQueue.MsgQueue.Put(msg)

//也可以调用加载器中提供的方法（参数分别是消息级别，消息内容，发出消息的插件）
func PutScanResults(Level string,Content string,plugin string)
```

另外消息队列有两个除了发送扫描结果的消息队列还有展示HTTP历史记录的消息队列
```
MassageQueue.HttpHistoryQueue       //http历史记录队列
MassageQueue.MsgQueue               //扫描消息队列
```
## 使用方法
> 编写好插件放入plugin目录，并且配置loader.go中的Loader函数加载插件

```
//比如添加以下三个插件

/*****【实例化插件】******/
//所有的插件写好要到这里来实例化
XSSScan := new(XSSScan)
HttpHistory := new(HttpHistory)
SensitiveInformation := new(SensitiveInformation)
/*****【注册插件】******/
//所有的插件写好要到这里来注册
plugin.register("XSSScan", XSSScan)
plugin.register("HttpHistory", HttpHistory)
plugin.register("SensitiveInformation", SensitiveInformation)
```

## 关于GetPluginInfo()
只需要根据注释添加相关信息即可
```
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
	info["vulType"] = "XSS"                           //漏洞类型		//todo 文档中整理标准漏洞类型，漏洞类型只能填文档中的
	info["vulDesc"] = ""                               //漏洞描述
	info["samples"] = "http://localhost"               //测试样例，利用该插件测试成功的网站
	info["install_requires"] = ""                      //插件需要提前安装的第三方库，填github地址
	info["desc"] = "简单扫描反射型XSS"                        //插件描述

	return info
}
```