package main

import (
	"fmt"
	"strconv"
	"time"

	service "github.com/budimanlai/go-service"
)

var (
	isBreak bool
	node    string
)

func StartFunc(context service.ServiceContext) {
	isBreak = false
	node = strconv.Itoa(context.Node)

	log("Service run as node: " + node)

	type Result struct {
		ID        int
		Handphone string
	}
	var result Result
	db := context.DB()
	db.Raw("SELECT id, handphone FROM user WHERE id = ?", 14).Scan(&result)

	fmt.Println(result)

	count := 0
	for {
		log("Infinite loop: " + strconv.Itoa(count))
		time.Sleep(time.Second * 2)
		count++

		if isBreak {
			log("Service stoped")
			break
		}
	}
}

func StopFunc(context service.ServiceContext) {
	log("try to stop node")
	defer func() {
		isBreak = true
	}()
}

func log(msg string) {
	t := time.Now()

	s := "[Node: " + node + " " + t.Format("2006-01-02 15:04:05") + "] " + msg
	fmt.Println(s)
}
