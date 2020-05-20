package main

import (
	"bufio"
	"log"
	"os"
	"time"
)

var (
	fileInfo *os.FileInfo
	err      error
)

func main() {
	//测试写入文件
	//t := time.Now()
	//log.Println(t.Format("2020052019"))
	logFileName := time.Now().Format("2006010215") + ".log"
	logDir := "../HttpBackUp/"

	file, err := os.OpenFile(logDir+logFileName, os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("日志文件的打开异常")
		log.Fatal(err)
	}
	defer file.Close()
	buffer := bufio.NewWriter(file)
	bw, err := buffer.WriteString("\n写入字符串")
	if err != nil {
		log.Fatal(err)
	}
	buffer.Flush()
	log.Printf("Bytes written: %d\n", bw)
}
