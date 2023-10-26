package main

import (
	service "kitex_study/kitex_gen/kitex_gen/service/hello"
	"log"
)

func main() {
	svr := service.NewServer(new(HelloImpl))
	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
