package plugin

import (
	"net/http"
)

type pluginfunc interface {
	request(Req *http.Request)
	respones(Resp string)
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
	plugin := new(plugins)
	plugin.init()

	plugin1 := new(plugin1)
	plugin2 := new(plugin2)
	plugin3 := new(plugin3)
	plugin.register("plugin1", plugin1)
	plugin.register("plugin2", plugin2)
	plugin.register("plugin3", plugin3)
	for _, plugin := range plugin.plist {
		plugin.request(Req)
		plugin.respones(Resp)
	}
}
