package main

import (
	"log"
	"regexp"
)

func main() {
	//匹配电话号

	Resp := "<h1>10.203.87.61</h1>" +
		"电话：18904081710<br>" +
		"邮箱：<br>" +
		"身份证号：150428199711100131"
	//开始匹配//ssssdl@qq.com

	//IP
	findIP := regexp.MustCompile("[\\d]+\\.[\\d]+\\.[\\d]+\\.[\\d]+")
	IP := findIP.FindString(Resp)
	log.Println("IP:" + IP)

	//电话 身份证号增加至相关位数
	findPH := regexp.MustCompile("\\d{11}")
	PH := findPH.FindString(Resp)
	log.Println("PH:" + PH)

	//邮箱
	findMail := regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
	Mail := findMail.FindString(Resp)
	log.Println(Mail)
}
