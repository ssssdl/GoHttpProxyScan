package web

/**
关于这个模块设计的想法
首先要有静态资源的一个目录static
数据传输使用sse
然后还有html文件的目录
控制器模块，前端js加载数据 后端返回json
存放数据的接口【这边是个核心难度比较高的东西，如果使用文件的话要控制写入和读取不冲突，使用队列的话要建立队列的传入方式，感觉还是数据库比较好实现但是涉及到画表】
		数据库的话将反序列化对象存到数据库中，响应和请求合成一个对象，插件的消息队列和响应请求的队列一对多
*/
import (
	"../MassageQueue"
	"../proxy"
	"container/list"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	//注意: 由于某种原因，最新的vscode-go语言扩展不能提供足够智能帮助（参数文档并转到定义功能）
	//对于`iris.Context`别名，因此如果您使用VS Code，则导入`Context`的原始导入路径，它将执行此操作：
	"github.com/kataras/iris/context"
)

//sse消息配置
//Broker拥有开放的客户端连接
//在其Notifier频道上侦听传入事件
//并将事件数据广播到所有已注册的连接
type Broker struct {
	//主要事件收集例程将事件推送到此频道
	Notifier chan []byte
	//新的客户端连接
	newClients chan chan []byte
	//关闭客户端连接
	closingClients chan chan []byte
	//客户端连接注册表
	clients map[chan []byte]bool
}

// NewBroker返回一个新的代理工厂
func NewBroker() *Broker {
	b := &Broker{
		Notifier:       make(chan []byte, 1),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}
	//设置它正在运行 - 收听和广播事件
	go b.listen()
	return b
}

//听取不同的频道并采取相应应对
func (b *Broker) listen() {
	for {
		select {
		case s := <-b.newClients:
			//新客户端已连接
			//注册他们的消息频道
			b.clients[s] = true
			golog.Infof("Client added. %d registered clients", len(b.clients))
		case s := <-b.closingClients:
			//客户端已离线，我们希望停止向其发送消息。
			delete(b.clients, s)
			golog.Warnf("Removed client. %d registered clients", len(b.clients))
		case event := <-b.Notifier:
			//我们从外面得到了一个新事件
			//向所有连接的客户端发送事件
			for clientMessageChan := range b.clients {
				clientMessageChan <- event
			}
		}
	}
}

func (b *Broker) ServeHTTP(ctx context.Context) {
	//确保编写器支持刷新
	flusher, ok := ctx.ResponseWriter().Flusher()
	if !ok {
		ctx.StatusCode(iris.StatusHTTPVersionNotSupported)
		ctx.WriteString("Streaming unsupported!")
		return
	}
	//设置与事件流相关的header，如果发送纯文本，则可以省略“application/json”
	//如果你开发了一个go客户端，你必须设置：“Accept”：“application/json，text/event-stream”header
	ctx.ContentType("application/json, text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	//我们还添加了跨源资源共享标头，以便不同域上的浏览器仍然可以连接
	ctx.Header("Access-Control-Allow-Origin", "*")
	//每个连接都使用Broker的连接注册表注册自己的消息通道
	messageChan := make(chan []byte)
	//通知我们有新连接的Broker
	b.newClients <- messageChan
	//监听连接关闭以及整个请求处理程序链退出时（此处理程序）并取消注册messageChan。
	ctx.OnClose(func() {
		//从已连接客户端的map中删除此客户端,当这个处理程序退出时
		b.closingClients <- messageChan
	})
	//阻止等待在此连接的消息上广播的消息
	for {
		//写入ResponseWriter
		// Server Sent Events兼容
		ctx.Writef("data: %s\n\n", <-messageChan)
		//或json：data：{obj}
		//立即刷新数据而不是稍后缓冲它
		flusher.Flush()
	}
}

//todo 以下两处js的链接域名统一配置  配置文件
//http历史记录的js
const script = `<script type="text/javascript">
if(typeof(EventSource) !== "undefined") {
    console.log("server-sent events supported");
    var client = new EventSource("http://localhost:8080/events");
    var index = 1;
    client.onmessage = function (evt) {
        console.log(evt);
        // it's not required that you send and receive JSON, you can just output the "evt.data" as well.
        dataJSON = JSON.parse(evt.data)
        var table = document.getElementById("messagesTable");
        var row = table.insertRow(index);
        var cellTimestamp = row.insertCell(0);
        var cellURL = row.insertCell(1);
		var cellStatus = row.insertCell(2);
        cellTimestamp.innerHTML = dataJSON.Time;
        cellURL.innerHTML = dataJSON.URL;
		cellStatus.innerHTML = dataJSON.Status;
        index++;
        window.scrollTo(0,document.body.scrollHeight);
    };
} else {
    document.getElementById("header").innerHTML = "<h2>SSE not supported by this client-protocol</h2>";
}
</script>`

//扫描结果界面的js
const scanScript = `<script type="text/javascript">
if(typeof(EventSource) !== "undefined") {
    console.log("server-sent events supported");
    var client = new EventSource("http://localhost:8080/eventsscan");
    var index = 1;
    client.onmessage = function (evt) {
        console.log(evt);
        // it's not required that you send and receive JSON, you can just output the "evt.data" as well.
        dataJSON = JSON.parse(evt.data)
        var table = document.getElementById("messagesTable");
        var row = table.insertRow(index);
        var cellTime = row.insertCell(0);
        var cellLevel = row.insertCell(1);
		var cellContent = row.insertCell(2);
        cellTime.innerHTML = dataJSON.Time;
        cellLevel.innerHTML = dataJSON.Level;
		cellContent.innerHTML = dataJSON.Content;
        index++;
        window.scrollTo(0,document.body.scrollHeight);
    };
} else {
    document.getElementById("header").innerHTML = "<h2>SSE not supported by this client-protocol</h2>";
}
</script>`

//http历史记录消息
type event struct {
	Time   string `json:"Time"`
	URL    string `json:"URL"`
	Status string `json:"Status"`
}

//扫描消息
type eventScan struct {
	Time    string `json:"Time`
	Level   string `json:"Level"`
	Content string `json:"Content"`
}

//关于代理相关的函数
func http_proxy_simple(addr string) {
	proxy := proxy.New()
	server := &http.Server{
		Addr:         addr,
		Handler:      proxy,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// https解密 需要导入根证书 mitm-proxy.crt
// 实现证书缓存接口
type Cache struct {
	m sync.Map
}

func (c *Cache) Set(host string, cert *tls.Certificate) {
	c.m.Store(host, cert)
}
func (c *Cache) Get(host string) *tls.Certificate {
	v, ok := c.m.Load(host)
	if !ok {
		return nil
	}

	return v.(*tls.Certificate)
}
func https_Decrypt(addr string) {
	proxy := proxy.New(proxy.WithDecryptHTTPS(&Cache{}))
	server := &http.Server{
		Addr:         addr,
		Handler:      proxy,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

var Addr = "localhost:8080" //默认访问地址
var ProxyPort = 8081

func Server(addr string) {
	/***** 【配置web服务】 *****/
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())
	app.Logger().SetLevel("debug") //设置日志级别为调试级别

	//todo 消息展示除了使用网页的方式 还要添加命令行的方式，main函数中配置命令行是否显示http历史记录的开关
	/***** 【配置sse服务】 *****/
	//第一个消息sse ： 用于推送http历史记录
	broker := NewBroker()
	//发送消息
	go func() {
		//如果消息队列不为空则向浏览器推送消息
		msg := make(map[string]string)
		for {
			//todo config：可以在这里调整刷新频率，配置
			// time.Sleep(200 * time.Millisecond)
			if MassageQueue.HttpHistoryQueue.Size() != 0 {
				msg = MassageQueue.HttpHistoryQueue.Get() //根除这个问题不能通过判断massage是否为空  找到size突然不变的原因
				now := time.Now()
				evt := event{
					Time: now.Format(time.RFC1123),
					//Timestamp: MassageQueue.MsgQueue.Size(),
					URL:    fmt.Sprintf("%s", msg["URL"]),
					Status: fmt.Sprintf("%s", msg["Status"]),
				}
				evtBytes, err := json.Marshal(evt)
				if err != nil {
					golog.Error(err)
					continue
				}
				broker.Notifier <- evtBytes
			}
		}
	}()
	//http历史记录推送
	app.Get("/ssehttp", func(ctx context.Context) {
		ctx.HTML(
			`<html><head><title>HTTP历史记录</title>` + script + `</head>
                <body>
                    <h1 id="header">HTTP历史记录</h1>
                    <table id="messagesTable" border="1">
                        <tr>
                            <th>时间</th> 
                            <th>url</th>
							<th>状态</th>
                        </tr>
                    </table>
                </body>
             </html>`)
	})
	app.Get("/events", broker.ServeHTTP)

	//第二个消息sse，用于推送扫描结果
	brokerScanMsg := NewBroker()
	go func() {
		//如果消息队列不为空则向浏览器推送消息
		msg := make(map[string]string)
		for {
			//todo config：可以在这里调整刷新频率，配置
			// time.Sleep(200 * time.Millisecond)
			if MassageQueue.MsgQueue.Size() != 0 {
				msg = MassageQueue.MsgQueue.Get() //根除这个问题不能通过判断massage是否为空  找到size突然不变的原因
				now := time.Now()
				evt := eventScan{
					Time: now.Format(time.RFC1123),
					//Timestamp: MassageQueue.MsgQueue.Size(),
					Level:   fmt.Sprintf("%s", msg["Level"]),
					Content: fmt.Sprintf("%s", msg["Content"]),
				}
				evtBytes, err := json.Marshal(evt)
				if err != nil {
					golog.Error(err)
					continue
				}
				brokerScanMsg.Notifier <- evtBytes
			}
		}
	}()
	//扫描结果推送
	app.Get("/ssescan", func(ctx context.Context) {
		ctx.HTML(
			`<html><head><title>扫描结果</title>` + scanScript + `</head>
                <body>
                    <h1 id="header">扫描结果</h1>
                    <table id="messagesTable" border="1">
                        <tr>
                            <th>时间</th>
							<th>等级</th>
                            <th>消息</th>
                        </tr>
                    </table>
                </body>
             </html>`)
	})
	app.Get("/eventsscan", brokerScanMsg.ServeHTTP)

	/***** 【配置静态资源】 *****/
	//todo  shortcut.js和core.js需要单独控制，
	app.Favicon("./web/static/favicon.ico", "/favicon.ico") //网站图标
	// 上面可以 这样访问  localhost:8080/favicon.ico
	app.HandleDir("/static", "./web/static") //静态资源

	/******* 【配置网页视图】 *******/
	tmpl := iris.HTML("./web/view", ".html")
	tmpl.Reload(true) //在每个请求上重新加载模板（开发模式）
	app.RegisterView(tmpl)

	/******* 【页面交互】 *******/
	//todo  剔除网页上没有用的功能，修改桌面，最后收尾的时候做
	app.Get("/html", func(ctx context.Context) {
		ctx.ViewData("title", "被动扫描器") //模板中添加数据，前面是模板中的key，后面是要渲染的内容，可以添加多个
		ctx.View("index.html")
	})

	ProxyList := list.New()
	//todo 制作启动代理的接口 返回代理接口
	app.Get("/createProxy", func(ctx context.Context) {
		//启动
		comma := strings.Index(Addr, ":")
		go func() {
			app.Logger().Println("创建成功,代理地址：" + Addr[:comma] + ":" + strconv.Itoa(ProxyPort))
			ProxyList.PushBack(Addr[:comma] + ":" + strconv.Itoa(ProxyPort))
			//http_proxy_simple(Addr[:comma] + ":" + strconv.Itoa(ProxyPort))
			https_Decrypt(Addr[:comma] + ":" + strconv.Itoa(ProxyPort))
		}()
		ProxyPort++ //会先执行这条语句后执行创建代理
		ctx.WriteString("创建成功,代理地址：" + Addr[:comma] + ":" + strconv.Itoa(ProxyPort))
	})

	//todo 制作代理查询接口，查看都那些端口上启动了代理  使用map存储
	app.Get("/showproxy", func(ctx context.Context) {
		HtmlTmp := "<html><head><title>代理列表</title></head>" +
			"<body>" +
			"<h1 id=\"header\">代理列表</h1>" +
			"<table id=\"messagesTable\" border=\"1\">" +
			"<tr>" +
			"<th>序号</th>" +
			"<th>代理地址</th>" +
			"</tr>"
		j := 0
		for i := ProxyList.Front(); i != nil; i = i.Next() {
			HtmlTmp += "<tr><td>" + strconv.Itoa(j) + "</td><td>" + fmt.Sprint(i.Value) + "</td></tr>"
		}
		HtmlTmp += `</table>
		</body>
		</html>`
		ctx.HTML(HtmlTmp)
	})

	app.Run(iris.Addr(addr), iris.WithoutServerError(iris.ErrServerClosed))
}
