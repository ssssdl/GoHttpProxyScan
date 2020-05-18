package plugin

import (
	"../MassageQueue"
	"net/http"
)

type plugin2 struct{}

func (p *plugin2) request(Req *http.Request) {
	//time.Sleep(200*time.Millisecond)
	MassageQueue.MsgQueue.Put("plugin2 request")
}
func (p *plugin2) respones(Resp string) {
	//time.Sleep(200*time.Millisecond)
	MassageQueue.MsgQueue.Put(string([]byte(Resp)[:20]))
}
