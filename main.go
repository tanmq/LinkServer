package main

import (
	"linkServer/config"
	"linkServer/logger"
	"linkServer/server"
)

func main() {
	initLogger()
	port, err := config.CF.Int("serverport")
	if err != nil {
		logger.Error("Can not resolve port. Please setting your listing port.", err)
		return
	}
	logger.Info("linkServer is starting, listening on port ", port)
	s := server.NewTCPServer()
	s.ListenAndServe(port)

}

//read logs setting
func initLogger() {
	logLevel := config.CF.String("log::level")
	switch logLevel {
	case "debug":
		logger.SetLevel(logger.DEBUG)
	case "info":
		logger.SetLevel(logger.INFO)
	case "error":
		logger.SetLevel(logger.ERROR)
	case "fatal":
		logger.SetLevel(logger.FATAL)
	default:
		logger.SetLevel(logger.INFO)
	}

	console := config.CF.DefaultString("log::console", "true")
	if console == "true" {
		logger.SetConsole(true)
	} else {
		logger.SetConsole(false)
	}

	dir := config.CF.DefaultString("log::dir", "logs")
	logfile := config.CF.DefaultString("log::file", "linkServer.log")
	maxFileNum := config.CF.DefaultInt("log::maxfilenum", 10)
	maxFileSize := config.CF.DefaultInt64("log::maxfilesize", 10)

	logger.SetRollingFile(dir, logfile, int32(maxFileNum), maxFileSize, logger.MB)

}
