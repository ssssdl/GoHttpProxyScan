package plugin

import (
	"../MassageQueue"
	"net/http"
	"strconv"
)

type plugin1 struct{}

func (p *plugin1) request(Req *http.Request) {
	//time.Sleep(200*time.Millisecond)
	MassageQueue.MsgQueue.Put(Req.Host)
}
func (p *plugin1) respones(Resp string) {
	//time.Sleep(200*time.Millisecond)
	MassageQueue.MsgQueue.Put(strconv.Itoa(len(Resp)))
}
