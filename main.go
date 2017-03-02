package main

import (

	"fmt"


	"linkServer/config"
	"linkServer/server"
)

func main() {
	cf := config.NewConfiger()
	fmt.Println(cf.Int("serverport"))

	s := server.NewTcpServer()
	s.ListenAndServe(9999)



}


