package main

import (

	"linkServer/config"
	"linkServer/server"
	"linkServer/logger"
)

func main() {
	port, err := config.CF.Int("serverport")
	if err != nil {
		logger.Error("Can not resolve port. Please setting your listing port.", err)
		return
	}
	logger.Info("linkServer is starting, listening on port ", port)
	s := server.NewTcpServer()
	s.ListenAndServe(port)

}
