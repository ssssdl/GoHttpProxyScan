package plugin

import (
	"net/http"
)

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
	//plugin2 := new(plugin2)
	//plugin3 := new(plugin3)

	/*****【注册插件】******/
	//所有的插件写好要到这里来注册
	plugin.register("XSSScan", XSSScan)
	plugin.register("HttpHistory", HttpHistory)
	//plugin.register("plugin2", plugin2)
	//plugin.register("plugin3", plugin3)

	/*****【加载插件入口函数】******/
	//后续可以在这里添加其他的函数，会被默认调用
	for _, plugin := range plugin.plist {
		plugin.GetHttp(Req, Resp)
	}
}
