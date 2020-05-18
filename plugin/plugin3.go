package plugin

import (
	"../MassageQueue"
	"net/http"
)

type plugin3 struct{}

func (p *plugin3) request(Req *http.Request) {
	//time.Sleep(200*time.Millisecond)
	MassageQueue.MsgQueue.Put("plugin3 request")
}
func (p *plugin3) respones(Resp string) {
	//time.Sleep(200*time.Millisecond)
	MassageQueue.MsgQueue.Put("plugin3 resp")
}
