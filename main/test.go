package main

import (
	"container/list"
	"log"
)

func main() {
	l := list.New()
	l.PushBack("1")
	l.PushBack("2")
	l.PushBack("3")
	for i := l.Front(); i != nil; i = i.Next() {
		log.Println(i.Value)
	}
}
