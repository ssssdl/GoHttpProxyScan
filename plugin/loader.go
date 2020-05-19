package plugin

import (
	"../MassageQueue"
	"log"
	"net/http"
)

//todo 再写一个敏感信息扫描的插件
//todo 还可以写一个扫描点击劫持的插件
//todo 还可以写一个扫描httponly的插件
//todo  F5泄露内网IP的插件
//todo  整理插件文档
//MassageQueue.HttpHistoryQueue		http历史记录消息队列
//MassageQueue.MsgQueue				插件扫描消息队列

type pluginfunc interface {
	//request(Req *http.Request)
	//respones(Resp string)
	GetHttp(Req *http.Request, Resp string) //获取http流量传入到插件中   加载插件自动执行该函数
	GetPluginInfo() map[string]string       //todo 插件信息				后续整理插件信息到打印到网页上
}

// 定义一个类，来存放我们的插件
type plugins struct {
	plist map[string]pluginfunc
}

// 初始化插件
func (p *plugins) init() {
	p.plist = make(map[string]pluginfunc)
}

// 注册插件
func (p *plugins) register(name string, plugin pluginfunc) {
	p.plist[name] = plugin
	//p.plist = append(p.plist, a)

}

// Req 请求 resp 响应  模块中要导出的函数，必须首字母大写。
//func Loader(Req *http.Request,Resp string,MsgQueue *MassageQueue.MassageQueue){
func Loader(Req *http.Request, Resp string) {
	/*****【初始化】******/
	plugin := new(plugins)
	plugin.init()

	/*****【实例化插件】******/
	//所有的插件写好要到这里来实例化
	XSSScan := new(XSSScan)
	HttpHistory := new(HttpHistory)
	SensitiveInformation := new(SensitiveInformation)
	//plugin2 := new(plugin2)
	//plugin3 := new(plugin3)

	/*****【注册插件】******/
	//所有的插件写好要到这里来注册
	plugin.register("XSSScan", XSSScan)
	plugin.register("HttpHistory", HttpHistory)
	plugin.register("SensitiveInformation", SensitiveInformation)
	//plugin.register("plugin2", plugin2)
	//plugin.register("plugin3", plugin3)

	/*****【加载插件入口函数】******/
	//后续可以在这里添加其他的函数，会被默认调用
	for _, plugin := range plugin.plist {
		plugin.GetHttp(Req, Resp)
	}
}

//处理扫描结果的接口 参数分别是消息级别，消息内容，插件名称
func PutScanResults(Level string, Content string, plugin string) {
	var msg map[string]string
	msg = make(map[string]string)

	msg["Level"] = Level
	msg["Content"] = "[+] " + plugin + ":" + Content

	/*****【将消息发送到推送队列】******/
	MassageQueue.MsgQueue.Put(msg)

	//打印消息
	log.Println("[" + msg["Level"] + "] " + plugin + ":" + msg["Content"])
}
