package MassageQueue

import (
	"log"
	"sync"
)

//使用链表实现消息队列

//消息队列

type MassageQueue struct {
	FirstMsg     *Massage
	LastMsg      *Massage
	size         int
	sync.RWMutex //sync.RWMutex 使用互斥锁保证线程安全
}

//消息节点
type Massage struct {
	NextMsg *Massage
	content string //消息内容
}

func New() *MassageQueue {
	MsgQ := new(MassageQueue)
	//MsgQ.Put("basic")
	return MsgQ
}

//func 后面的括号的意思是实现MassageQueue的接口方法
//向消息队列中添加消息
//参数：消息内容（content string）
//todo：将消息写入文件作为备份
func (MassageQueue *MassageQueue) Put(content string) {
	MassageQueue.Lock()         //锁定队列
	defer MassageQueue.Unlock() //使用defer来自动解锁消息队列
	msg := &Massage{            //创建消息节点
		content: content,
	}
	if MassageQueue.FirstMsg == nil { //如果消息队列为空则将msg节点作为第一个节点
		MassageQueue.FirstMsg = msg
	} else { //否则将msg节点追加到最后一个节点末尾
		MassageQueue.LastMsg.NextMsg = msg
	}
	MassageQueue.LastMsg = msg //更新最后一个节点为新添加进来的msg节点
	MassageQueue.size++        //消息数自增
}

//获取消息
//返回：消息内容（string）
func (MassageQueue *MassageQueue) Get() string {
	MassageQueue.Lock()               //锁定消息队列
	defer MassageQueue.Unlock()       //使用defer来解锁消息队列
	if MassageQueue.FirstMsg == nil { //如果FirstMsg是空的说明消息队列为空，直接返回空
		log.Println(MassageQueue.size)
		return ""
	}
	content := MassageQueue.FirstMsg.content              //获取第一个消息内容
	MassageQueue.FirstMsg = MassageQueue.FirstMsg.NextMsg //第一个消息已经被读取，将FirstMsg指向下一个消息
	MassageQueue.size--                                   //消息数自减
	return content                                        //返回消息内容
}

//获取消息个函数
//返回：消息个数（int）
func (MassageQueue *MassageQueue) Size() int {
	MassageQueue.Lock()         //锁定消息队列
	defer MassageQueue.Unlock() //使用defer解锁消息队列
	return MassageQueue.size    //返回消息长度
}

var MsgQueue = New() //用于存放扫描消息

var HttpHistoryQueue = New() //用于存放http历史记录
